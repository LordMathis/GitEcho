package common

import (
	gitclient "github.com/LordMathis/GitEcho/pkg/git"
	"github.com/go-git/go-git/v5"
)

type BackupRepo struct {
	SrcRepo      *git.Repository
	PullInterval int
	S3bucket     string
	LocalPath    string
}

// NewBackupRepo creates a new BackupRepo instance
// srcRepoURL: the URL of the source repo to backup
// pullInterval: the interval (in seconds) between each pull operation
// s3bucket: the name of the S3 bucket to store the backups
// localPath: the local path where the backups will be stored
func NewBackupRepo(srcRepoURL string, pullInterval int, s3bucket string, localPath string) (*BackupRepo, error) {

	// Clone the source repo to the local path
	srcRepo, err := gitclient.Clone(&git.CloneOptions{
		URL: srcRepoURL,
	}, localPath)

	if err != nil {
		return nil, err
	}

	// Return a new BackupRepo instance with the cloned source repo, pull interval, S3 bucket, and local path
	return &BackupRepo{
		SrcRepo:      srcRepo,
		PullInterval: pullInterval,
		S3bucket:     s3bucket,
		LocalPath:    localPath,
	}, nil
}
