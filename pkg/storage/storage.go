package storage

import (
	"context"
	"log"

	"github.com/rclone/rclone/backend/local"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config/configmap"
	"github.com/rclone/rclone/fs/sync"
)

type Storage struct {
	fremote    fs.Fs  `yaml:"-"`
	RemoteName string `yaml:"remote_name"`
	configPath string `yaml:"config_path"`
}

func (s *Storage) InitializeStorage() error {
	fremote, err := fs.NewFs(context.Background(), s.RemoteName+":"+s.configPath)
	if err != nil {
		return err
	}

	s.fremote = fremote
	return nil
}

func (s *Storage) UploadDirectory(ctx context.Context, repoName, localPath string) error {

	flocal, err := local.NewFs(ctx, localPath, repoName, configmap.New())
	if err != nil {
		return err
	}

	err = sync.Sync(ctx, s.fremote, flocal, true)
	if err != nil {
		return err
	}

	log.Printf("Directory '%s' uploaded to '%s' successfully.", localPath, s.fremote.String())
	return nil
}

func (s *Storage) DownloadDirectory(ctx context.Context, repoName, localPath string) error {

	flocal, err := local.NewFs(ctx, localPath, repoName, configmap.New())
	if err != nil {
		return err
	}

	err = sync.Sync(ctx, flocal, s.fremote, true)
	if err != nil {
		return err
	}

	log.Printf("Directory '%s' uploaded to '%s' successfully.", localPath, s.fremote.String())
	return nil

}
