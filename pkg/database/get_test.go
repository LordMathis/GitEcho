package database_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type MockBackupRepoProcessor struct {
	ProcessBackupRepoFn func(backupRepoData *backuprepo.BackupRepoData) (*backuprepo.BackupRepo, error)
}

func (m *MockBackupRepoProcessor) ProcessBackupRepo(backupRepoData *backuprepo.BackupRepoData) (*backuprepo.BackupRepo, error) {
	return m.ProcessBackupRepoFn(backupRepoData)
}

var mockBackupRepoProcessor = &MockBackupRepoProcessor{
	ProcessBackupRepoFn: func(backupRepoData *backuprepo.BackupRepoData) (*backuprepo.BackupRepo, error) {
		storageInstance := &storage.S3Storage{}
		backupRepoData.BackupRepo.Storage = storageInstance
		return backupRepoData.BackupRepo, nil
	},
}

func TestGetBackupRepoByName(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Create the Database instance using the mock connection
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	database := &database.Database{DB: sqlxDB}

	s3storage := &storage.S3Storage{}
	testBackupRepo := getTestBackupRepo(t, s3storage)

	mock.ExpectPrepare(regexp.QuoteMeta(`SELECT backup_repo.*, storage.type, storage.data
						FROM backup_repo
						INNER JOIN storage ON backup_repo.storage_id = storage.id
						WHERE backup_repo.name = $1`)).
		ExpectQuery().
		WithArgs("test-repo").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"name", "pull_interval", "storage_id", "local_path", "remote_url", "git_username", "git_password", "git_key_path", "storage.type", "storage.data",
			}).AddRow(
				testBackupRepo.Name,
				testBackupRepo.PullInterval,
				"0",
				testBackupRepo.LocalPath,
				testBackupRepo.RemoteURL,
				testBackupRepo.GitUsername,
				testBackupRepo.GitPassword,
				testBackupRepo.GitKeyPath,
				"s3",
				"data",
			),
		)

	// Assign the mock BackupRepoProcessor to the Database
	database.BackupRepoProcessor = mockBackupRepoProcessor

	// Call the function under test
	result, err := database.GetBackupRepoByName("test-repo")
	assert.NoError(t, err)
	assert.Equal(t, testBackupRepo, *result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAllBackupRepos(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Create the Database instance using the mock connection
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	database := &database.Database{DB: sqlxDB}

	s3storage := &storage.S3Storage{}
	testBackupRepo1 := getTestBackupRepo(t, s3storage)
	testBackupRepo2 := getTestBackupRepo(t, s3storage)
	testBackupRepo2.StorageID = 1

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT backup_repo.*, storage.type, storage.data
		FROM backup_repo
		INNER JOIN storage ON backup_repo.storage_id = storage.id
	`)).WillReturnRows(
		sqlmock.NewRows([]string{
			"name", "pull_interval", "storage_id", "local_path", "remote_url", "git_username", "git_password", "git_key_path", "storage.type", "storage.data",
		}).
			AddRow(
				testBackupRepo1.Name,
				testBackupRepo1.PullInterval,
				"0",
				testBackupRepo1.LocalPath,
				testBackupRepo1.RemoteURL,
				testBackupRepo1.GitUsername,
				testBackupRepo1.GitPassword,
				testBackupRepo1.GitKeyPath,
				"s3",
				"data",
			).
			AddRow(
				testBackupRepo2.Name,
				testBackupRepo2.PullInterval,
				"1",
				testBackupRepo2.LocalPath,
				testBackupRepo2.RemoteURL,
				testBackupRepo2.GitUsername,
				testBackupRepo2.GitPassword,
				testBackupRepo2.GitKeyPath,
				"s3",
				"data",
			),
	)

	// Assign the mock BackupRepoProcessor to the Database
	database.BackupRepoProcessor = mockBackupRepoProcessor

	// Call the function under test
	result, err := database.GetAllBackupRepos()
	assert.NoError(t, err)

	// Prepare the expected result
	expectedResult := []*backuprepo.BackupRepo{&testBackupRepo1, &testBackupRepo2}

	// Assert the expected result
	assert.Equal(t, expectedResult, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
