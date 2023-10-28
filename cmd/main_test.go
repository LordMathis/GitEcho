package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/config"
	"github.com/LordMathis/GitEcho/pkg/gitutil"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {

	configPath := "../testdata/config-test.yaml"

	config, err := config.ReadConfig(configPath)
	if err != nil {
		panic(err)
	}

	assert.NoError(t, err)

	scheduler := backup.NewBackupScheduler()
	scheduler.Start()

	repo := config.Repositories["test-repo"]
	scheduler.ScheduleBackup(repo)

	//TODO: Find better way to wait for backup
	time.Sleep(2 * 60 * time.Second)

	tempDir, err := os.MkdirTemp("", "test-repo-restore")
	assert.NoError(t, err)

	stor := repo.Storages["test-storage"]

	err = stor.DownloadDirectory(context.Background(), "gitecho/test-repo", tempDir)
	assert.NoError(t, err)

	gitClient := gitutil.NewGitClient("", "", "")
	restoredRepo, err := gitClient.OpenRepository(tempDir)
	assert.NoError(t, err)

	err = gitClient.PullChanges(restoredRepo)
	assert.NoError(t, err)
}
