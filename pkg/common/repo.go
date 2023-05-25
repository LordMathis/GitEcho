package common

import "github.com/go-git/go-git/v5"

type RepositoryBackupConfig struct {
	SrcRepo      *git.Repository
	PullInterval int
	S3bucket     string
	LocalPath    string
}
