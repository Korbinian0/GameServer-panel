-- Schema für Gateway-Datenbank

CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  role TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE roles (
  name TEXT PRIMARY KEY,
  permissions TEXT[] NOT NULL
);

CREATE TABLE nodes (
  id TEXT PRIMARY KEY,
  platform TEXT NOT NULL,
  hostname TEXT NOT NULL,
  ip_address TEXT NOT NULL,
  capabilities TEXT[] NOT NULL,
  version TEXT NOT NULL,
  last_seen TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE server_instances (
  id TEXT PRIMARY KEY,
  node_id TEXT REFERENCES nodes(id) ON DELETE CASCADE,
  game_type TEXT NOT NULL,
  status TEXT NOT NULL,
  pid TEXT,
  uptime_seconds BIGINT DEFAULT 0,
  player_count INT DEFAULT 0,
  max_players INT DEFAULT 0
);
