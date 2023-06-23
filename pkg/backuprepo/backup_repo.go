package backuprepo

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

type BackupRepo struct {
	Name         string `json:"name" db:"name"`
	SrcRepo      *git.Repository
	RemoteUrl    string `json:"remote_url" db:"remote_url"`
	PullInterval int    `json:"pull_interval" db:"pull_interval"`
	S3URL        string `json:"s3_url" db:"s3_url"`
	S3Bucket     string `json:"s3_bucket" db:"s3_bucket"`
	LocalPath    string `db:"local_path"`
}

// NewBackupRepo creates a new BackupRepo instance
// srcRepoURL: the URL of the source repo to backup
// pullInterval: the interval (in seconds) between each pull operation
// s3bucket: the name of the S3 bucket to store the backups
// localPath: the local path where the backups will be stored
func NewBackupRepo(name string, remoteURL string, pullInterval int, s3URL string, s3Bucket string, localPath string) (*BackupRepo, error) {

	// Extract repository name from remote URL if not provided
	if name == "" {
		urlParts := strings.Split(remoteURL, "/")
		name = strings.TrimSuffix(urlParts[len(urlParts)-1], ".git")
	}

	backup_repo := &BackupRepo{
		Name:         name,
		PullInterval: pullInterval,
		S3URL:        s3URL,
		S3Bucket:     s3Bucket,
		LocalPath:    localPath,
	}

	err := backup_repo.InitializeRepo()
	if err != nil {
		return nil, err
	}

	return backup_repo, nil
}

func (b *BackupRepo) InitializeRepo() error {
	_, err := git.PlainOpen(b.LocalPath)
	if err == nil {
		// If repository exists, pull latest changes
		fmt.Printf("Pulling repository at %s\n", b.LocalPath)

		srcRepo, err := git.PlainOpen(b.LocalPath)
		if err != nil {
			return fmt.Errorf("failed to open repository: %v", err)
		}

		wt, err := srcRepo.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree: %v", err)
		}

		err = wt.Pull(&git.PullOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("failed to pull repository: %v", err)
		}

		b.SrcRepo = srcRepo

		return nil

	}

	// Repository doesn't exist, clone it
	fmt.Printf("Cloning repository from %s to %s\n", b.RemoteUrl, b.LocalPath)

	_, err = git.PlainClone(b.LocalPath, false, &git.CloneOptions{
		URL:      b.RemoteUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	srcRepo, err := git.PlainOpen(b.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %v", err)
	}

	b.SrcRepo = srcRepo

	return nil
}
