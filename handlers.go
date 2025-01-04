package main

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *AnalyticsHandler) Home(ctx *gin.Context) {
	tmpl, err := template.ParseFiles(("templates/index.html"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
		return
	}

	err = tmpl.Execute(ctx.Writer, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
		return
	}
}
func (h *AnalyticsHandler) SignIn(ctx *gin.Context) {
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return ctx.Param("provider"), nil
	}
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (h *AnalyticsHandler) Callback(ctx *gin.Context) {
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return ctx.Param("provider"), nil
	}

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.SignIn(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "login successful",
		"accessToken": token,
	})
}

func (h *AnalyticsHandler) TrackEvent(ctx *gin.Context) {
	encodedData := ctx.Query("data")
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64 data"})
		return
	}

	var payload types.EventPayload
	if err := json.Unmarshal(decodedData, &payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON data"})
		return
	}

	remoteAddr := ctx.ClientIP()
	geoLocation, err := h.service.ResolveGeoLocation(remoteAddr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve geolocation"})
		return
	}

	payload.Tracking.Country = geoLocation.Country
	if err := h.service.TrackEvent(ctx, payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event tracked successfully"})
}
