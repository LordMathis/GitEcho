package database_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type MockS3StorageMarshaler struct {
}

func (m *MockS3StorageMarshaler) MarshalS3Storage(s3Storage *storage.S3Storage) ([]byte, error) {
	return []byte("testdata"), nil
}

type storagesInserterMock struct{}

func (m *storagesInserterMock) InsertOrUpdateStorages(tx *sqlx.Tx, backupRepoName string, storages map[string]storage.Storage) error {
	return nil
}

func TestInsertOrUpdateBackupRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	// Create a Database instance with the mock StoragesInserter
	test_database := &database.Database{
		DB:               sqlx.NewDb(db, "sqlmock"),
		StoragesInserter: &storagesInserterMock{},
	}

	mock.ExpectBegin()

	// Prepare the expectation for the INSERT statement for backup_repo
	mock.ExpectPrepare("INSERT INTO backup_repo (.+)").ExpectExec().WithArgs(
		"backup1", 60, "/path/to/local", "https://remote.git", "username", "password", "/path/to/key",
	).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	// Create a sample BackupRepo
	backupRepo := &backuprepo.BackupRepo{
		Name:         "backup1",
		PullInterval: 60,
		LocalPath:    "/path/to/local",
		RemoteURL:    "https://remote.git",
		Credentials: backuprepo.Credentials{
			GitUsername: "username",
			GitPassword: "password",
			GitKeyPath:  "/path/to/key",
		},
		Storages: make(map[string]storage.Storage),
	}

	// Execute the function under test
	err = test_database.InsertOrUpdateBackupRepo(backupRepo)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestInsertOrUpdateStorages(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	// Create a new instance of the Database with the mock DB
	test_database := &database.Database{
		DB:               sqlx.NewDb(db, "sqlmock"),
		StoragesInserter: &database.StoragesInserterImpl{},
	}

	mockS3StorageMarshaler := &MockS3StorageMarshaler{}

	// Define the expected SQL queries and their results
	mock.ExpectBegin()
	mockStorageInsert := mock.ExpectPrepare("INSERT INTO storage (.+)")
	mockRepoStorageInsert := mock.ExpectPrepare("INSERT INTO backup_repo_storage (.+)")

	mockStorageInsert.ExpectExec().WithArgs("storage1", "s3", "testdata").WillReturnResult(sqlmock.NewResult(0, 1))
	mockRepoStorageInsert.ExpectExec().WithArgs("backuprepo1", "storage1").WillReturnResult(sqlmock.NewResult(0, 1))
	// Create the test data
	backupRepoName := "backuprepo1"
	storages := map[string]storage.Storage{
		"storage1": &storage.S3Storage{
			S3StorageMarshaler: mockS3StorageMarshaler,
			Endpoint:           "http://example.com",
			Region:             "us-west-1",
			AccessKey:          "access_key",
			SecretKey:          "secret_key",
			BucketName:         "my-bucket",
		},
	}

	mockTx, err := test_database.Beginx()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	// Call the function under test
	err = test_database.StoragesInserter.InsertOrUpdateStorages(mockTx, backupRepoName, storages)
	assert.NoError(t, err)

	mock.ExpectCommit()

	// Commit the mock transaction
	err = mockTx.Commit()
	assert.NoError(t, err)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
