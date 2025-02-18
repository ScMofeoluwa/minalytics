-- name: GetOrCreateUser :one
INSERT INTO users (
  email
) VALUES ( $1 )
ON CONFLICT ( email ) DO UPDATE
SET email = EXCLUDED.email
RETURNING id;

-- name: CreateApp :one
INSERT INTO apps (
  name, user_id
) VALUES ( $1, $2 )
RETURNING *;

-- name: CreateEvent :exec
INSERT INTO events (
  visitor_id, tracking_id, event_type, url, referrer, country, browser, device, operating_system, details
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 );

-- name: UpdateApp :one
UPDATE apps
SET name = $1
WHERE tracking_id = $2
RETURNING *;

-- name: DeleteApp :exec
DELETE FROM apps WHERE tracking_id = $1;

-- name: GetAppByTrackingID :one
SELECT * FROM apps WHERE tracking_id = $1;

-- name: CheckAppExists :one
SELECT * FROM apps WHERE user_id = $1 AND name = $2;

-- name: GetApps :many
SELECT * FROM apps WHERE user_id = $1;

-- name: GetVisitors :many
SELECT time_bucket($4, timestamp::timestamptz)::timestamptz AS time, COUNT(DISTINCT visitor_id) AS visitors
FROM events WHERE tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY time;

-- name: GetPageViews :many
SELECT time_bucket($4, timestamp::timestamptz)::timestamptz AS time, COUNT(url) AS views
FROM events WHERE tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY time;

-- name: GetReferrals :many
SELECT referrer, COUNT(DISTINCT visitor_id) AS visitor_count
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE referrer IS NOT NULL AND a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY referrer
ORDER BY visitor_count DESC;

-- name: GetPages :many
SELECT url, COUNT(DISTINCT visitor_id) AS visitor_count
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE url IS NOT NULL AND e.event_type = 'pageview' AND a.tracking_id = $1 AND 
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY url
ORDER BY visitor_count DESC;

-- name: GetCountries :many
SELECT country, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY country
ORDER BY percentage DESC;

-- name: GetBrowsers :many
SELECT browser, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY browser
ORDER BY percentage DESC;

-- name: GetDevices :many
SELECT device, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY device
ORDER BY percentage DESC;

-- name: GetOS :many
SELECT operating_system, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM apps a JOIN events e ON a.tracking_id = e.tracking_id
WHERE a.tracking_id = $1 AND
(
  ($2::timestamptz IS NULL AND $3::timestamptz IS NULL AND timestamp >= NOW() - INTERVAL '24 hours') OR
  (timestamp BETWEEN $2 AND $3)
)
GROUP BY operating_system
ORDER BY percentage DESC;

