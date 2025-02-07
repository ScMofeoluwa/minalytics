CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE apps (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  tracking_id UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
  user_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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
  details JSONB,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_app FOREIGN KEY (tracking_id) REFERENCES apps(tracking_id) ON DELETE CASCADE
);

CREATE INDEX idx_visitor_id ON events(visitor_id);
