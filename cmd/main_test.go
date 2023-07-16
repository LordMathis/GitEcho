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
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/gitutil"
	"github.com/LordMathis/GitEcho/pkg/server"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {

	setupTestEnvVars(t)

	db, err := database.ConnectDB()
	if err != nil {
		t.Skip("error connecting to database", err)
	}

	err = db.MigrateDB()
	if err != nil {
		t.Skip("error migrating database", err)
	}

	defer db.CloseDB()

	dispatcher := backup.NewBackupDispatcher()
	assert.NoError(t, err)
	dispatcher.Start()

	templatesDir := getTemplatesDirectory()

	apiHandler := server.NewAPIHandler(dispatcher, db, templatesDir)

	go func() {
		err := http.ListenAndServe(":8080", server.SetupRouter(apiHandler))
		if err != nil {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	s3Storage := &storage.S3Storage{
		Endpoint:   "http://127.0.0.1:9000",
		Region:     "",
		AccessKey:  "gitecho",
		SecretKey:  "gitechokey",
		BucketName: "gitecho",
	}

	err = s3Storage.InitializeS3Storage()
	assert.NoError(t, err)

	data := map[string]interface{}{
		"name":          "test-repo",
		"remote_url":    "https://github.com/LordMathis/GitEcho",
		"pull_interval": 1,
		"credentials": map[string]string{
			"git_username": "",
			"git_password": "",
			"git_key_path": "",
		},
		"storage": map[string]interface{}{
			"test": map[string]string{
				"name":        "test",
				"type":        "s3",
				"endpoint":    "http://127.0.0.1:9000",
				"region":      "",
				"access_key":  "gitecho",
				"secret_key":  "gitechokey",
				"bucket_name": "gitecho",
			},
		},
	}

	// Encode the data to JSON
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = createBackupRepo(t, jsonData)
	assert.NoError(t, err)

	//TODO: Find better way to wait for backup
	time.Sleep(2 * 60 * time.Second)

	tempDir, err := os.MkdirTemp("", "test-repo-restore")
	assert.NoError(t, err)

	err = s3Storage.DownloadDirectory("test-repo", tempDir)
	assert.NoError(t, err)

	gitClient := gitutil.NewGitClient("", "", "")
	repo, err := gitClient.OpenRepository(tempDir)
	assert.NoError(t, err)

	err = gitClient.PullChanges(repo)
	assert.NoError(t, err)
}

func setupTestEnvVars(t *testing.T) {
	encryption.SetEncryptionKey([]byte("12345678901234567890123456789012"))
	os.Setenv("DB_TYPE", "postgres")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "gitecho")
	os.Setenv("DB_PASSWORD", "gitecho")
	os.Setenv("DB_NAME", "gitecho")
	os.Setenv("GITECHO_DATA_PATH", "/tmp")
}

func createBackupRepo(t *testing.T, jsonData []byte) error {
	// Perform the HTTP request to create the backup repository
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/repository", bytes.NewBuffer(jsonData))
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
