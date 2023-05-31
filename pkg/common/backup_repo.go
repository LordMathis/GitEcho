package common

import (
	"os"

	gitclient "github.com/LordMathis/GitEcho/pkg/git"
	"github.com/go-git/go-git/v5"
)

type S3Credentials struct {
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
}

type BackupRepo struct {
	Name          string
	SrcRepo       *git.Repository
	PullInterval  int
	S3url         string
	S3bucket      string
	S3credentials *S3Credentials
	LocalPath     string
}

func NewS3Credentials(AWS_ACCESS_KEY_ID string, AWS_SECRET_ACCESS_KEY string) *S3Credentials {

	if AWS_ACCESS_KEY_ID == "" {
		AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
	}

	if AWS_SECRET_ACCESS_KEY == "" {
		AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	return &S3Credentials{
		AWS_ACCESS_KEY_ID:     AWS_ACCESS_KEY_ID,
		AWS_SECRET_ACCESS_KEY: AWS_SECRET_ACCESS_KEY,
	}
}

// NewBackupRepo creates a new BackupRepo instance
// srcRepoURL: the URL of the source repo to backup
// pullInterval: the interval (in seconds) between each pull operation
// s3bucket: the name of the S3 bucket to store the backups
// localPath: the local path where the backups will be stored
func NewBackupRepo(Name string, srcRepoURL string, pullInterval int, s3url string, s3bucket string, localPath string) (*BackupRepo, error) {

	// Clone the source repo to the local path
	srcRepo, err := gitclient.Clone(&git.CloneOptions{
		URL: srcRepoURL,
	}, localPath)

	if err != nil {
		return nil, err
	}

	// Return a new BackupRepo instance with the cloned source repo, pull interval, S3 bucket, and local path
	return &BackupRepo{
		Name:         Name,
		SrcRepo:      srcRepo,
		PullInterval: pullInterval,
		S3url:        s3url,
		S3bucket:     s3bucket,
		LocalPath:    localPath,
	}, nil
}
