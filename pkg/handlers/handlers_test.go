package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/backuprepo/testdata"
	"github.com/LordMathis/GitEcho/pkg/handlers"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/stretchr/testify/assert"
)

type MockBackupRepoProcessor struct{}

func (m *MockBackupRepoProcessor) ProcessBackupRepo(backupRepoData *backuprepo.BackupRepoData) (*backuprepo.BackupRepo, error) {
	// Return a mock BackupRepo instance
	return &backuprepo.BackupRepo{
		// Set the required fields
	}, nil
}

// MockBackupRepoInserter is a mock implementation of the BackupRepoInserter interface
type MockBackupRepoInserter struct{}

func (m *MockBackupRepoInserter) InsertBackupRepo(backupRepo *backuprepo.BackupRepo, storageID int) error {
	// Implement the mock behavior
	return nil
}

type MockRepositoryAdder struct{}

func (m *MockRepositoryAdder) AddRepository(backupRepo *backuprepo.BackupRepo) {}

type MockStorageInserter struct{}

func (m *MockStorageInserter) InsertStorage(s *storage.Storage) (int, error) {
	return 0, nil
}

func TestHandleCreateBackupRepo(t *testing.T) {
	// Create the APIHandler instance with mock dependencies
	apiHandler := &handlers.APIHandler{
		RepositoryAdder:     &MockRepositoryAdder{},
		Dispatcher:          &backup.BackupDispatcher{},
		BackupRepoInserter:  &MockBackupRepoInserter{},
		BackupRepoProcessor: &MockBackupRepoProcessor{},
		StorageInserter:     &MockStorageInserter{},
		TemplatesDir:        "",
	}

	backupRepo := testdata.GetTestBackupRepo(t, &storage.S3Storage{})

	backupRepoData := &backuprepo.BackupRepoData{
		BackupRepo:  &backupRepo,
		StorageType: "s3",
		StorageData: "data",
	}

	// Prepare the request payload
	payload, err := json.Marshal(backupRepoData)
	assert.NoError(t, err)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, "/create-backup-repo", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the handler function
	apiHandler.HandleCreateBackupRepo(recorder, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body
	expectedBody := `{"message":"Backup repository config created successfully"}`
	assert.Equal(t, expectedBody, recorder.Body.String())
}
