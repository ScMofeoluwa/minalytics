package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ScMofeoluwa/minalytics/mocks"
	types "github.com/ScMofeoluwa/minalytics/shared"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HandlerSuite struct {
	suite.Suite
	mockService *mocks.AnalyticsService
	handler     *AnalyticsHandler
	logger      *zap.Logger
}

func createGinContext(req *http.Request, rr *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	return ctx
}

func (suite *HandlerSuite) SetupSuite() {
	suite.mockService = mocks.NewAnalyticsService(suite.T())
	suite.logger = zap.NewExample()
	suite.handler = NewAnalyticsHandler(suite.mockService, suite.logger)
}

func (suite *HandlerSuite) TestTrackEvent() {
	validPayload := types.EventPayload{
		Type: "pageview",
		Tracking: types.TrackingData{
			VisitorID:  uuid.NewString(),
			TrackingID: uuid.New(),
			Url:        faker.URL(),
			Referrer:   faker.URL(),
			Country:    faker.GetCountryInfo().Name,
			Ua:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
			Details:    map[string]interface{}{},
		},
	}

	payloadBytes, _ := json.Marshal(validPayload)
	encodedValidPayload := base64.StdEncoding.EncodeToString(payloadBytes)

	invalidJSON := []byte(`{"invalid": json`)
	encodedInvalidJSON := base64.StdEncoding.EncodeToString(invalidJSON)

	testCases := []struct {
		name       string
		mockSetup  func()
		query      string
		statusCode int
	}{
		{
			name:       "invalid base64 data",
			mockSetup:  func() {},
			query:      "invalid-base64",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "invalid JSON data",
			mockSetup:  func() {},
			query:      encodedInvalidJSON,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "failed to resolve geolocation",
			mockSetup: func() {
				suite.mockService.EXPECT().ResolveGeoLocation(mock.Anything).Return(&types.GeoLocation{}, fmt.Errorf("geolocation error")).Once()
			},
			query:      encodedValidPayload,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "failed to track event",
			mockSetup: func() {
				suite.mockService.EXPECT().ResolveGeoLocation(mock.Anything).Return(&types.GeoLocation{Country: "USA"}, nil).Once()
				suite.mockService.EXPECT().TrackEvent(mock.Anything, mock.Anything).Return(fmt.Errorf("failed to track event")).Once()
			},
			query:      encodedValidPayload,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "event successfully tracked",
			mockSetup: func() {
				suite.mockService.EXPECT().ResolveGeoLocation(mock.Anything).Return(&types.GeoLocation{Country: "USA"}, nil).Once()
				suite.mockService.EXPECT().TrackEvent(mock.Anything, mock.Anything).Return(nil).Once()
			},
			query:      encodedValidPayload,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/analytics/track?data="+tc.query, nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := createGinContext(req, rr)
			handlerFunc := WrapHandler(suite.handler.TrackEvent)
			handlerFunc(ctx)

			suite.Equal(tc.statusCode, rr.Code)
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *HandlerSuite) TestCreateApp() {
	testCases := []struct {
		name       string
		mockSetup  func()
		req        types.CreateAppRequest
		statusCode int
	}{
		{
			name:       "userID not found in context",
			mockSetup:  func() {},
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "invalid request body",
			mockSetup:  func() {},
			req:        types.CreateAppRequest{},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "failed to create app",
			mockSetup: func() {
				suite.mockService.EXPECT().CreateApp(mock.Anything, mock.Anything, mock.Anything).Return(&types.App{}, fmt.Errorf("failed to create app")).Once()
			},
			req: types.CreateAppRequest{
				Name: "test app",
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "app created successfully",
			mockSetup: func() {
				suite.mockService.EXPECT().CreateApp(mock.Anything, mock.Anything, mock.Anything).Return(&types.App{}, nil).Once()
			},
			req: types.CreateAppRequest{
				Name: "test app",
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()
			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(tc.req)
			suite.NoError(err)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/apps", b)
			req.Header.Add("Content-Type", "application/json")

			ctx := createGinContext(req, rr)
			if tc.statusCode != http.StatusUnauthorized {
				ctx.Set("userID", uuid.New())
			}

			handlerFunc := WrapHandler(suite.handler.CreateApp)
			handlerFunc(ctx)

			suite.Equal(tc.statusCode, rr.Code)
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *HandlerSuite) TestGetApps() {
	testCases := []struct {
		name       string
		mockSetup  func()
		statusCode int
	}{
		{
			name:       "userID not found in context",
			mockSetup:  func() {},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "failed to fetch apps",
			mockSetup: func() {
				suite.mockService.EXPECT().GetApps(mock.Anything, mock.Anything).Return([]types.App{}, fmt.Errorf("failed to fetch apps")).Once()
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "apps fetched successfully",
			mockSetup: func() {
				suite.mockService.EXPECT().GetApps(mock.Anything, mock.Anything).Return([]types.App{}, nil).Once()
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/apps", nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := createGinContext(req, rr)
			if tc.statusCode != http.StatusUnauthorized {
				ctx.Set("userID", uuid.New())
			}

			handlerFunc := WrapHandler(suite.handler.GetApps)
			handlerFunc(ctx)

			suite.Equal(tc.statusCode, rr.Code)
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *HandlerSuite) TestUpdateApp() {
	trackingID := uuid.New()
	testCases := []struct {
		name       string
		mockSetup  func()
		req        types.CreateAppRequest
		statusCode int
	}{
		{
			name:       "userID not found in context",
			mockSetup:  func() {},
			req:        types.CreateAppRequest{Name: "Updated App"},
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "invalid request body",
			mockSetup:  func() {},
			req:        types.CreateAppRequest{},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "failed to update app",
			mockSetup: func() {
				suite.mockService.EXPECT().UpdateApp(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to update app")).Once()
			},
			req:        types.CreateAppRequest{Name: "Updated App"},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "app successfully updated",
			mockSetup: func() {
				suite.mockService.EXPECT().UpdateApp(mock.Anything, mock.Anything).Return(&types.App{}, nil).Once()
			},
			req:        types.CreateAppRequest{Name: "Updated App"},
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(tc.req)
			suite.NoError(err)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPatch, "/apps/"+trackingID.String(), b)
			req.Header.Add("Content-Type", "application/json")

			ctx := createGinContext(req, rr)
			if tc.statusCode != http.StatusUnauthorized {
				ctx.Set("userID", uuid.New())
			}
			ctx.Params = gin.Params{{Key: "trackingID", Value: trackingID.String()}}

			handlerFunc := WrapHandler(suite.handler.UpdateApp)
			handlerFunc(ctx)

			suite.Equal(tc.statusCode, rr.Code)
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *HandlerSuite) TestDeleteApp() {
	trackingID := uuid.New()
	testCases := []struct {
		name       string
		mockSetup  func()
		statusCode int
	}{
		{
			name:       "userID not found in context",
			mockSetup:  func() {},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "failed to delete app",
			mockSetup: func() {
				suite.mockService.EXPECT().DeleteApp(mock.Anything, mock.Anything).Return(fmt.Errorf("failed to delete app")).Once()
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "app deleted successfully",
			mockSetup: func() {
				suite.mockService.EXPECT().DeleteApp(mock.Anything, mock.Anything).Return(nil).Once()
			},
			statusCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockSetup()

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodDelete, "/apps/"+trackingID.String(), nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := createGinContext(req, rr)
			if tc.statusCode != http.StatusUnauthorized {
				ctx.Set("userID", uuid.New())
			}
			ctx.Params = gin.Params{{Key: "trackingID", Value: trackingID.String()}}

			handlerFunc := WrapHandler(suite.handler.DeleteApp)
			handlerFunc(ctx)

			suite.Equal(tc.statusCode, rr.Code)
			suite.mockService.AssertExpectations(suite.T())
		})
	}
}

func (suite *HandlerSuite) TestGetAnalyticsEndpoints() {
	type analyticsTest struct {
		name       string
		startDate  string
		endDate    string
		mockSetup  func()
		statusCode int
	}

	testEndpoint := func(endpoint string, funcName string, handler func(*gin.Context) types.APIResponse, mockResult interface{}) {
		testCases := []analyticsTest{
			{
				name:       "trackingID not found in context",
				startDate:  "",
				endDate:    "",
				mockSetup:  func() {},
				statusCode: http.StatusUnauthorized,
			},
			{
				name:       "only start date provided",
				startDate:  "2025-03-01",
				mockSetup:  func() {},
				statusCode: http.StatusBadRequest,
			},
			{
				name:       "only end date provided",
				endDate:    "2025-03-06",
				mockSetup:  func() {},
				statusCode: http.StatusBadRequest,
			},
			{
				name:       "invalid date format",
				startDate:  "2025/03/01",
				endDate:    "2025/03/06",
				mockSetup:  func() {},
				statusCode: http.StatusBadRequest,
			},
			{
				name:       "start date after end date",
				startDate:  "2025-03-06",
				endDate:    "2025-03-01",
				mockSetup:  func() {},
				statusCode: http.StatusBadRequest,
			},
			{
				name:       "same dates",
				startDate:  "2025-03-01",
				endDate:    "2025-03-01",
				mockSetup:  func() {},
				statusCode: http.StatusBadRequest,
			},
			{
				name:      "failed to fetch" + endpoint,
				startDate: "2025-03-01",
				endDate:   "2025-03-06",
				mockSetup: func() {
					suite.mockService.On(funcName, mock.Anything, mock.Anything).Return(nil, errors.New("database error")).Once()
				},
				statusCode: http.StatusInternalServerError,
			},
			{
				name:      "successful stats retrieval",
				startDate: "2025-03-01",
				endDate:   "2025-03-06",
				mockSetup: func() {
					suite.mockService.On(funcName, mock.Anything, mock.Anything).Return(mockResult, nil).Once()
				},
				statusCode: http.StatusOK,
			},
		}

		for _, tc := range testCases {
			suite.Run(endpoint+" - "+tc.name, func() {
				tc.mockSetup()

				rr := httptest.NewRecorder()
				url := "/analytics/" + endpoint

				// Add query parameters if provided
				if tc.startDate != "" || tc.endDate != "" {
					url += "?"
					if tc.startDate != "" {
						url += "startDate=" + tc.startDate
					}
					if tc.endDate != "" {
						if tc.startDate != "" {
							url += "&"
						}
						url += "endDate=" + tc.endDate
					}
				}

				req := httptest.NewRequest(http.MethodGet, url, nil)
				req.Header.Add("Content-Type", "application/json")

				ctx := createGinContext(req, rr)
				if tc.statusCode != http.StatusUnauthorized {
					ctx.Set("trackingID", uuid.New())
				}

				handlerFunc := WrapHandler(handler)
				handlerFunc(ctx)

				suite.Equal(tc.statusCode, rr.Code)
				suite.mockService.AssertExpectations(suite.T())
			})
		}
	}

	// Test all analytics endpoints
	testEndpoint("referrals", "GetReferrals", suite.handler.GetReferrals, []types.ReferralStats{})
	testEndpoint("pages", "GetPages", suite.handler.GetPages, []types.PageStats{})
	testEndpoint("browsers", "GetBrowsers", suite.handler.GetBrowsers, []types.BrowserStats{})
	testEndpoint("countries", "GetCountries", suite.handler.GetCountries, []types.CountryStats{})
	testEndpoint("devices", "GetDevices", suite.handler.GetDevices, []types.DeviceStats{})
	testEndpoint("os", "GetOS", suite.handler.GetOS, []types.OSStats{})
	testEndpoint("visitors", "GetVisitors", suite.handler.GetVisitors, []types.VisitorStats{})
	testEndpoint("pageviews", "GetPageViews", suite.handler.GetPageViews, []types.PageViewStats{})
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
