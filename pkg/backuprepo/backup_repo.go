package backuprepo

import (
	"fmt"

	"github.com/go-git/go-git/v5"

	"github.com/LordMathis/GitEcho/pkg/gitutil"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type BackupRepo struct {
	Name         string                 `yaml:"name"`
	SrcRepo      *git.Repository        `yaml:"-"`
	RemoteURL    string                 `yaml:"remote_url"`
	Schedule     string                 `yaml:"schedule"`
	StorageNames []string               `yaml:"sorage_names"`
	Storages     []*storage.BaseStorage `yaml:"-"`
	LocalPath    string                 `yaml:"-"`
	Credentials  `yaml:"credentials"`
}

type Credentials struct {
	GitUsername string `yaml:"username"`
	GitPassword string `yaml:"password"`
	GitKeyPath  string `yaml:"key_path"`
}

func (b *BackupRepo) BackupAndUpload() error {
	gitclient := gitutil.NewGitClient(b.Credentials.GitUsername, b.Credentials.GitPassword, b.Credentials.GitKeyPath)
	err := gitclient.PullChanges(b.SrcRepo)
	if err != nil {
		return err
	}

	// Upload the local directory to S3
	for _, stor := range b.Storages {

		err := stor.UploadDirectory(b.LocalPath)
		if err != nil {
			return err
		}
	}

	return nil
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

// func ValidateBackupRepo(backupRepo BackupRepo) error {
// 	// Define regular expression patterns for validation
// 	namePattern := `^[a-zA-Z0-9_-]+$`
// 	// s3URLPattern := `^https?://.+`

// 	// Validate the Name field
// 	if backupRepo.Name == "" {
// 		return errors.New("name field is required")
// 	}
// 	if matched, _ := regexp.MatchString(namePattern, backupRepo.Name); !matched {
// 		return errors.New("name field must consist of alphanumeric characters, hyphens, and underscores only")
// 	}

// 	// Validate the PullInterval field
// 	if backupRepo.PullInterval <= 0 {
// 		return errors.New("pullInterval field must be a positive integer")
// 	}

// 	return nil
// }
