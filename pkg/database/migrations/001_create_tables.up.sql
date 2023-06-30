-- +migrate Up

-- Create the storage table
CREATE TABLE IF NOT EXISTS storage (
  id SERIAL PRIMARY KEY,
  type TEXT NOT NULL,
  data JSONB
);

-- Create the backup_repo table
CREATE TABLE IF NOT EXISTS backup_repo (
  name TEXT PRIMARY KEY,
  pull_interval INT NOT NULL,
  storage_id INT NOT NULL REFERENCES storage(id),
  local_path TEXT NOT NULL,
  remote_url TEXT NOT NULL.
  git_username TEXT,
  git_password TEXT,
  git_key_path TEXT
);