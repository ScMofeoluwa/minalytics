package main

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mileusna/useragent"
	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/viper"

	database "github.com/ScMofeoluwa/minalytics/database/sqlc"
	"github.com/ScMofeoluwa/minalytics/types"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidTrackingID = errors.New("invalid tracking ID")
	ErrUserNotCreated    = errors.New("failed to create user")
	ErrEventNotCreated   = errors.New("failed to create event")
)

type AnalyticsService struct {
	Queries *database.Queries
	GeoDB   *geoip2.Reader
}

func NewAnalyticsService(queries *database.Queries, geoDB *geoip2.Reader) *AnalyticsService {
	return &AnalyticsService{
		Queries: queries,
		GeoDB:   geoDB,
	}
}

func (s *AnalyticsService) SignIn(ctx context.Context, email string) (string, error) {
	userId, err := s.Queries.FindOrCreateUser(ctx, email)
	if err != nil {
		return "", ErrUserNotCreated
	}
	return CreateJWT(userId.String())
}

func (s *AnalyticsService) TrackEvent(ctx context.Context, data types.EventPayload) error {
	if _, err := s.Queries.FindUserByTrackingID(ctx, data.Tracking.TrackingId); err != nil {
		return ErrInvalidTrackingID
	}

	uaDetails := s.ParseUserAgent(data.Tracking.Ua)

	params := database.CreateEventParams{
		VisitorID:       data.Tracking.VisitorId,
		TrackingID:      data.Tracking.TrackingId,
		EventType:       data.Type,
		Url:             &data.Tracking.Url,
		Referrer:        &data.Tracking.Referrer,
		Country:         data.Tracking.Country,
		Browser:         uaDetails.Browser,
		Device:          uaDetails.Device,
		OperatingSystem: uaDetails.OperatingSystem,
		Details:         data.Tracking.Details,
	}

	if err := s.Queries.CreateEvent(ctx, params); err != nil {
		return ErrEventNotCreated
	}
	return nil
}

func (s *AnalyticsService) ResolveGeoLocation(remoteAddr string) (types.GeoLocation, error) {
	ip := net.ParseIP(remoteAddr)
	record, err := s.GeoDB.City(ip)
	if err != nil {
		return types.GeoLocation{}, err
	}

	geoLocation := types.GeoLocation{
		Country:   record.Country.Names["en"],
		City:      record.City.Names["en"],
		Longitude: record.Location.Longitude,
		Latitude:  record.Location.Latitude,
	}

	return geoLocation, nil
}

func (s *AnalyticsService) ParseUserAgent(ua string) types.UserAgentDetails {
	parsedUA := useragent.Parse(ua)
	return types.UserAgentDetails{
		Browser:         parsedUA.Name,
		Device:          parsedUA.Device,
		OperatingSystem: parsedUA.OS,
	}
}
func CreateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(viper.GetString("TokenSecret"))
	return token.SignedString(secret)
}

func VerifyJWT(token string) (jwt.MapClaims, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(viper.GetString("TokenSecret")), nil
	}

	jwtToken, err := jwt.Parse(token, keyfunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
