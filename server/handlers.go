package server

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	types "github.com/ScMofeoluwa/minalytics/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

type AnalyticsHandler struct {
	service types.AnalyticsService
	logger  *zap.Logger
}

func NewAnalyticsHandler(service types.AnalyticsService, logger *zap.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
		logger:  logger,
	}
}

// @Summary User Sign-In
// @Description Initiates OAuth authentication with the specified provider and returns a JWT token upon successful login.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param provider path string true "OAuth provider (e.g., google, github)"
// @Success 200 {string} string "JWT token"
// @Failure 400 {object} types.APIStatus "Invalid provider or missing provider"
// @Failure 500 {object} types.APIStatus "Internal server error"
// @Router /auth/{provider} [get]
func (h *AnalyticsHandler) SignIn(ctx *gin.Context) {
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return ctx.Param("provider"), nil
	}
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (h *AnalyticsHandler) Callback(ctx *gin.Context) types.APIResponse {
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return ctx.Param("provider"), nil
	}

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		h.logger.Error("authentication failed", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "authentication failed")
	}

	token, err := h.service.SignIn(ctx, user.Email)
	if err != nil {
		h.logger.Error("failed to sign in", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to sign in")
	}

	return types.NewSuccessResponse(token, http.StatusOK, "login successful")
}

// @Summary Track an event
// @Description Tracks an event based on encoded data
// @Tags Analytics
// @Accept json
// @Produce json
// @Param data query string true "Base64 encoded event data"
// @Success 200 {object} types.APIStatus "Event tracked successfully"
// @Failure 400 {object} types.APIStatus "Invalid base64 or JSON data"
// @Failure 500 {object} types.APIStatus "Failed to resolve geolocation or track event"
// @Router /analytics/track [get]
func (h *AnalyticsHandler) TrackEvent(ctx *gin.Context) types.APIResponse {
	encodedData := ctx.Query("data")
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return types.NewErrorResponse(http.StatusBadRequest, "invalid base64 data")
	}

	var payload types.EventPayload
	if err := json.Unmarshal(decodedData, &payload); err != nil {
		return types.NewErrorResponse(http.StatusBadRequest, "invalid JSON data")
	}

	geoLocation, err := h.service.ResolveGeoLocation(ctx.ClientIP())
	if err != nil {
		h.logger.Error("failed to resolve geolocation", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to resolve geolocation")
	}

	payload.Tracking.Country = geoLocation.Country
	if err := h.service.TrackEvent(ctx, payload); err != nil {
		h.logger.Error("failed to track event", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to track event")
	}

	return types.NewSuccessResponse(nil, http.StatusOK, "event tracked successfully")
}

// @Summary Create App
// @Description creates an app
// @Tags Apps
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body types.CreateAppRequest true "app name"
// @Success 200 {object} types.AppResponse "app created successfully"
// @Failure 400 {object} types.APIStatus "invalid request body"
// @Failure 401 {object} types.APIStatus "userID not found in context"
// @Failure 500 {object} types.APIStatus "failed to create app"
// @Router /apps [post]
func (h *AnalyticsHandler) CreateApp(ctx *gin.Context) types.APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	var req types.CreateAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, "invalid request body")
	}

	app, err := h.service.CreateApp(ctx, user, req.Name)
	if err != nil {
		h.logger.Error("failed to create app", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to create app")
	}

	return types.NewSuccessResponse(app, http.StatusOK, "app created successfully")
}

// @Summary Get Apps
// @Description Retrieves user apps
// @Tags Apps
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} types.AppResponse "apps fetched successfully"
// @Failure 401 {object} types.APIStatus "userID not found in context"
// @Failure 500 {object} types.APIStatus "failed to fetch apps"
// @Router /apps [get]
func (h *AnalyticsHandler) GetApps(ctx *gin.Context) types.APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	apps, err := h.service.GetApps(ctx, user)
	if err != nil {
		h.logger.Error("failed to fetch apps", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch apps")
	}

	return types.NewSuccessResponse(apps, http.StatusOK, "apps fetched successfully")
}

// @Summary Update App
// @Description Updates app by tracking ID
// @Tags Apps
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param trackingID path string true "Tracking ID of the app to delete"
// @Param request body types.CreateAppRequest true "app name"
// @Success 200 {object} types.AppResponse "app name successfully changed"
// @Failure 400 {object} types.APIStatus "trackingID is required"
// @Failure 400 {object} types.APIStatus "invalid request body"
// @Failure 401 {object} types.APIStatus "userID not found in context"
// @Failure 500 {object} types.APIStatus "failed to update app"
// @Router /apps/{trackingID} [patch]
func (h *AnalyticsHandler) UpdateApp(ctx *gin.Context) types.APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	trackingID_ := ctx.Param("trackingID")
	if trackingID_ == "" {
		return types.NewErrorResponse(http.StatusBadRequest, "trackingID is required")
	}

	trackingID := uuid.MustParse(trackingID_)

	var req types.CreateAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, "invalid request body")
	}

	app, err := h.service.UpdateApp(ctx, req.Name, trackingID, user)
	if err != nil {
		h.logger.Error("failed to update app", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to update app")
	}

	return types.NewSuccessResponse(app, http.StatusOK, "app successfully updated")
}

// @Summary Delete App
// @Description Updates app by tracking ID
// @Tags Apps
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param trackingID path string true "Tracking ID of the app to delete"
// @Success 204 {string} string "app successfully deleted"
// @Failure 400 {object} types.APIStatus "trackingID is required"
// @Failure 400 {object} types.APIStatus "invalid request body"
// @Failure 401 {object} types.APIStatus "userID not found in context"
// @Failure 500 {object} types.APIStatus "failed to delete app"
// @Router /apps/{trackingID} [delete]
func (h *AnalyticsHandler) DeleteApp(ctx *gin.Context) types.APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	trackingID_ := ctx.Param("trackingID")
	if trackingID_ == "" {
		return types.NewErrorResponse(http.StatusBadRequest, "trackingID is required")
	}

	trackingID := uuid.MustParse(trackingID_)

	if err := h.service.DeleteApp(ctx, trackingID, user); err != nil {
		h.logger.Error("failed to delete app", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to delete app")
	}

	return types.NewSuccessResponse(nil, http.StatusNoContent, "app successfully updated")
}

// @Summary Get Referrals
// @Description Retrieves referral stats
// @Tags Analytics
// @Accept  json
// @Produce  json
// @Param trackingID query string true "app tracking ID"
// @Param startDate query string false "start date"
// @Param endDate query string false "end date"
// @Security BearerAuth
// @Success 200 {object} types.ReferralResponse "stats fetched successfully"
// @Failure 400 {object} types.APIStatus "invalid request paramaters"
// @Failure 500 {object} types.APIStatus "failed to fetch referrals"
// @Router /analytics/referrals [get]
func (h *AnalyticsHandler) GetReferrals(ctx *gin.Context) types.APIResponse {
	trackingID_, exists := ctx.Get("trackingID")
	if !exists {
		h.logger.Warn("trackingID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "trackingID not found in context")
	}
	trackingID := trackingID_.(uuid.UUID)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	payload, err := createRequestPayload(trackingID, startDate, endDate)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetReferrals(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch referrals", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch referrals")
	}

	return types.NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

// @Summary Get Pages
// @Description Retrieves page stats
// @Tags Analytics
// @Accept  json
// @Produce  json
// @Param trackingID query string true "app tracking ID"
// @Param startDate query string false "start date"
// @Param endDate query string false "end date"
// @Security BearerAuth
// @Success 200 {object} types.PageResponse "stats fetched successfully"
// @Failure 400 {object} types.APIStatus "invalid request paramaters"
// @Failure 500 {object} types.APIStatus "failed to fetch pages"
// @Router /analytics/pages [get]
func (h *AnalyticsHandler) GetPages(ctx *gin.Context) types.APIResponse {
	trackingID_, exists := ctx.Get("trackingID")
	if !exists {
		h.logger.Warn("trackingID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "trackingID not found in context")
	}
	trackingID := trackingID_.(uuid.UUID)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	payload, err := createRequestPayload(trackingID, startDate, endDate)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetPages(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch pages", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch pages")
	}

	return types.NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

// @Summary Get Browsers
// @Description Retrieves browser stats
// @Tags Analytics
// @Accept  json
// @Produce  json
// @Param trackingID query string true "app tracking ID"
// @Param startDate query string false "start date"
// @Param endDate query string false "end date"
// @Security BearerAuth
// @Success 200 {object} types.BrowserResponse "stats fetched successfully"
// @Failure 400 {object} types.APIStatus "invalid request paramaters"
// @Failure 500 {object} types.APIStatus "failed to fetch browsers"
// @Router /analytics/browsers [get]
func (h *AnalyticsHandler) GetBrowsers(ctx *gin.Context) types.APIResponse {
	trackingID_, exists := ctx.Get("trackingID")
	if !exists {
		h.logger.Warn("trackingID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "trackingID not found in context")
	}
	trackingID := trackingID_.(uuid.UUID)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	payload, err := createRequestPayload(trackingID, startDate, endDate)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetBrowsers(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch browsers", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch browsers")
	}

	return types.NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

// @Summary Get Countries
// @Description Retrieves country stats
// @Tags Analytics
// @Accept  json
// @Produce  json
// @Param trackingID query string true "app tracking ID"
// @Param startDate query string false "start date"
// @Param endDate query string false "end date"
// @Security BearerAuth
// @Success 200 {object} types.CountryResponse "stats fetched successfully"
// @Failure 400 {object} types.APIStatus "invalid request paramaters"
// @Failure 500 {object} types.APIStatus "failed to fetch countries"
// @Router /analytics/countries [get]
func (h *AnalyticsHandler) GetCountries(ctx *gin.Context) types.APIResponse {
	trackingID_, exists := ctx.Get("trackingID")
	if !exists {
		h.logger.Warn("trackingID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "trackingID not found in context")
	}
	trackingID := trackingID_.(uuid.UUID)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	payload, err := createRequestPayload(trackingID, startDate, endDate)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetCountries(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch countries", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch countries")
	}

	return types.NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

// @Summary Get Devices
// @Description Retrieves device stats
// @Tags Analytics
// @Accept  json
// @Produce  json
// @Param trackingID query string true "app tracking ID"
// @Param startDate query string false "start date"
// @Param endDate query string false "end date"
// @Security BearerAuth
// @Success 200 {object} types.DeviceResponse "stats fetched successfully"
// @Failure 400 {object} types.APIStatus "invalid request paramaters"
// @Failure 500 {object} types.APIStatus "failed to fetch devices"
// @Router /analytics/devices [get]
func (h *AnalyticsHandler) GetDevices(ctx *gin.Context) types.APIResponse {
	trackingID_, exists := ctx.Get("trackingID")
	if !exists {
		h.logger.Warn("trackingID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "trackingID not found in context")
	}
	trackingID := trackingID_.(uuid.UUID)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	payload, err := createRequestPayload(trackingID, startDate, endDate)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetDevices(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch devices", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch devices")
	}

	return types.NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

// @Summary Get OS
// @Description Retrieves operating system stats
// @Tags Analytics
// @Accept  json
// @Produce  json
// @Param trackingID query string true "app tracking ID"
// @Param startDate query string false "start date"
// @Param endDate query string false "end date"
// @Security BearerAuth
// @Success 200 {object} types.OSResponse "stats fetched successfully"
// @Failure 400 {object} types.APIStatus "invalid request paramaters"
// @Failure 500 {object} types.APIStatus "failed to fetch operating systems"
// @Router /analytics/os [get]
func (h *AnalyticsHandler) GetOS(ctx *gin.Context) types.APIResponse {
	trackingID_, exists := ctx.Get("trackingID")
	if !exists {
		h.logger.Warn("trackingID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "trackingID not found in context")
	}
	trackingID := trackingID_.(uuid.UUID)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	payload, err := createRequestPayload(trackingID, startDate, endDate)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetOS(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch operating systems", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch operating systems")
	}

	return types.NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

// @Summary Get Visitors
// @Description Retrieves visitors
// @Tags Analytics
// @Accept  json
// @Produce  json
// @Param trackingID query string true "app tracking ID"
// @Param startDate query string false "start date"
// @Param endDate query string false "end date"
// @Security BearerAuth
// @Success 200 {object} types.VisitorResponse "stats fetched successfully"
// @Failure 400 {object} types.APIStatus "invalid request paramaters"
// @Failure 500 {object} types.APIStatus "failed to fetch visitors"
// @Router /analytics/visitors [get]
func (h *AnalyticsHandler) GetVisitors(ctx *gin.Context) types.APIResponse {
	trackingID_, exists := ctx.Get("trackingID")
	if !exists {
		h.logger.Warn("trackingID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "trackingID not found in context")
	}
	trackingID := trackingID_.(uuid.UUID)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	payload, err := createRequestPayload(trackingID, startDate, endDate)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetVisitors(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch visitors", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch visitors")
	}

	return types.NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

// @Summary Get PageViews
// @Description Retrieves page views
// @Tags Analytics
// @Accept  json
// @Produce  json
// @Param trackingID query string true "app tracking ID"
// @Param startDate query string false "start date"
// @Param endDate query string false "end date"
// @Security BearerAuth
// @Success 200 {object} types.PageViewResponse "stats fetched successfully"
// @Failure 400 {object} types.APIStatus "invalid request paramaters"
// @Failure 500 {object} types.APIStatus "failed to fetch visitors"
// @Router /analytics/pageviews [get]
func (h *AnalyticsHandler) GetPageViews(ctx *gin.Context) types.APIResponse {
	trackingID_, exists := ctx.Get("trackingID")
	if !exists {
		h.logger.Warn("trackingID not found in context")
		return types.NewErrorResponse(http.StatusUnauthorized, "trackingID not found in context")
	}
	trackingID := trackingID_.(uuid.UUID)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	payload, err := createRequestPayload(trackingID, startDate, endDate)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return types.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetPageViews(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch page views", zap.Error(err))
		return types.NewErrorResponse(http.StatusInternalServerError, "failed to fetch page views")
	}

	return types.NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

func createRequestPayload(trackingID uuid.UUID, startDateStr, endDateStr string) (types.RequestPayload, error) {
	if (startDateStr == "" && endDateStr != "") || (startDateStr != "" && endDateStr == "") {
		return types.RequestPayload{}, fmt.Errorf("either specify both startDate and endDate, or specify neither")
	}

	if startDateStr == "" && endDateStr == "" {
		return types.RequestPayload{
			TrackingID: trackingID,
			BucketSize: "1 hour",
			StartDate:  sql.NullTime{},
			EndDate:    sql.NullTime{},
		}, nil
	}
	parsedTimes, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		return types.RequestPayload{}, err
	}

	startDate, endDate := parsedTimes[0], parsedTimes[1]

	if startDate.After(endDate) {
		return types.RequestPayload{}, fmt.Errorf("startDate cannot be after endDate")
	}

	if startDate.Equal(endDate) {
		return types.RequestPayload{}, fmt.Errorf("startDate and endDate cannot be the same")
	}

	return types.RequestPayload{
		TrackingID: trackingID,
		BucketSize: "1 day",
		StartDate:  sql.NullTime{Time: startDate, Valid: true},
		EndDate:    sql.NullTime{Time: endDate.Add(24 * time.Hour), Valid: true},
	}, nil
}

func parseDates(dateStrings ...string) ([]time.Time, error) {
	const layout = "2006-01-02"
	parsedTimes := make([]time.Time, 0, len(dateStrings))

	for _, dateString := range dateStrings {
		parsedTime, err := time.Parse(layout, dateString)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for %q, expect format: YY-MM-DD", dateString)
		}
		parsedTimes = append(parsedTimes, parsedTime)
	}

	return parsedTimes, nil
}
