package main

import (
	"os"
	"testing"
	"time"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/config"
	"github.com/LordMathis/GitEcho/pkg/gitutil"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

var yamlConfig = `
data_path: /tmp
repositories:
- name: test-repo
  remote_url: "https://github.com/LordMathis/GitEcho"
  schedule:  "*/1 * * * *"
  storages:
  - test-storage
storages:
- name: test-storage
  type: s3
  config:
    endpoint:   "http://127.0.0.1:9000"
    access_key:  gitecho
    secret_key:  gitechokey
    bucket_name: gitecho
    disable_ssl: true
    force_path_style: true
`

func TestIntegration(t *testing.T) {

	config, err := config.ParseConfigFile([]byte(yamlConfig))
	assert.NoError(t, err)

	if !isS3StorageAvailable(config.Storages["test-storage"].Config.(*storage.S3StorageConfig)) {
		t.Skip("Minio is not available. Skipping integration test.")
	}

	scheduler := backup.NewBackupScheduler()
	scheduler.Start()

	for _, repo := range config.Repositories {

		repo.Storages = make([]storage.Storage, len(repo.StorageNames))

		for i, storageName := range repo.StorageNames {
			stor := config.Storages[storageName]
			repo.Storages[i] = stor.Config
		}

		scheduler.ScheduleBackup(repo)
	}

	//TODO: Find better way to wait for backup
	time.Sleep(2 * 60 * time.Second)

	tempDir, err := os.MkdirTemp("", "test-repo-restore")
	assert.NoError(t, err)

	stor := config.Storages["test-storage"].Config

	err = stor.DownloadDirectory("test-repo", tempDir)
	assert.NoError(t, err)

	gitClient := gitutil.NewGitClient("", "", "")
	repo, err := gitClient.OpenRepository(tempDir)
	assert.NoError(t, err)

	err = gitClient.PullChanges(repo)
	assert.NoError(t, err)
}

func isS3StorageAvailable(stor *storage.S3StorageConfig) bool {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Endpoint:    aws.String(stor.Endpoint),
		Credentials: credentials.NewStaticCredentials(stor.AccessKey, stor.SecretKey, ""),
	})
	if err != nil {
		return false
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Perform a simple S3 operation to check if Minio is reachable
	_, err = svc.ListBuckets(nil)
	return err == nil
}
