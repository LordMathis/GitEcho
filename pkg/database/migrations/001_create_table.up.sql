-- +migrate Up
CREATE TABLE backup_repo (
    name TEXT PRIMARY KEY,
    remote_url TEXT,
    pull_interval INT,
    s3_url TEXT,
    s3_bucket TEXT,
    local_path TEXT
);