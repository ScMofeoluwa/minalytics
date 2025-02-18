package server

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
)

var ErrInvalidToken = errors.New("invalid token")
var ErrAppNotFound = errors.New("app not found")
var ErrAppExists = errors.New("app already exists")

type analyticsService struct {
	Querier database.Querier
	GeoDB   *geoip2.Reader
}

func NewAnalyticsService(querier database.Querier, geoDB *geoip2.Reader) AnalyticsService {
	return &analyticsService{
		Querier: querier,
		GeoDB:   geoDB,
	}
}

func (s *analyticsService) SignIn(ctx context.Context, email string) (string, error) {
	userId, err := s.Querier.GetOrCreateUser(ctx, email)
	if err != nil {
		return "", err
	}
	return CreateJWT(userId.String())
}

func (s *analyticsService) TrackEvent(ctx context.Context, data EventPayload) error {
	if _, err := s.Querier.GetAppByTrackingID(ctx, data.Tracking.TrackingID); err != nil {
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

	if err := s.Querier.CreateEvent(ctx, params); err != nil {
		return err
	}
	return nil
}

func (s *analyticsService) CreateApp(ctx context.Context, userID uuid.UUID, name string) (*App, error) {
	params := database.CheckAppExistsParams{
		UserID: userID,
		Name:   name,
	}
	if app, err := s.Querier.CheckAppExists(ctx, params); err == nil && app.TrackingID != uuid.Nil {
		return &App{}, ErrAppExists
	}

	createParams := database.CreateAppParams{
		UserID: userID,
		Name:   name,
	}

	app_, err := s.Querier.CreateApp(ctx, createParams)
	if err != nil {
		return &App{}, err
	}

	app := &App{
		Name:       app_.Name,
		TrackingID: app_.TrackingID,
		CreatedAt:  app_.CreatedAt.Time,
	}
	return app, nil
}

func (s *analyticsService) GetApps(ctx context.Context, userID uuid.UUID) ([]App, error) {
	apps_, err := s.Querier.GetApps(ctx, userID)
	if err != nil {
		return []App{}, err
	}

	apps := make([]App, 0, len(apps_))
	for _, row := range apps_ {
		apps = append(apps, App{
			Name:       row.Name,
			CreatedAt:  row.CreatedAt.Time,
			TrackingID: row.TrackingID,
		})
	}
	return apps, nil
}

func (s *analyticsService) UpdateApp(ctx context.Context, name string, trackingID, userID uuid.UUID) (*App, error) {
	if err := s.ValidateAppAccess(ctx, userID, trackingID); err != nil {
		return &App{}, err
	}

	params := database.UpdateAppParams{
		TrackingID: trackingID,
		Name:       name,
	}

	app_, err := s.Querier.UpdateApp(ctx, params)
	if err != nil {
		return &App{}, err
	}

	app := &App{
		Name:       app_.Name,
		TrackingID: app_.TrackingID,
		CreatedAt:  app_.CreatedAt.Time,
	}
	return app, nil
}

func (s *analyticsService) DeleteApp(ctx context.Context, trackingID, userID uuid.UUID) error {
	if err := s.ValidateAppAccess(ctx, userID, trackingID); err != nil {
		return err
	}

	return s.Querier.DeleteApp(ctx, trackingID)
}

func (s *analyticsService) GetReferrals(ctx context.Context, data RequestPayload) ([]ReferralStats, error) {
	params := database.GetReferralsParams{
		TrackingID: data.TrackingID,
		Column2:    data.StartDate,
		Column3:    data.EndDate,
	}

	stats, err := s.Querier.GetReferrals(ctx, params)
	if err != nil {
		return []ReferralStats{}, err
	}

	referralStats := make([]ReferralStats, 0, len(stats))
	for _, row := range stats {
		referralStats = append(referralStats, ReferralStats{
			Referrer:     *row.Referrer,
			VisitorCount: int(row.VisitorCount),
		})
	}

	return referralStats, nil
}

func (s *analyticsService) GetPages(ctx context.Context, data RequestPayload) ([]PageStats, error) {
	params := database.GetPagesParams{
		TrackingID: data.TrackingID,
		Column2:    data.StartDate,
		Column3:    data.EndDate,
	}

	stats, err := s.Querier.GetPages(ctx, params)
	if err != nil {
		return []PageStats{}, err
	}

	pageStats := make([]PageStats, 0, len(stats))
	for _, row := range stats {
		url, err := url.Parse(*row.Url)
		if err != nil {
			return nil, err
		}

		path := url.Path
		if url.Path == "" {
			path = "/"
		}

		pageStats = append(pageStats, PageStats{
			Path:         path,
			VisitorCount: int(row.VisitorCount),
		})
	}

	return pageStats, nil
}

func (s *analyticsService) GetBrowsers(ctx context.Context, data RequestPayload) ([]BrowserStats, error) {
	params := database.GetBrowsersParams{
		TrackingID: data.TrackingID,
		Column2:    data.StartDate,
		Column3:    data.EndDate,
	}

	stats, err := s.Querier.GetBrowsers(ctx, params)
	if err != nil {
		return []BrowserStats{}, err
	}

	browserStats := make([]BrowserStats, 0, len(stats))
	for _, row := range stats {
		browserStats = append(browserStats, BrowserStats{
			Browser:    row.Browser,
			Percentage: int(row.Percentage),
		})
	}

	return browserStats, nil
}

func (s *analyticsService) GetCountries(ctx context.Context, data RequestPayload) ([]CountryStats, error) {
	params := database.GetCountriesParams{
		TrackingID: data.TrackingID,
		Column2:    data.StartDate,
		Column3:    data.EndDate,
	}

	stats, err := s.Querier.GetCountries(ctx, params)
	if err != nil {
		return []CountryStats{}, err
	}

	countryStats := make([]CountryStats, 0, len(stats))
	for _, row := range stats {
		countryStats = append(countryStats, CountryStats{
			Country:    row.Country,
			Percentage: int(row.Percentage),
		})
	}

	return countryStats, nil
}

func (s *analyticsService) GetDevices(ctx context.Context, data RequestPayload) ([]DeviceStats, error) {
	params := database.GetDevicesParams{
		TrackingID: data.TrackingID,
		Column2:    data.StartDate,
		Column3:    data.EndDate,
	}

	stats, err := s.Querier.GetDevices(ctx, params)
	if err != nil {
		return []DeviceStats{}, err
	}

	deviceStats := make([]DeviceStats, 0, len(stats))
	for _, row := range stats {
		deviceStats = append(deviceStats, DeviceStats{
			Device:     row.Device,
			Percentage: int(row.Percentage),
		})
	}

	return deviceStats, nil
}

func (s *analyticsService) GetOS(ctx context.Context, data RequestPayload) ([]OSStats, error) {
	params := database.GetOSParams{
		TrackingID: data.TrackingID,
		Column2:    data.StartDate,
		Column3:    data.EndDate,
	}

	stats, err := s.Querier.GetOS(ctx, params)
	if err != nil {
		return []OSStats{}, err
	}

	osstats := make([]OSStats, 0, len(stats))
	for _, row := range stats {
		osstats = append(osstats, OSStats{
			OS:         row.OperatingSystem,
			Percentage: int(row.Percentage),
		})
	}

	return osstats, nil
}

func (s *analyticsService) GetVisitors(ctx context.Context, data RequestPayload) ([]VisitorStats, error) {
	params := database.GetVisitorsParams{
		TrackingID: data.TrackingID,
		Column2:    data.StartDate,
		Column3:    data.EndDate,
		TimeBucket: data.BucketSize,
	}

	stats, err := s.Querier.GetVisitors(ctx, params)
	if err != nil {
		return []VisitorStats{}, err
	}

	visitorStats := make([]VisitorStats, 0, len(stats))
	for _, row := range stats {
		visitorStats = append(visitorStats, VisitorStats{
			Time:     row.Time.Time.String(),
			Visitors: int(row.Visitors),
		})
	}

	return visitorStats, nil
}

func (s *analyticsService) GetPageViews(ctx context.Context, data RequestPayload) ([]PageViewStats, error) {
	params := database.GetPageViewsParams{
		TrackingID: data.TrackingID,
		Column2:    data.StartDate,
		Column3:    data.EndDate,
		TimeBucket: data.BucketSize,
	}

	stats, err := s.Querier.GetPageViews(ctx, params)
	if err != nil {
		return []PageViewStats{}, err
	}

	pageViewStats := make([]PageViewStats, 0, len(stats))
	for _, row := range stats {
		pageViewStats = append(pageViewStats, PageViewStats{
			Time:  row.Time.Time.String(),
			Views: int(row.Views),
		})
	}

	return pageViewStats, nil
}

func (s *analyticsService) ResolveGeoLocation(remoteAddr string) (*GeoLocation, error) {
	ip := net.ParseIP(remoteAddr)
	record, err := s.GeoDB.City(ip)
	if err != nil {
		return &GeoLocation{}, err
	}

	geoLocation := &GeoLocation{
		Country:   record.Country.Names["en"],
		City:      record.City.Names["en"],
		Longitude: record.Location.Longitude,
		Latitude:  record.Location.Latitude,
	}

	return geoLocation, nil
}

func (s *analyticsService) ParseUserAgent(ua string) *UserAgentDetails {
	parsedUA := useragent.Parse(ua)
	return &UserAgentDetails{
		Browser:         parsedUA.Name,
		Device:          parsedUA.Device,
		OperatingSystem: parsedUA.OS,
	}
}

func (s *analyticsService) ValidateAppAccess(ctx context.Context, userID, trackingID uuid.UUID) error {
	app, err := s.Querier.GetAppByTrackingID(ctx, trackingID)
	if err != nil {
		return err
	}
	if app.UserID != userID {
		return ErrAppNotFound
	}
	return nil
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
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
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
