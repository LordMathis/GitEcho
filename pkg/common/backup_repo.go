package common

import (
	"fmt"
	"os"
	"strings"

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
func NewBackupRepo(name string, srcRepoURL string, pullInterval int, s3url string, s3bucket string, localPath string) (*BackupRepo, error) {

	_, err := git.PlainOpen(localPath)
	if err == nil {
		// If repository exists, pull latest changes
		fmt.Printf("Pulling repository at %s\n", localPath)

		srcRepo, err := git.PlainOpen(localPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open repository: %v", err)
		}

		wt, err := srcRepo.Worktree()
		if err != nil {
			return nil, fmt.Errorf("failed to get worktree: %v", err)
		}

		err = wt.Pull(&git.PullOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, fmt.Errorf("failed to pull repository: %v", err)
		}

		// Extract repository name from remote URL if not provided
		if name == "" {
			urlParts := strings.Split(srcRepoURL, "/")
			name = strings.TrimSuffix(urlParts[len(urlParts)-1], ".git")
		}

		return &BackupRepo{
			Name:         name,
			SrcRepo:      srcRepo,
			PullInterval: pullInterval,
			S3url:        s3url,
			S3bucket:     s3bucket,
			LocalPath:    localPath,
		}, nil
	}

	// Repository doesn't exist, clone it
	fmt.Printf("Cloning repository from %s to %s\n", srcRepoURL, localPath)

	_, err = git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      srcRepoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %v", err)
	}

	srcRepo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %v", err)
	}

	// Extract repository name from remote URL if not provided
	if name == "" {
		urlParts := strings.Split(srcRepoURL, "/")
		name = strings.TrimSuffix(urlParts[len(urlParts)-1], ".git")
	}

	// Return a new BackupRepo instance with the cloned source repo, pull interval, S3 bucket, and local path
	return &BackupRepo{
		Name:         name,
		SrcRepo:      srcRepo,
		PullInterval: pullInterval,
		S3url:        s3url,
		S3bucket:     s3bucket,
		LocalPath:    localPath,
	}, nil
}
