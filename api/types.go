package main

import (
	"time"

	"github.com/google/uuid"
)

type TrackingData struct {
	VisitorID  string                 `json:"visitorID"`
	TrackingID uuid.UUID              `json:"trackingID"`
	Url        string                 `json:"url"`
	Referrer   string                 `json:"referrer"`
	Country    string                 `json:"country"`
	Ua         string                 `json:"ua"`
	Details    map[string]interface{} `json:"details"`
}

type EventPayload struct {
	Tracking TrackingData `json:"tracking"`
	Type     string       `json:"type"`
}

type GeoLocation struct {
	Country   string
	City      string
	Longitude float64
	Latitude  float64
}

type UserAgentDetails struct {
	Browser         string
	Device          string
	OperatingSystem string
}

type App struct {
	Name       string    `json:"name"`
	TrackingID uuid.UUID `json:"trackingID"`
	CreatedAt  time.Time `json:"created_at"`
}

type ReferralStats struct {
	Referrer     string `json:"referrer"`
	VisitorCount int    `json:"visitor_count"`
}

type PageStats struct {
	Path         string `json:"path"`
	VisitorCount int    `json:"visitor_count"`
}

type BrowserStats struct {
	Browser    string `json:"browser"`
	Percentage int    `json:"percentage"`
}

type CountryStats struct {
	Country    string `json:"country"`
	Percentage int    `json:"percentage"`
}

type DeviceStats struct {
	Device     string `json:"device"`
	Percentage int    `json:"percentage"`
}

type OSStats struct {
	OS         string `json:"operating_system"`
	Percentage int    `json:"percentage"`
}

type RequestPayload struct {
	UserID     uuid.UUID
	TrackingID uuid.UUID
	StartDate  time.Time
	EndDate    time.Time
}

type APIResponse struct {
	Data       interface{} `json:"data,omitempty"`
	statusCode int         `json:"-"`
	Message    string      `json:"message"`
}

type APIStatus struct {
	statusCode int    `json:"-"`
	Message    string `json:"message"`
}

type AppResponse struct {
	Data App
	APIStatus
}

type ReferralResponse struct {
	Data ReferralStats
	APIStatus
}

type PageResponse struct {
	Data PageStats
	APIStatus
}

type BrowserResponse struct {
	Data BrowserStats
	APIStatus
}

type CountryResponse struct {
	Data CountryStats
	APIStatus
}

type DeviceResponse struct {
	Data DeviceStats
	APIStatus
}

type OSResponse struct {
	Data OSStats
	APIStatus
}

type CreateAppRequest struct {
	Name string `json:"name"`
}

func NewSuccessResponse(data interface{}, code int, message string) APIResponse {
	return APIResponse{
		Data:       data,
		statusCode: code,
		Message:    message,
	}
}

func NewErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		statusCode: code,
		Message:    message,
	}
}
