package types

import (
	"time"

	"github.com/google/uuid"
)

type EventDetails map[string]interface{}

type TrackingData struct {
	VisitorId string;
	TrackingId uuid.UUID;
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

type ReferralPayload struct {
	UserID uuid.UUID
	StartTime time.Time
	EndTime time.Time
}
