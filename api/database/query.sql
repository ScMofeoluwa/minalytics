-- name: GetOrCreateUser :one
INSERT INTO users (
  email
) VALUES ( $1 )
ON CONFLICT ( email ) DO UPDATE
SET email = EXCLUDED.email
RETURNING id;

-- name: GetUserByTrackingID :one
SELECT * FROM users WHERE users.tracking_id = $1;

-- name: GetUserTrackingID :one
SELECT tracking_id FROM users WHERE users.id = $1;

-- name: GetReferrals :many
SELECT referrer, COUNT(DISTINCT visitor_id) AS visitor_count
FROM users u JOIN events e ON u.tracking_id = e.tracking_id
WHERE referrer IS NOT NULL AND u.id = $1 AND timestamp BETWEEN $2 AND $3 
GROUP BY referrer
ORDER BY visitor_count DESC;

-- name: GetPages :many
SELECT url, COUNT(DISTINCT visitor_id) AS visitor_count
FROM users u JOIN events e ON u.tracking_id = e.tracking_id
WHERE url IS NOT NULL AND e.event_type = 'pageview' AND u.id = $1 AND timestamp BETWEEN $2 AND $3 
GROUP BY url
ORDER BY visitor_count DESC;

-- name: GetCountries :many
SELECT country, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM users u JOIN events e ON u.tracking_id = e.tracking_id
WHERE u.id = $1 AND timestamp BETWEEN $2 AND $3 
GROUP BY country
ORDER BY percentage DESC;

-- name: GetBrowsers :many
SELECT browser, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM users u JOIN events e ON u.tracking_id = e.tracking_id
WHERE u.id = $1 AND timestamp BETWEEN $2 AND $3 
GROUP BY browser
ORDER BY percentage DESC;

-- name: GetDevices :many
SELECT device, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM users u JOIN events e ON u.tracking_id = e.tracking_id
WHERE u.id = $1 AND timestamp BETWEEN $2 AND $3 
GROUP BY device
ORDER BY percentage DESC;

-- name: GetOS :many
SELECT operating_system, ROUND((COUNT(DISTINCT visitor_id) * 100.0) / SUM(COUNT(DISTINCT visitor_id)) OVER (), 0) as percentage
FROM users u JOIN events e ON u.tracking_id = e.tracking_id
WHERE u.id = $1 AND timestamp BETWEEN $2 AND $3 
GROUP BY operating_system
ORDER BY percentage DESC;

-- name: CreateEvent :exec
INSERT INTO events (
  visitor_id, tracking_id, event_type, url, referrer, country, browser, device, operating_system, details
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 );
