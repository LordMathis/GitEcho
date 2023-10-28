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

## Configuration

GiEcho is configured using yaml config file. By default GitEcho looks for `config.yaml` file in the working directory. You can specify path to config file using the `-f` option 

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
    credentials:  # Credentials for git remote. Can be ommited if the repo is public
      username: gitecho  # git username
      password: gitecho  # git password
      key_path: /ssh/id_ed25519.pub  # Path to ssh key for authentication
```

For rclone config refer to the official [rclone documentation](https://rclone.org/docs/)


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