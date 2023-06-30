package backuprepo

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/LordMathis/GitEcho/pkg/gitutil"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type BackupRepo struct {
	Name         string `json:"name" db:"name"`
	SrcRepo      *git.Repository
	RemoteURL    string `json:"remote_url" db:"remote_url"`
	PullInterval int    `json:"pull_interval" db:"pull_interval"`
	Storage      storage.Storage
	StorageID    int         `db:"storage_id"`
	LocalPath    string      `db:"local_path"`
	Credentials  Credentials `json:"credentials"`
}

type Credentials struct {
	Username string `json:"username" db:"git_username"`
	Password string `json:"password" db:"git_password"`
	KeyPath  string `json:"key_path" db:"git_key_path"`
}

// Utility struct BackupRepoData for db and api calls
type BackupRepoData struct {
	*BackupRepo
	StorageType string `db:"type"`
	StorageData string `db:"data"`
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
	gitclient := gitutil.NewGitClient(b.Credentials.Username, b.Credentials.Password, b.Credentials.KeyPath)
	repo, err := gitclient.OpenRepository(b.LocalPath)

	if err == nil {
		// If repository exists, pull the latest changes
		fmt.Printf("Pulling repository at %s\n", b.LocalPath)
		b.SrcRepo = repo
		return gitclient.PullChanges(repo)
	}

	// Repository doesn't exist, clone it
	fmt.Printf("Cloning repository from %s to %s\n", b.RemoteURL, b.LocalPath)
	repo, err = gitclient.CloneRepository(b.RemoteURL, b.LocalPath)

	if err != nil {
		return err
	}

	b.SrcRepo = repo
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

	return nil
}
