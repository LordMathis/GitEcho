# GitEcho

![Go Tests](https://github.com/LordMathis/GitEcho/actions/workflows/go.yml/badge.svg)

:warning: **Work in Progress**

:warning: **Expect Breaking Changes**


GitEcho is a backup tool for git repositories. It uses [rclone](https://github.com/rclone/rclone) as a storage backend.

## Usage

```
  -f <path> Path to the config file 
  -h Print help and exit
  -r <repository_name> <storage_name> <local_path> Restore repository from storage backup to local path
```

## Installation

Gitecho already includes rclone as a dependency. You don't need to install it separately.

**As standalone binary**

```
go build -o gitecho ./cmd
./gitecho -f ./config.yaml
```

**With docker**

```
docker pull ghcr.io/lordmathis/gitecho
docker run gitecho
```

or build your own image

```
docker build -t gitecho .
docker run gitecho
```

## Configuration

GitEcho is configured using yaml config file. By default GitEcho looks for `config.yaml` file in the working directory. You can specify path to config file using the `-f` option 

**Full example:**

```yaml
data_path: /data  # Path where GitEcho stores git repositories
rclone_config_path: "./rclone.conf"  # Path to rclone config
repositories:  # Kist of repositories to backup
  - name: test-repo  # Unique name of the repository
    remote_url: "https://github.com/LordMathis/GitEcho"  # Remote git url, either https or ssh
    schedule:  "*/1 * * * *"  # Backup schedule, supports single number (minutes) or cron syntax
    storages:  # Dictionary of storages to backup to
      test-storage:  # Storage name
        remote_name: minio  # Remote name from rclone config
        remote_path: gitecho/test-repo  # Backup target path
    webhook:
      vendor: github  # One of "github", "gitea", "gitlab"
      secret: "test-secret"  # Webhook secret
      events: ["push", "create", "pull_request", "release"]  # Webhook events
    credentials:  # Credentials for git remote. Can be ommited if the repo is public
      username: gitecho  # git username
      password: gitecho  # git password
      key_path: /ssh/id_ed25519.pub  # Path to ssh key for authentication
```

For rclone config refer to the official [rclone documentation](https://rclone.org/docs/)

### Storage Config

A repository can have multiple storages (backup targets). Each storage must have a unique name (key). Storage configuration options must match options from rclone config.

**Minio example:**

rclone.conf
```conf
[minio]
env_auth = false
type = s3
provider = Minio
access_key_id = gitecho
secret_access_key = gitechokey
region = us-east-1
endpoint = http://minio:9000
location_constraint = 
server_side_encryption = 
```

config.yaml
```yaml
    storages:
      test-storage:
        remote_name: minio  # Matches [minio] from rclone.conf
        remote_path: gitecho/test-repo  # <bucket_name>/<path>
```

### Webhook config

GitEcho supports receiving webhooks from GitHub, Gitea and GitLab. Only a limited number of events are supported.

Vendor | Supported events
-------|---
[GitHub](https://docs.github.com/en/webhooks/webhook-events-and-payloads) | "create", "push", "pull_request", "release"  
[Gitea](https://docs.gitea.com/usage/webhooks) | "create", "push", "pull_request", "release"  
[GitLab](https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html) | "Push Hook", "Tag Push Hook", "Merge Request Hook"  

The webhook server runs on port 8080. The url is "/api/v1/webhooks/<repository_name>"

### Credentials

In case you GitEcho needs to authenticate to access the repository, you can specify either username and password or path the private ssh key. If the repository is public, you can omit this section

## Contributing

Issues and pull requests are always welcome




