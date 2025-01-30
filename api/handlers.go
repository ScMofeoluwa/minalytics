package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"

	"github.com/ScMofeoluwa/minalytics/types"
)

type AnalyticsHandler struct {
	service *AnalyticsService
	logger  *zap.Logger
}

func NewAnalyticsHandler(service *AnalyticsService, logger *zap.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AnalyticsHandler) SignIn(ctx *gin.Context) {
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return ctx.Param("provider"), nil
	}
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (h *AnalyticsHandler) Callback(ctx *gin.Context) APIResponse {
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return ctx.Param("provider"), nil
	}

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		h.logger.Error("authentication failed", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "authentication failed")
	}

	token, err := h.service.SignIn(ctx, user.Email)
	if err != nil {
		h.logger.Error("failed to sign in", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to sign in")
	}

	return NewSuccessResponse(token, http.StatusOK, "login successful")
}

func (h *AnalyticsHandler) TrackEvent(ctx *gin.Context) APIResponse {
	encodedData := ctx.Query("data")
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return NewErrorResponse(http.StatusBadRequest, "invalid base64 data")
	}

	var payload types.EventPayload
	if err := json.Unmarshal(decodedData, &payload); err != nil {
		return NewErrorResponse(http.StatusBadRequest, "invalid JSON data")
	}

	geoLocation, err := h.service.ResolveGeoLocation(ctx.ClientIP())
	if err != nil {
		h.logger.Error("failed to resolve geolocation", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to resolve geolocation")
	}

	payload.Tracking.Country = geoLocation.Country
	if err := h.service.TrackEvent(ctx, payload); err != nil {
		h.logger.Error("failed to track event", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to track event")
	}

	return NewSuccessResponse(nil, http.StatusOK, "event tracked successfully")
}

func (h *AnalyticsHandler) GetTrackingID(ctx *gin.Context) APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}

	user := userID.(uuid.UUID)
	trackingID, err := h.service.GetTrackingID(ctx, user)
	if err != nil {
		h.logger.Error("failed to get trackingID", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to fetch trackingID")
	}

	return NewSuccessResponse(trackingID, http.StatusOK, "tracking ID fetched successfully")
}

func (h *AnalyticsHandler) GetReferrals(ctx *gin.Context) APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	payload, err := createRequestPayload(user, startTime, endTime)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetReferrals(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch referrals", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to fetch referrals")
	}

	return NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

func (h *AnalyticsHandler) GetPages(ctx *gin.Context) APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	payload, err := createRequestPayload(user, startTime, endTime)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetPages(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch pages", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to fetch pages")
	}

	return NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

func (h *AnalyticsHandler) GetBrowsers(ctx *gin.Context) APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	payload, err := createRequestPayload(user, startTime, endTime)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetBrowsers(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch browsers", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to fetch browsers")
	}

	return NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

func (h *AnalyticsHandler) GetCountries(ctx *gin.Context) APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	payload, err := createRequestPayload(user, startTime, endTime)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetCountries(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch countries", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to fetch countries")
	}

	return NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

func (h *AnalyticsHandler) GetDevices(ctx *gin.Context) APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	payload, err := createRequestPayload(user, startTime, endTime)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetDevices(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch devices", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to fetch devices")
	}

	return NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

func (h *AnalyticsHandler) GetOS(ctx *gin.Context) APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		h.logger.Warn("userID not found in context")
		return NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	payload, err := createRequestPayload(user, startTime, endTime)
	if err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	stats, err := h.service.GetOS(ctx, payload)
	if err != nil {
		h.logger.Error("failed to fetch operating systems", zap.Error(err))
		return NewErrorResponse(http.StatusInternalServerError, "failed to fetch operating systems")
	}

	return NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
}

func createRequestPayload(userID uuid.UUID, startTimeStr, endTimeStr string) (types.RequestPayload, error) {
	parsedTimes, err := parseDates(startTimeStr, endTimeStr)
	if err != nil {
		return types.RequestPayload{}, err
	}

	startTime, endTime := parsedTimes[0], parsedTimes[1]

	if startTime.After(endTime) {
		return types.RequestPayload{}, fmt.Errorf("startTime cannot be after endTime")
	}

	if startTime.Equal(endTime) {
		return types.RequestPayload{}, fmt.Errorf("startTime and endTime cannot be the same")
	}

	return types.RequestPayload{
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime.Add(24 * time.Hour),
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
