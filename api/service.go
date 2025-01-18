package main

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/viper"

	database "github.com/ScMofeoluwa/minalytics/database/sqlc"
	"github.com/ScMofeoluwa/minalytics/types"
)

var ErrInvalidToken = errors.New("invalid token")

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
	userId, err := s.Queries.GetOrCreateUser(ctx, email)
	if err != nil {
		return "", err
	}
	return CreateJWT(userId.String())
}

func (s *AnalyticsService) TrackEvent(ctx context.Context, data types.EventPayload) error {
	if _, err := s.Queries.GetUserByTrackingID(ctx, data.Tracking.TrackingID); err != nil {
		return err
	}

	uaDetails := s.ParseUserAgent(data.Tracking.Ua)

	params := database.CreateEventParams{
		VisitorID:       data.Tracking.VisitorID,
		TrackingID:      data.Tracking.TrackingID,
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
		return err
	}
	return nil
}

func (s *AnalyticsService) GetTrackingID(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := s.Queries.GetUserTrackingID(ctx, userID)
	if err != nil {
		return "", err
	}

	return user.String(), nil
}

func (s *AnalyticsService) GetReferrals(ctx context.Context, data types.RequestPayload) ([]types.ReferralStats, error) {
	params := database.GetReferralsParams{
		ID:          data.UserID,
		Timestamp:   data.StartTime,
		Timestamp_2: data.EndTime,
	}

	stats, err := s.Queries.GetReferrals(ctx, params)
	if err != nil {
		return []types.ReferralStats{}, err
	}

	referralStats := make([]types.ReferralStats, 0, len(stats))
	for _, row := range stats {
		referralStats = append(referralStats, types.ReferralStats{
			Referrer:     *row.Referrer,
			VisitorCount: int(row.VisitorCount),
		})
	}

	return referralStats, nil
}

func (s *AnalyticsService) GetPages(ctx context.Context, data types.RequestPayload) ([]types.PageStats, error) {
	params := database.GetPagesParams{
		ID:          data.UserID,
		Timestamp:   data.StartTime,
		Timestamp_2: data.EndTime,
	}

	stats, err := s.Queries.GetPages(ctx, params)
	if err != nil {
		return []types.PageStats{}, err
	}

	pageStats := make([]types.PageStats, 0, len(stats))
	for _, row := range stats {
		url, err := url.Parse(*row.Url)
		if err != nil {
			return nil, err
		}

		pageStats = append(pageStats, types.PageStats{
			Path:         url.Path,
			VisitorCount: int(row.VisitorCount),
		})
	}

	return pageStats, nil
}

func (s *AnalyticsService) GetBrowsers(ctx context.Context, data types.RequestPayload) ([]types.BrowserStats, error) {
	params := database.GetBrowsersParams{
		ID:          data.UserID,
		Timestamp:   data.StartTime,
		Timestamp_2: data.EndTime,
	}

	stats, err := s.Queries.GetBrowsers(ctx, params)
	if err != nil {
		return []types.BrowserStats{}, err
	}

	browserStats := make([]types.BrowserStats, 0, len(stats))
	for _, row := range stats {
		browserStats = append(browserStats, types.BrowserStats{
			Browser:    row.Browser,
			Percentage: int(row.Percentage),
		})
	}

	return browserStats, nil
}

func (s *AnalyticsService) GetCountries(ctx context.Context, data types.RequestPayload) ([]types.CountryStats, error) {
	params := database.GetCountriesParams{
		ID:          data.UserID,
		Timestamp:   data.StartTime,
		Timestamp_2: data.EndTime,
	}

	stats, err := s.Queries.GetCountries(ctx, params)
	if err != nil {
		return []types.CountryStats{}, err
	}

	countryStats := make([]types.CountryStats, 0, len(stats))
	for _, row := range stats {
		countryStats = append(countryStats, types.CountryStats{
			Country:    row.Country,
			Percentage: int(row.Percentage),
		})
	}

	return countryStats, nil
}

func (s *AnalyticsService) GetDevices(ctx context.Context, data types.RequestPayload) ([]types.DeviceStats, error) {
	params := database.GetDevicesParams{
		ID:          data.UserID,
		Timestamp:   data.StartTime,
		Timestamp_2: data.EndTime,
	}

	stats, err := s.Queries.GetDevices(ctx, params)
	if err != nil {
		return []types.DeviceStats{}, err
	}

	deviceStats := make([]types.DeviceStats, 0, len(stats))
	for _, row := range stats {
		deviceStats = append(deviceStats, types.DeviceStats{
			Device:     row.Device,
			Percentage: int(row.Percentage),
		})
	}

	return deviceStats, nil
}

func (s *AnalyticsService) GetOS(ctx context.Context, data types.RequestPayload) ([]types.OSStats, error) {
	params := database.GetOSParams{
		ID:          data.UserID,
		Timestamp:   data.StartTime,
		Timestamp_2: data.EndTime,
	}

	stats, err := s.Queries.GetOS(ctx, params)
	if err != nil {
		return []types.OSStats{}, err
	}

	OSStats := make([]types.OSStats, 0, len(stats))
	for _, row := range stats {
		OSStats = append(OSStats, types.OSStats{
			OS:         row.OperatingSystem,
			Percentage: int(row.Percentage),
		})
	}

	return OSStats, nil
}

func (s *AnalyticsService) ResolveGeoLocation(remoteAddr string) (*types.GeoLocation, error) {
	ip := net.ParseIP(remoteAddr)
	record, err := s.GeoDB.City(ip)
	if err != nil {
		return &types.GeoLocation{}, err
	}

	geoLocation := &types.GeoLocation{
		Country:   record.Country.Names["en"],
		City:      record.City.Names["en"],
		Longitude: record.Location.Longitude,
		Latitude:  record.Location.Latitude,
	}

	return geoLocation, nil
}

func (s *AnalyticsService) ParseUserAgent(ua string) *types.UserAgentDetails {
	parsedUA := useragent.Parse(ua)
	return &types.UserAgentDetails{
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
