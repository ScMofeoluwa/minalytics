package types

import (
	"time"

	"github.com/google/uuid"
)

type EventDetails map[string]interface{}

type TrackingData struct {
	VisitorID  string       `json:"visitorID"`
	TrackingID uuid.UUID    `json:"trackingID"`
	Url        string       `json:"url"`
	Referrer   string       `json:"referrer"`
	Country    string       `json:"country"`
	Ua         string       `json:"ua"`
	Details    EventDetails `json:"details"`
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
