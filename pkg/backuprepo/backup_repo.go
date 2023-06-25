package backuprepo

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/LordMathis/GitEcho/pkg/storage"
)

type BackupRepo struct {
	Name         string `json:"name" db:"name"`
	SrcRepo      *git.Repository
	RemoteURL    string `json:"remote_url" db:"remote_url"`
	PullInterval int    `json:"pull_interval" db:"pull_interval"`
	Storage      storage.Storage
	StorageID    int    `db:"storage_id"`
	LocalPath    string `db:"local_path"`
}

// NewBackupRepo creates a new BackupRepo instance
// srcRepoURL: the URL of the source repo to backup
// pullInterval: the interval (in seconds) between each pull operation
// s3bucket: the name of the S3 bucket to store the backups
// localPath: the local path where the backups will be stored
func NewBackupRepo(name string, remoteURL string, pullInterval int, localPath string, storage storage.Storage) (*BackupRepo, error) {

	// Extract repository name from remote URL if not provided
	if name == "" {
		urlParts := strings.Split(remoteURL, "/")
		name = strings.TrimSuffix(urlParts[len(urlParts)-1], ".git")
	}

	backup_repo := &BackupRepo{
		Name:         name,
		RemoteURL:    remoteURL,
		PullInterval: pullInterval,
		Storage:      storage,
		LocalPath:    localPath,
	}

	err := ValidateBackupRepo(*backup_repo)
	if err != nil {
		return nil, err
	}

	err = backup_repo.InitializeRepo()
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
	fmt.Printf("Cloning repository from %s to %s\n", b.RemoteURL, b.LocalPath)

	_, err = git.PlainClone(b.LocalPath, false, &git.CloneOptions{
		URL:      b.RemoteURL,
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

func ValidateBackupRepo(backupRepo BackupRepo) error {
	// Define regular expression patterns for validation
	namePattern := `^[a-zA-Z0-9_-]+$`
	// s3URLPattern := `^https?://.+`

	// Validate the Name field
	if backupRepo.Name == "" {
		return errors.New("name field is required")
	}
	if matched, _ := regexp.MatchString(namePattern, backupRepo.Name); !matched {
		return errors.New("name field must consist of alphanumeric characters, hyphens, and underscores only")
	}

	// Validate the PullInterval field
	if backupRepo.PullInterval <= 0 {
		return errors.New("pullInterval field must be a positive integer")
	}

	// // Validate the S3URL field (example pattern)
	// if backupRepo.S3URL != "" {
	// 	if matched, _ := regexp.MatchString(s3URLPattern, backupRepo.S3URL); !matched {
	// 		return errors.New("S3URL field must be a valid HTTP or HTTPS URL")
	// 	}
	// }

	return nil
}
