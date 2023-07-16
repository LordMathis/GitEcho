# GitEcho

![Go Tests](https://github.com/LordMathis/GitEcho/actions/workflows/go.yml/badge.svg)

:warning: **Work in Progress**

:warning: **Expect Breaking Changes**


GitEcho is a backup tool for git repositories.

## Features

- Periodically pulls changes from the repositories
- Uploads the changes to any S3 compatible storage
- REST API for managing backup configurations

## Configuration

Configure GitEcho via environment variables

**Set up database**

GitEcho supports postgres and sqlite databases.

```env
DB_TYPE=sqlite3  # "sqlite3" or "postgres"
```

1. Sqlite settings

```env
DB_PATH=./sqlite.db
```

2. Postgres settings

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=gitecho
DB_PASSWORD=gitecho
DB_NAME=gitecho
```

**Set up data path**

GitEcho clones your repository to data path

```
GITECHO_DATA_PATH=./data
```

**Set up encryption key**

GitEcho must store your git and/or s3 credentials. In order not to store them in plain text in the database it encrypts them. You cen generate a key by running `gitecho -g`

Copy the generated key and put it in `GITECHO_ENCRYPTION_KEY` environment variable. You can also use your own 16, 24 or 32 byte key

**Customize port**

GitEcho runs by default on port 8080. You can override it with `GITECHO_PORT` environment variable

## Deployment

**As standalone binary**

```
go build -o gitecho ./cmd/server
./gitecho
```

**With docker**

```
docker build -t gitecho .
docker run -p 8080 gitecho
```

**With docker-compose**

Adapt the example `docker-compose.yaml` to your needs and launch it

```
docker-compose up -d
```

## Usage

Navigate to `http://localhost:8080`. Add your backup repository configuration. The repository will be backed up every `Pull Interval` minutes