// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	types "github.com/ScMofeoluwa/minalytics/types"
	"github.com/google/uuid"
)

type App struct {
	ID         uuid.UUID `json:"id"`
	TrackingID uuid.UUID `json:"tracking_id"`
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
}

type Event struct {
	ID              uuid.UUID          `json:"id"`
	TrackingID      uuid.UUID          `json:"tracking_id"`
	VisitorID       string             `json:"visitor_id"`
	EventType       string             `json:"event_type"`
	Url             *string            `json:"url"`
	Referrer        *string            `json:"referrer"`
	Country         string             `json:"country"`
	Browser         string             `json:"browser"`
	Device          string             `json:"device"`
	OperatingSystem string             `json:"operating_system"`
	Details         types.EventDetails `json:"details"`
	Timestamp       time.Time          `json:"timestamp"`
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
