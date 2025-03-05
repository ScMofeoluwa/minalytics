package server

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AnalyticsService interface {
	SignIn(context.Context, string) (string, error)
	TrackEvent(context.Context, EventPayload) error
	CreateApp(context.Context, uuid.UUID, string) (*App, error)
	UpdateApp(context.Context, AppPayload) (*App, error)
	DeleteApp(context.Context, AppPayload) error
	GetApps(context.Context, uuid.UUID) ([]App, error)
	GetReferrals(context.Context, RequestPayload) ([]ReferralStats, error)
	GetPages(context.Context, RequestPayload) ([]PageStats, error)
	GetBrowsers(context.Context, RequestPayload) ([]BrowserStats, error)
	GetCountries(context.Context, RequestPayload) ([]CountryStats, error)
	GetDevices(context.Context, RequestPayload) ([]DeviceStats, error)
	GetOS(context.Context, RequestPayload) ([]OSStats, error)
	GetVisitors(context.Context, RequestPayload) ([]VisitorStats, error)
	GetPageViews(context.Context, RequestPayload) ([]PageViewStats, error)
	ValidateAppAccess(context.Context, uuid.UUID, uuid.UUID) error
	ResolveGeoLocation(string) (*GeoLocation, error)
	ParseUserAgent(string) *UserAgentDetails
}

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

type AppPayload struct {
	Name       string
	TrackingID uuid.UUID
	UserID     uuid.UUID
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

type VisitorStats struct {
	Time     string `json:"time"`
	Visitors int    `json:"visitors"`
}

type PageViewStats struct {
	Time  string `json:"time"`
	Views int    `json:"views"`
}

type RequestPayload struct {
	TrackingID uuid.UUID
	BucketSize string
	StartDate  sql.NullTime
	EndDate    sql.NullTime
}

type APIResponse struct {
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"-"`
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

type VisitorResponse struct {
	Data VisitorStats
	APIStatus
}

type PageViewResponse struct {
	Data PageViewStats
	APIStatus
}

type CreateAppRequest struct {
	Name string `json:"name" binding:"required"`
}

func NewSuccessResponse(data interface{}, code int, message string) APIResponse {
	return APIResponse{
		Data:       data,
		StatusCode: code,
		Message:    message,
	}
}

func NewErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		StatusCode: code,
		Message:    message,
	}
}
