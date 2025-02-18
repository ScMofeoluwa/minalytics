package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	database "github.com/ScMofeoluwa/minalytics/database/sqlc"
	"github.com/ScMofeoluwa/minalytics/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	mockRepo *mocks.Querier
	geoDB    *geoip2.Reader
	service  AnalyticsService
	ctx      context.Context
}

func (suite *ServiceSuite) SetupSuite() {
	suite.ctx = context.Background()

	//initialise Geo Database
	geoDB, err := geoip2.Open("database/GeoLite2-City.mmdb")
	if err != nil {
		suite.FailNowf("failed to open GeoIP2 database", err.Error())
	}
	suite.geoDB = geoDB

	suite.mockRepo = mocks.NewQuerier(suite.T())
	suite.service = NewAnalyticsService(suite.mockRepo, geoDB)
}

func (suite *ServiceSuite) TearDownSuite() {
	suite.geoDB.Close()
}

func (suite *ServiceSuite) TestSignIn() {
	mockUser := struct {
		email string
		id    uuid.UUID
	}{
		email: "testuser@gmail.com",
		id:    uuid.New(),
	}

	suite.mockRepo.EXPECT().GetOrCreateUser(mock.Anything, mock.Anything).Return(mockUser.id, nil).Once()
	token, err := suite.service.SignIn(suite.ctx, mockUser.email)

	suite.NoError(err)
	suite.NotEmpty(token)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestTrackEvent() {
	testCases := []struct {
		name        string
		data        EventPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "event tracked successfully",
			data: EventPayload{
				Type: "pageview",
				Tracking: TrackingData{
					TrackingID: uuid.New(),
					VisitorID:  faker.UUIDDigit(),
					Ua:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
					Url:        faker.URL(),
					Referrer:   faker.URL(),
					Country:    faker.GetCountryInfo().Name,
					Details:    map[string]interface{}{},
				},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetAppByTrackingID(mock.Anything, mock.Anything).Return(database.App{}, nil).Once()
				suite.mockRepo.EXPECT().CreateEvent(mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "app not found",
			data: EventPayload{
				Type: "pageview",
				Tracking: TrackingData{
					TrackingID: uuid.New(),
					VisitorID:  faker.UUIDDigit(),
					Ua:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
					Url:        faker.URL(),
					Referrer:   faker.URL(),
					Country:    faker.GetCountryInfo().Name,
					Details:    map[string]interface{}{},
				},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetAppByTrackingID(mock.Anything, mock.Anything).Return(database.App{}, errors.New("app not found")).Once()
			},
			expectedErr: errors.New("app not found"),
		},
		{
			name: "failed to create event",
			data: EventPayload{
				Type: "pageview",
				Tracking: TrackingData{
					TrackingID: uuid.New(),
					VisitorID:  faker.UUIDDigit(),
					Ua:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
					Url:        faker.URL(),
					Referrer:   faker.URL(),
					Country:    faker.GetCountryInfo().Name,
					Details:    map[string]interface{}{},
				},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetAppByTrackingID(mock.Anything, mock.Anything).Return(database.App{}, nil).Once()
				suite.mockRepo.EXPECT().CreateEvent(mock.Anything, mock.Anything).Return(errors.New("create event failed")).Once()
			},
			expectedErr: errors.New("create event failed"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			err := suite.service.TrackEvent(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestCreateApp() {
	testCases := []struct {
		name        string
		userID      uuid.UUID
		appName     string
		mockSetup   func()
		expectedErr error
	}{
		{
			name:    "app created successfully",
			userID:  uuid.New(),
			appName: "Test App",
			mockSetup: func() {
				suite.mockRepo.EXPECT().CheckAppExists(mock.Anything, mock.Anything).Return(database.App{}, nil).Once()
				suite.mockRepo.EXPECT().CreateApp(mock.Anything, mock.Anything).Return(database.App{
					Name:       "Test App",
					TrackingID: uuid.New(),
					CreatedAt:  sql.NullTime{},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:    "app already exists",
			userID:  uuid.New(),
			appName: "Test App",
			mockSetup: func() {
				suite.mockRepo.EXPECT().CheckAppExists(mock.Anything, mock.Anything).Return(database.App{
					Name:       "Test App",
					TrackingID: uuid.New(),
					CreatedAt:  sql.NullTime{},
				}, nil).Once()
			},
			expectedErr: errors.New("app already exists"),
		},
		{
			name:    "failed to create app",
			userID:  uuid.New(),
			appName: "Test App",
			mockSetup: func() {
				suite.mockRepo.EXPECT().CheckAppExists(mock.Anything, mock.Anything).Return(database.App{}, nil).Once()
				suite.mockRepo.EXPECT().CreateApp(mock.Anything, mock.Anything).Return(database.App{}, errors.New("create app failed")).Once()
			},
			expectedErr: errors.New("create app failed"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			app, err := suite.service.CreateApp(suite.ctx, tc.userID, tc.appName)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.Equal(tc.appName, app.Name)
			suite.NotEmpty(app.TrackingID)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetApps() {
	testCases := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expectedErr error
	}{
		{
			name:   "app successfully retrieved",
			userID: uuid.New(),
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetApps(mock.Anything, mock.Anything).Return([]database.App{
					{
						Name:       faker.Name(),
						TrackingID: uuid.New(),
						CreatedAt:  sql.NullTime{},
					},
					{
						Name:       faker.Name(),
						TrackingID: uuid.New(),
						CreatedAt:  sql.NullTime{},
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "failed to fetch apps",
			userID: uuid.New(),
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetApps(mock.Anything, mock.Anything).Return([]database.App{}, errors.New("failed to fetch apps")).Once()
			},
			expectedErr: errors.New("failed to fetch apps"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			apps, err := suite.service.GetApps(suite.ctx, tc.userID)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(apps)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetReferrals() {
	testCases := []struct {
		name        string
		data        RequestPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "referrals successfully retrieved",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetReferrals(mock.Anything, mock.Anything).Return([]database.GetReferralsRow{
					{
						Referrer:     stringPtr(faker.URL()),
						VisitorCount: 10,
					},
					{
						Referrer:     stringPtr(faker.URL()),
						VisitorCount: 20,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "failed to fetch referrals",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetReferrals(mock.Anything, mock.Anything).Return([]database.GetReferralsRow{}, errors.New("failed to fetch referrals")).Once()
			},
			expectedErr: errors.New("failed to fetch referrals"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			referrals, err := suite.service.GetReferrals(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(referrals)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetPages() {
	testCases := []struct {
		name        string
		data        RequestPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "pages successfully retrieved",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetPages(mock.Anything, mock.Anything).Return([]database.GetPagesRow{
					{
						Url:          stringPtr(faker.URL()),
						VisitorCount: 10,
					},
					{
						Url:          stringPtr(faker.URL()),
						VisitorCount: 20,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "failed to fetch pages",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetPages(mock.Anything, mock.Anything).Return([]database.GetPagesRow{}, errors.New("failed to fetch pages")).Once()
			},
			expectedErr: errors.New("failed to fetch pages"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			pages, err := suite.service.GetPages(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(pages)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetBrowsers() {
	testCases := []struct {
		name        string
		data        RequestPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "browsers successfully retrieved",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetBrowsers(mock.Anything, mock.Anything).Return([]database.GetBrowsersRow{
					{
						Browser:    "Chrome",
						Percentage: 50,
					},
					{
						Browser:    "Firefox",
						Percentage: 30,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "failed to fetch browsers",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetBrowsers(mock.Anything, mock.Anything).Return([]database.GetBrowsersRow{}, errors.New("failed to fetch browsers")).Once()
			},
			expectedErr: errors.New("failed to fetch browsers"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			browsers, err := suite.service.GetBrowsers(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(browsers)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetCountries() {
	testCases := []struct {
		name        string
		data        RequestPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "countries successfully retrieved",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetCountries(mock.Anything, mock.Anything).Return([]database.GetCountriesRow{
					{
						Country:    faker.GetCountryInfo().Name,
						Percentage: 50,
					},
					{
						Country:    faker.GetCountryInfo().Name,
						Percentage: 30,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "failed to fetch countries",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetCountries(mock.Anything, mock.Anything).Return([]database.GetCountriesRow{}, errors.New("failed to fetch countries")).Once()
			},
			expectedErr: errors.New("failed to fetch countries"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			countries, err := suite.service.GetCountries(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(countries)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetDevices() {
	testCases := []struct {
		name        string
		data        RequestPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "devices successfully retrieved",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetDevices(mock.Anything, mock.Anything).Return([]database.GetDevicesRow{
					{
						Device:     "Desktop",
						Percentage: 50,
					},
					{
						Device:     "Mobile",
						Percentage: 30,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "failed to fetch devices",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetDevices(mock.Anything, mock.Anything).Return([]database.GetDevicesRow{}, errors.New("failed to fetch devices")).Once()
			},
			expectedErr: errors.New("failed to fetch devices"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			devices, err := suite.service.GetDevices(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(devices)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetOS() {
	testCases := []struct {
		name        string
		data        RequestPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "os successfully retrieved",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetOS(mock.Anything, mock.Anything).Return([]database.GetOSRow{
					{
						OperatingSystem: "Windows",
						Percentage:      50,
					},
					{
						OperatingSystem: "MacOS",
						Percentage:      30,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "failed to fetch os",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetOS(mock.Anything, mock.Anything).Return([]database.GetOSRow{}, errors.New("failed to fetch os")).Once()
			},
			expectedErr: errors.New("failed to fetch os"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			osStats, err := suite.service.GetOS(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(osStats)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetVisitors() {
	testCases := []struct {
		name        string
		data        RequestPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "visitors successfully retrieved",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
				BucketSize: "1h",
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetVisitors(mock.Anything, mock.Anything).Return([]database.GetVisitorsRow{
					{
						Time:     sql.NullTime{Time: time.Now(), Valid: true},
						Visitors: 100,
					},
					{
						Time:     sql.NullTime{Time: time.Now().Add(-1 * time.Hour), Valid: true},
						Visitors: 50,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "failed to fetch visitors",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
				BucketSize: "1h",
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetVisitors(mock.Anything, mock.Anything).Return([]database.GetVisitorsRow{}, errors.New("failed to fetch visitors")).Once()
			},
			expectedErr: errors.New("failed to fetch visitors"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			visitors, err := suite.service.GetVisitors(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(visitors)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestGetPageViews() {
	testCases := []struct {
		name        string
		data        RequestPayload
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "pageviews successfully retrieved",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
				BucketSize: "1h",
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetPageViews(mock.Anything, mock.Anything).Return([]database.GetPageViewsRow{
					{
						Time:  sql.NullTime{Time: time.Now(), Valid: true},
						Views: 200,
					},
					{
						Time:  sql.NullTime{Time: time.Now().Add(-1 * time.Hour), Valid: true},
						Views: 100,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "failed to fetch pageviews",
			data: RequestPayload{
				TrackingID: uuid.New(),
				StartDate:  sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
				EndDate:    sql.NullTime{Time: time.Now(), Valid: true},
				BucketSize: "1h",
			},
			mockSetup: func() {
				suite.mockRepo.EXPECT().GetPageViews(mock.Anything, mock.Anything).Return([]database.GetPageViewsRow{}, errors.New("failed to fetch pageviews")).Once()
			},
			expectedErr: errors.New("failed to fetch pageviews"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			pageViews, err := suite.service.GetPageViews(suite.ctx, tc.data)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.NotEmpty(pageViews)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *ServiceSuite) TestResolveGeoLocation() {
	testCases := []struct {
		name        string
		ipAddress   string
		expectError bool
	}{
		{
			name:        "valid IP address",
			ipAddress:   "8.8.8.8",
			expectError: false,
		},
		{
			name:        "invalid IP address",
			ipAddress:   "invalid-ip",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			location, err := suite.service.ResolveGeoLocation(tc.ipAddress)

			if tc.expectError {
				suite.Error(err)
				suite.Empty(location)
				return
			}
			suite.NoError(err)
			suite.NotEmpty(location.Country)
			suite.NotZero(location.Latitude)
			suite.NotZero(location.Longitude)
		})
	}
}

func (suite *ServiceSuite) TestParseUserAgent() {
	testCases := []struct {
		name      string
		userAgent string
		expected  *UserAgentDetails
	}{
		{
			name:      "Chrome Windows browser",
			userAgent: "Mozilla/5.0 (Windows; Windows NT 6.2; WOW64) AppleWebKit/533.25 (KHTML, like Gecko) Chrome/53.0.1644.308 Safari/536",
			expected: &UserAgentDetails{
				Browser:         "Chrome",
				Device:          "",
				OperatingSystem: "Windows",
			},
		},
		{
			name:      "Chrome iPhone browser",
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_2; like Mac OS X) AppleWebKit/601.7 (KHTML, like Gecko)  Chrome/48.0.2837.308 Mobile Safari/600.3",
			expected: &UserAgentDetails{
				Browser:         "Chrome",
				Device:          "iPhone",
				OperatingSystem: "iOS",
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := suite.service.ParseUserAgent(tc.userAgent)
			suite.NotNil(result)
			suite.Equal(tc.expected.Browser, result.Browser)
			suite.Equal(tc.expected.Device, result.Device)
			suite.Equal(tc.expected.OperatingSystem, result.OperatingSystem)
		})
	}
}

func (suite *ServiceSuite) TestValidateAppAccess() {
	testCases := []struct {
		name        string
		userID      uuid.UUID
		trackingID  uuid.UUID
		mockSetup   func(trackingID uuid.UUID)
		expectedErr error
	}{
		{
			name:       "successful validation",
			userID:     uuid.New(),
			trackingID: uuid.New(),
			mockSetup: func(trackingID uuid.UUID) {
				suite.mockRepo.EXPECT().GetAppByTrackingID(mock.Anything, mock.Anything).Return(database.App{TrackingID: trackingID}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:       "app not found",
			userID:     uuid.New(),
			trackingID: uuid.New(),
			mockSetup: func(trackingID uuid.UUID) {
				suite.mockRepo.EXPECT().GetAppByTrackingID(mock.Anything, mock.Anything).Return(database.App{}, errors.New("app not found")).Once()
			},
			expectedErr: errors.New("app not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup(tc.trackingID)
			err := suite.service.ValidateAppAccess(suite.ctx, tc.userID, tc.trackingID)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
