package types

import "github.com/google/uuid"

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
