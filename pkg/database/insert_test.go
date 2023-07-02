package database_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// MockS3StorageInserter implements the S3StorageInserter interface
type MockStorageInserter struct {
	InsertStorageFn func(s *storage.Storage) (int, error)
}

// InsertS3Storage calls the mock function
func (m *MockStorageInserter) InsertStorage(s *storage.Storage) (int, error) {
	return m.InsertStorageFn(s)
}

func TestInsertBackupRepo(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	setEncryptionKey(t)

	// Set up the expected InsertS3Storage return value
	mockStorageID := 123
	mockS3Storage := getTestS3Storage(t)
	mockInsertStorage := func(s *storage.Storage) (int, error) {
		return mockStorageID, nil
	}
	inserter := &MockStorageInserter{
		InsertStorageFn: mockInsertStorage,
	}

	// Create the Database instance using the mock connection
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	database := &database.Database{
		DB:              sqlxDB,
		StorageInserter: inserter,
	}

	// Prepare the expected INSERT statement for backup_repo
	mock.ExpectPrepare("INSERT INTO backup_repo").
		ExpectExec().
		WithArgs(
			sqlmock.AnyArg(), // :name
			sqlmock.AnyArg(), // :pull_interval
			mockStorageID,    // :storage_id
			sqlmock.AnyArg(), // :local_path
			sqlmock.AnyArg(), // :git_username
			sqlmock.AnyArg(), // :git_password
			sqlmock.AnyArg(), // :git_key_path
			sqlmock.AnyArg(), // :remote_url
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	backupRepo := getTestBackupRepo(t, &mockS3Storage)

	err = database.InsertBackupRepo(&backupRepo)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestInsertS3Storage(t *testing.T) {
	// Create a new mock DB and expect the necessary query and arguments
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	setEncryptionKey(t)

	mock.ExpectPrepare("INSERT INTO storage").
		ExpectQuery().
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Create a new instance of the Database and call the function
	database := &database.Database{DB: sqlxDB}
	s3Storage := getTestS3Storage(t)

	storageID, err := database.InsertS3Storage(&s3Storage)
	assert.NoError(t, err)
	assert.Equal(t, 1, storageID)

	// Assert that all the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
