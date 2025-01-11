package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"

	"github.com/ScMofeoluwa/minalytics/types"
)

type AnalyticsHandler struct {
	service *AnalyticsService
}

func NewAnalyticsHandler(service *AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}
func (h *AnalyticsHandler) Home(ctx *gin.Context) APIResponse {
	tmpl, err := template.ParseFiles(("templates/index.html"))
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to load template")
	}

	err = tmpl.Execute(ctx.Writer, nil)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to render template")
	}
	return APIResponse{}
}
func (h *AnalyticsHandler) SignIn(ctx *gin.Context) {
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return ctx.Param("provider"), nil
	}
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (h *AnalyticsHandler) Callback(ctx *gin.Context) APIResponse{
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return ctx.Param("provider"), nil
	}

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "authentication failed")
	}

	token, err := h.service.SignIn(ctx, user.Email)
	if err != nil {
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
		return NewErrorResponse(http.StatusInternalServerError, "failed to resolve geolocation")
	}

	payload.Tracking.Country = geoLocation.Country
	if err := h.service.TrackEvent(ctx, payload); err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to track event")
	}

	return NewSuccessResponse(nil, http.StatusOK, "event tracked successfully")
}

func (h *AnalyticsHandler) GetReferrals(ctx *gin.Context) APIResponse {
	userID, exists := ctx.Get("userID")
	if !exists {
		return NewErrorResponse(http.StatusUnauthorized, "userID not found in context")
	}
	user := userID.(uuid.UUID)

	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	parsedTimes, err := parseDates(startTime, endTime)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if parsedTimes[0].After(parsedTimes[1]){
		return NewErrorResponse(http.StatusBadRequest, "startTime cannot be after endTime")
	}

	if parsedTimes[0].Equal(parsedTimes[1]){
		return NewErrorResponse(http.StatusBadRequest, "startTime and endTime cannot be the same")
	}

	payload := types.ReferralPayload{
		UserID:    user,
		StartTime: parsedTimes[0],
		EndTime:   parsedTimes[1].Add(24 * time.Hour),
	}

	stats, err := h.service.GetReferrals(ctx, payload)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to fetch referrals")
	}

	return NewSuccessResponse(stats, http.StatusOK, "stats fetched successfully")
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
