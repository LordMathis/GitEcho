# GitEcho

![Go Tests](https://github.com/LordMathis/GitEcho/actions/workflows/go.yml/badge.svg)

:warning: **Work in Progress**

:warning: **Expect Breaking Changes**


GitEcho is a backup tool for git repositories.

## Usage

```
  -f <path> Path to the config file 
  -g Generate encryption key and exit
  -h Print help and exit
  -r <repository_name> <storage_name> <local_path> Restore repository from storage backup to local path
```

## Configuration

GiEcho is configured using yaml config file. By default GitEcho looks for `config.yaml` file in the working directory. You can specify path to config file using the `-f` option 

```yaml
data_path: /data  # Path where GitEcho stores git repositories
repositories:  # Kist of repositories to backup
  - name: test-repo  # Unique name of the repository
    remote_url: "https://github.com/LordMathis/GitEcho"  # Remote git url, either https or ssh
    schedule:  "*/1 * * * *"  # Backup schedule, supports single number (minutes) or cron syntax
    storages:  # List of storage names to backup to
      - test-storage
    credentials:  # Credentials for git remote. Can be ommited if the repo is public
      username: gitecho  # git username
      password: gitecho  # git password
      key_path: /ssh/id_ed25519.pub  # Path to ssh key for authentication
storages:
  - name: test-storage  # Storage name
    type: s3  # Storage type, currently only s3 is supported
    config:  # Config for storage type
      endpoint:   "http://127.0.0.1:9000"  # S3 endpoint
      region: us-east-1  # AWS region
      access_key:  gitecho  
      secret_key:  gitechokey
      bucket_name: gitecho
      disable_ssl: true  # Set disable SSL to true if you are using local minio over http
      force_path_style: true  # Force s3 api path style
      encryption:  # Client side encryption settings
        enabled: true
        key: 12345678901234567890123456789012  # 32 byte encryption key

```


## Deployment

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