// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createApp = `-- name: CreateApp :one
INSERT INTO apps (
  name, user_id
) VALUES ( $1, $2 )
RETURNING id, tracking_id, user_id, name, created_at
`

type CreateAppParams struct {
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateApp(ctx context.Context, arg CreateAppParams) (App, error) {
	row := q.db.QueryRow(ctx, createApp, arg.Name, arg.UserID)
	var i App
	err := row.Scan(
		&i.ID,
		&i.TrackingID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
	)
	return i, err
}

const createEvent = `-- name: CreateEvent :exec
INSERT INTO events (
  visitor_id, tracking_id, event_type, url, referrer, country, browser, device, operating_system, details
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 )
`

type CreateEventParams struct {
	VisitorID       string       `json:"visitor_id"`
	TrackingID      uuid.UUID    `json:"tracking_id"`
	EventType       string       `json:"event_type"`
	Url             *string      `json:"url"`
	Referrer        *string      `json:"referrer"`
	Country         string       `json:"country"`
	Browser         string       `json:"browser"`
	Device          string       `json:"device"`
	OperatingSystem string       `json:"operating_system"`
	Details         EventDetails `json:"details"`
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) error {
	_, err := q.db.Exec(ctx, createEvent,
		arg.VisitorID,
		arg.TrackingID,
		arg.EventType,
		arg.Url,
		arg.Referrer,
		arg.Country,
		arg.Browser,
		arg.Device,
		arg.OperatingSystem,
		arg.Details,
	)
	return err
}

const getAppByTrackingID = `-- name: GetAppByTrackingID :one
SELECT id, tracking_id, user_id, name, created_at FROM apps WHERE apps.tracking_id = $1
`

func (q *Queries) GetAppByTrackingID(ctx context.Context, trackingID uuid.UUID) (App, error) {
	row := q.db.QueryRow(ctx, getAppByTrackingID, trackingID)
	var i App
	err := row.Scan(
		&i.ID,
		&i.TrackingID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
	)
	return i, err
}

const getApps = `-- name: GetApps :many
SELECT id, tracking_id, user_id, name, created_at FROM apps WHERE apps.user_id = $1
`

func (q *Queries) GetApps(ctx context.Context, userID uuid.UUID) ([]App, error) {
	rows, err := q.db.Query(ctx, getApps, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []App{}
	for rows.Next() {
		var i App
		if err := rows.Scan(
			&i.ID,
			&i.TrackingID,
			&i.UserID,
			&i.Name,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBrowsers = `-- name: GetBrowsers :many
SELECT browser, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY browser
ORDER BY percentage DESC
`

type GetBrowsersParams struct {
	TrackingID uuid.UUID    `json:"tracking_id"`
	Column2    sql.NullTime `json:"column_2"`
	Column3    sql.NullTime `json:"column_3"`
}

type GetBrowsersRow struct {
	Browser    string `json:"browser"`
	Percentage int    `json:"percentage"`
}

func (q *Queries) GetBrowsers(ctx context.Context, arg GetBrowsersParams) ([]GetBrowsersRow, error) {
	rows, err := q.db.Query(ctx, getBrowsers, arg.TrackingID, arg.Column2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetBrowsersRow{}
	for rows.Next() {
		var i GetBrowsersRow
		if err := rows.Scan(&i.Browser, &i.Percentage); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCountries = `-- name: GetCountries :many
SELECT country, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY country
ORDER BY percentage DESC
`

type GetCountriesParams struct {
	TrackingID uuid.UUID    `json:"tracking_id"`
	Column2    sql.NullTime `json:"column_2"`
	Column3    sql.NullTime `json:"column_3"`
}

type GetCountriesRow struct {
	Country    string `json:"country"`
	Percentage int    `json:"percentage"`
}

func (q *Queries) GetCountries(ctx context.Context, arg GetCountriesParams) ([]GetCountriesRow, error) {
	rows, err := q.db.Query(ctx, getCountries, arg.TrackingID, arg.Column2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCountriesRow{}
	for rows.Next() {
		var i GetCountriesRow
		if err := rows.Scan(&i.Country, &i.Percentage); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDevices = `-- name: GetDevices :many
SELECT device, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY device
ORDER BY percentage DESC
`

type GetDevicesParams struct {
	TrackingID uuid.UUID    `json:"tracking_id"`
	Column2    sql.NullTime `json:"column_2"`
	Column3    sql.NullTime `json:"column_3"`
}

type GetDevicesRow struct {
	Device     string `json:"device"`
	Percentage int    `json:"percentage"`
}

func (q *Queries) GetDevices(ctx context.Context, arg GetDevicesParams) ([]GetDevicesRow, error) {
	rows, err := q.db.Query(ctx, getDevices, arg.TrackingID, arg.Column2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetDevicesRow{}
	for rows.Next() {
		var i GetDevicesRow
		if err := rows.Scan(&i.Device, &i.Percentage); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOS = `-- name: GetOS :many
SELECT operating_system, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY operating_system
ORDER BY percentage DESC
`

type GetOSParams struct {
	TrackingID uuid.UUID    `json:"tracking_id"`
	Column2    sql.NullTime `json:"column_2"`
	Column3    sql.NullTime `json:"column_3"`
}

type GetOSRow struct {
	OperatingSystem string `json:"operating_system"`
	Percentage      int    `json:"percentage"`
}

func (q *Queries) GetOS(ctx context.Context, arg GetOSParams) ([]GetOSRow, error) {
	rows, err := q.db.Query(ctx, getOS, arg.TrackingID, arg.Column2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetOSRow{}
	for rows.Next() {
		var i GetOSRow
		if err := rows.Scan(&i.OperatingSystem, &i.Percentage); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrCreateUser = `-- name: GetOrCreateUser :one
INSERT INTO users (
  email
) VALUES ( $1 )
ON CONFLICT ( email ) DO UPDATE
SET email = EXCLUDED.email
RETURNING id
`

func (q *Queries) GetOrCreateUser(ctx context.Context, email string) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, getOrCreateUser, email)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getPageViews = `-- name: GetPageViews :many
SELECT time_bucket($4, timestamp::timestamptz)::timestamptz AS time, COUNT(url) AS views
FROM events WHERE tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY time
`

type GetPageViewsParams struct {
	TrackingID uuid.UUID    `json:"tracking_id"`
	Column2    sql.NullTime `json:"column_2"`
	Column3    sql.NullTime `json:"column_3"`
	TimeBucket interface{}  `json:"time_bucket"`
}

type GetPageViewsRow struct {
	Time  sql.NullTime `json:"time"`
	Views int64        `json:"views"`
}

func (q *Queries) GetPageViews(ctx context.Context, arg GetPageViewsParams) ([]GetPageViewsRow, error) {
	rows, err := q.db.Query(ctx, getPageViews,
		arg.TrackingID,
		arg.Column2,
		arg.Column3,
		arg.TimeBucket,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPageViewsRow{}
	for rows.Next() {
		var i GetPageViewsRow
		if err := rows.Scan(&i.Time, &i.Views); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPages = `-- name: GetPages :many
SELECT url, COUNT(DISTINCT visitor_id) AS visitor_count
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE url IS NOT NULL AND e.event_type = 'pageview' AND a.tracking_id = $1 AND 
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY url
ORDER BY visitor_count DESC
`

type GetPagesParams struct {
	TrackingID uuid.UUID    `json:"tracking_id"`
	Column2    sql.NullTime `json:"column_2"`
	Column3    sql.NullTime `json:"column_3"`
}

type GetPagesRow struct {
	Url          *string `json:"url"`
	VisitorCount int64   `json:"visitor_count"`
}

func (q *Queries) GetPages(ctx context.Context, arg GetPagesParams) ([]GetPagesRow, error) {
	rows, err := q.db.Query(ctx, getPages, arg.TrackingID, arg.Column2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPagesRow{}
	for rows.Next() {
		var i GetPagesRow
		if err := rows.Scan(&i.Url, &i.VisitorCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReferrals = `-- name: GetReferrals :many
SELECT referrer, COUNT(DISTINCT visitor_id) AS visitor_count
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE referrer IS NOT NULL AND a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY referrer
ORDER BY visitor_count DESC
`

type GetReferralsParams struct {
	TrackingID uuid.UUID    `json:"tracking_id"`
	Column2    sql.NullTime `json:"column_2"`
	Column3    sql.NullTime `json:"column_3"`
}

type GetReferralsRow struct {
	Referrer     *string `json:"referrer"`
	VisitorCount int64   `json:"visitor_count"`
}

func (q *Queries) GetReferrals(ctx context.Context, arg GetReferralsParams) ([]GetReferralsRow, error) {
	rows, err := q.db.Query(ctx, getReferrals, arg.TrackingID, arg.Column2, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetReferralsRow{}
	for rows.Next() {
		var i GetReferralsRow
		if err := rows.Scan(&i.Referrer, &i.VisitorCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVisitors = `-- name: GetVisitors :many
SELECT time_bucket($4, timestamp::timestamptz)::timestamptz AS time, COUNT(DISTINCT visitor_id) AS visitors
FROM events WHERE tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY time
`

type GetVisitorsParams struct {
	TrackingID uuid.UUID    `json:"tracking_id"`
	Column2    sql.NullTime `json:"column_2"`
	Column3    sql.NullTime `json:"column_3"`
	TimeBucket interface{}  `json:"time_bucket"`
}

type GetVisitorsRow struct {
	Time     sql.NullTime `json:"time"`
	Visitors int64        `json:"visitors"`
}

func (q *Queries) GetVisitors(ctx context.Context, arg GetVisitorsParams) ([]GetVisitorsRow, error) {
	rows, err := q.db.Query(ctx, getVisitors,
		arg.TrackingID,
		arg.Column2,
		arg.Column3,
		arg.TimeBucket,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetVisitorsRow{}
	for rows.Next() {
		var i GetVisitorsRow
		if err := rows.Scan(&i.Time, &i.Visitors); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
