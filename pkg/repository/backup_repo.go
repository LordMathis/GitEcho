package repository

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"

	"github.com/LordMathis/GitEcho/pkg/gitutil"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/LordMathis/GitEcho/pkg/webhooks"
)

type BackupRepo struct {
	Name          string                      `yaml:"name"`
	SrcRepo       *git.Repository             `yaml:"-"`
	RemoteURL     string                      `yaml:"remote_url"`
	Schedule      string                      `yaml:"schedule"`
	WebhookConfig *webhooks.WebhookConfig     `yaml:"webhook"`
	Storages      map[string]*storage.Storage `yaml:"storages"`
	LocalPath     string                      `yaml:"-"`
	Credentials   `yaml:"credentials"`
}

type Credentials struct {
	GitUsername string `yaml:"username"`
	GitPassword string `yaml:"password"`
	GitKeyPath  string `yaml:"key_path"`
}

func (b *BackupRepo) BackupAndUpload(ctx context.Context) error {
	gitclient := gitutil.NewGitClient(b.Credentials.GitUsername, b.Credentials.GitPassword, b.Credentials.GitKeyPath)
	err := gitclient.PullChanges(b.SrcRepo)
	if err != nil {
		return err
	}

	// Upload the local directory to S3
	for _, stor := range b.Storages {

		err := stor.UploadDirectory(ctx, b.Name, b.LocalPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BackupRepo) Initialize() error {
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
