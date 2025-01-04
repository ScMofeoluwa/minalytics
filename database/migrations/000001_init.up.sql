CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  tracking_id UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE events (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  tracking_id UUID NOT NULL,
  visitor_id VARCHAR(64) NOT NULL,
  event_type VARCHAR(50) NOT NULL,
  url VARCHAR(255),
  referrer VARCHAR(255),
  country VARCHAR(100) NOT NULL,
  browser VARCHAR(100) NOT NULL,
  device VARCHAR(100) NOT NULL,
  operating_system VARCHAR(100) NOT NULL,
  details JSON,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_events_users FOREIGN KEY (tracking_id) REFERENCES users(tracking_id) ON DELETE CASCADE
);

CREATE INDEX idx_visitor_id ON events(visitor_id);
