-- name: FindOrCreateUser :one
INSERT INTO users (
  email
) VALUES ( $1 )
ON CONFLICT ( email ) DO UPDATE
SET email = EXCLUDED.email
RETURNING id;

-- name: CreateEvent :exec
INSERT INTO events (
  visitor_id, tracking_id, event_type, url, referrer, country, browser, operating_system, details, timestamp
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP );
