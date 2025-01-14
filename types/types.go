package types

import (
	"time"

	"github.com/google/uuid"
)

type EventDetails map[string]interface{}

type TrackingData struct {
	VisitorID string;
	TrackingID uuid.UUID;
	Url string;
	Referrer string;
	Country string;
	Ua string;
	Details EventDetails;
}

type EventPayload struct {
	Tracking TrackingData;
	Type string
}

type GeoLocation struct {
	Country string;
	City string;
	Longitude float64;
	Latitude float64
}

type UserAgentDetails struct {
	Browser string;
	Device string;
	OperatingSystem string;
}

type ReferralStats struct {
	Referrer string `json:"referrer"`;
	VisitorCount int `json:"visitor_count"`
}

type PageStats struct {
	Path string `json:"path"`;
	VisitorCount int `json:"visitor_count"`
}

type BrowserStats struct {
	Browser string `json:"browser"`;
	Percentage int `json:"percentage"`
}

type CountryStats struct {
	Country string `json:"country"`;
	Percentage int `json:"percentage"`
}

type DeviceStats struct {
	Device string `json:"device"`;
	Percentage int `json:"percentage"`
}

type OSStats struct {
	OS string `json:"operating_system"`;
	Percentage int `json:"percentage"`
}

type RequestPayload struct {
	UserID uuid.UUID
	StartTime time.Time
	EndTime time.Time
}
