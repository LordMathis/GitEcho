package common

import "github.com/go-git/go-git/v5"

type RepositoryBackupConfig struct {
	src_repo      *git.Repository
	pull_interval int
	s3bucket      string
	local_path    string
}
