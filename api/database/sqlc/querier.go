// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CheckAppExists(ctx context.Context, arg CheckAppExistsParams) (App, error)
	CreateApp(ctx context.Context, arg CreateAppParams) (App, error)
	CreateEvent(ctx context.Context, arg CreateEventParams) error
	GetAppByTrackingID(ctx context.Context, trackingID uuid.UUID) (App, error)
	GetApps(ctx context.Context, userID uuid.UUID) ([]App, error)
	GetBrowsers(ctx context.Context, arg GetBrowsersParams) ([]GetBrowsersRow, error)
	GetCountries(ctx context.Context, arg GetCountriesParams) ([]GetCountriesRow, error)
	GetDevices(ctx context.Context, arg GetDevicesParams) ([]GetDevicesRow, error)
	GetOS(ctx context.Context, arg GetOSParams) ([]GetOSRow, error)
	GetOrCreateUser(ctx context.Context, email string) (uuid.UUID, error)
	GetPageViews(ctx context.Context, arg GetPageViewsParams) ([]GetPageViewsRow, error)
	GetPages(ctx context.Context, arg GetPagesParams) ([]GetPagesRow, error)
	GetReferrals(ctx context.Context, arg GetReferralsParams) ([]GetReferralsRow, error)
	GetVisitors(ctx context.Context, arg GetVisitorsParams) ([]GetVisitorsRow, error)
}

var _ Querier = (*Queries)(nil)
