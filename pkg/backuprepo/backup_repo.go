package backuprepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"

	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/gitutil"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type BackupRepo struct {
	Name         string                     `json:"name" db:"name"`
	SrcRepo      *git.Repository            `json:"-"`
	RemoteURL    string                     `json:"remote_url" db:"remote_url"`
	PullInterval int                        `json:"pull_interval" db:"pull_interval"`
	Storages     map[string]storage.Storage `json:"storage"`
	LocalPath    string                     `json:"-" db:"local_path"`
	Credentials  `json:"credentials"`
}

type Credentials struct {
	GitUsername string `json:"username" db:"git_username"`
	GitPassword string `json:"password" db:"git_password"`
	GitKeyPath  string `json:"key_path" db:"git_key_path"`
}

type ParsedJSONRepo struct {
	Name         string                     `json:"name" db:"name"`
	RemoteURL    string                     `json:"remote_url" db:"remote_url"`
	PullInterval int                        `json:"pull_interval" db:"pull_interval"`
	LocalPath    string                     `json:"-" db:"local_path"`
	Storages     map[string]json.RawMessage `json:"storage"`
	Credentials  `json:"credentials"`
}

type BackupRepoProcessor interface {
	ProcessBackupRepo(parsedJSONRepo *ParsedJSONRepo) (*BackupRepo, error)
}

type BackupRepoProcessorImpl struct {
	StorageCreator storage.StorageCreator
}

func (b *BackupRepo) InitializeRepo() error {
	gitclient := gitutil.NewGitClient(b.GitUsername, b.GitPassword, b.GitKeyPath)
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

func (p *BackupRepoProcessorImpl) ProcessBackupRepo(parsedJSONRepo *ParsedJSONRepo) (*BackupRepo, error) {

	backupRepo := &BackupRepo{
		Name:         parsedJSONRepo.Name,
		RemoteURL:    parsedJSONRepo.RemoteURL,
		PullInterval: parsedJSONRepo.PullInterval,
		LocalPath:    parsedJSONRepo.LocalPath,
		Storages:     make(map[string]storage.Storage),
		Credentials:  parsedJSONRepo.Credentials,
	}

	for name, storage := range parsedJSONRepo.Storages {

		storageInstance, err := p.StorageCreator.CreateStorage(storage)
		if err != nil {
			return nil, err
		}

		backupRepo.Storages[name] = storageInstance
	}

	password := backupRepo.GitPassword
	if password != "" {
		decryptedPassword, err := encryption.Decrypt([]byte(password))
		if err != nil {
			return nil, err
		}
		backupRepo.GitPassword = string(decryptedPassword)
	}

	backupRepo.InitializeRepo()

	return backupRepo, nil
}
