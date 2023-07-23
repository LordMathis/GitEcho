package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/gitutil"
	"github.com/LordMathis/GitEcho/pkg/server"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/stretchr/testify/assert"
)

var testBackupRepo *backuprepo.BackupRepo = &backuprepo.BackupRepo{
	Name:         "test-repo",
	RemoteURL:    "https://github.com/LordMathis/GitEcho",
	PullInterval: 1,
	Credentials: backuprepo.Credentials{
		GitUsername: "",
		GitPassword: "",
		GitKeyPath:  "",
	},
}

var testStorage *storage.S3Storage = &storage.S3Storage{
	Name:       "test-storage",
	Endpoint:   "http://127.0.0.1:9000",
	Region:     "",
	AccessKey:  "gitecho",
	SecretKey:  "gitechokey",
	BucketName: "gitecho",
}

func TestIntegration(t *testing.T) {

	setupTestEnvVars(t)

	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.CloseDB()

	storageManager := initializeStorageManager(db)
	backupRepoManager := initializeBackupRepoManager(db, storageManager)
	scheduler := backup.NewBackupScheduler(backupRepoManager)

	templatesDir := getTemplatesDirectory()

	apiHandler := server.NewAPIHandler(db, backupRepoManager, storageManager, scheduler, templatesDir)

	scheduler.Start()

	go func() {
		err := http.ListenAndServe(":8080", server.SetupRouter(apiHandler))
		if err != nil {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	s3storageData, err := json.Marshal(testStorage)
	assert.NoError(t, err)

	s3Storage := &storage.BaseStorage{
		Name: "test-storage",
		Type: storage.S3StorageType,
		Data: string(s3storageData),
	}
	jsonStorage, err := json.Marshal(s3Storage)
	assert.NoError(t, err)

	err = sendPostRequest(t, "http://127.0.0.1:8080/api/v1/storage", jsonStorage)
	assert.NoError(t, err)

	jsonRepo, err := json.Marshal(testBackupRepo)
	assert.NoError(t, err)

	err = sendPostRequest(t, "http://127.0.0.1:8080/api/v1/repository", jsonRepo)
	assert.NoError(t, err)

	requestURL := fmt.Sprintf("http://127.0.0.1:8080/api/v1/repository/%s/storage/%s", testBackupRepo.Name, testStorage.Name)
	err = sendPostRequest(t, requestURL, nil)
	assert.NoError(t, err)

	//TODO: Find better way to wait for backup
	time.Sleep(2 * 60 * time.Second)

	tempDir, err := os.MkdirTemp("", "test-repo-restore")
	assert.NoError(t, err)

	stor := storageManager.GetStorage(testStorage.Name)

	err = stor.DownloadDirectory("test-repo", tempDir)
	assert.NoError(t, err)

	gitClient := gitutil.NewGitClient("", "", "")
	repo, err := gitClient.OpenRepository(tempDir)
	assert.NoError(t, err)

	err = gitClient.PullChanges(repo)
	assert.NoError(t, err)
}

func setupTestEnvVars(t *testing.T) {
	encryption.SetEncryptionKey([]byte("12345678901234567890123456789012"))
	os.Setenv("DB_TYPE", "sqlite3")
	os.Setenv("DB_PATH", "/test.db")
	os.Setenv("GITECHO_DATA_PATH", "/tmp")
}

func sendPostRequest(t *testing.T, url string, jsonData []byte) error {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create backup repository, status code: %d", resp.StatusCode)
	}

	return nil
}
