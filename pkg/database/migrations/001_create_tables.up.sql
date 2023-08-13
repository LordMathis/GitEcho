-- +migrate Up

-- Create the storage table
CREATE TABLE IF NOT EXISTS storage (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  type TEXT NOT NULL,
  data JSONB
);

-- Create the backup_repo table
CREATE TABLE IF NOT EXISTS backup_repo (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  schedule INT NOT NULL,
  local_path TEXT NOT NULL,
  remote_url TEXT NOT NULL,
  git_username TEXT,
  git_password TEXT,
  git_key_path TEXT
);

-- Create the backup_repo_storage table
CREATE TABLE IF NOT EXISTS backup_repo_storage (
  id SERIAL PRIMARY KEY,
  backup_repo_name TEXT REFERENCES backup_repo(name),
  storage_name TEXT REFERENCES storage(name),
  UNIQUE (backup_repo_name, storage_name)
);
