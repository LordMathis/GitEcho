package database_test

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/backuprepo/testdata"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/jmoiron/sqlx"
)

type MockBackupRepoProcessor struct {
	Result *backuprepo.BackupRepo
}

// ProcessBackupRepo is a mock implementation of the ProcessBackupRepo method
func (m *MockBackupRepoProcessor) ProcessBackupRepo(parsedJSONRepo *backuprepo.ParsedJSONRepo) (*backuprepo.BackupRepo, error) {
	return m.Result, nil
}

func TestGetBackupRepoByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	name := "backup1"

	// Prepare the expectation for the SELECT statement to fetch the backup repo
	mock.ExpectPrepare("SELECT name, pull_interval, local_path, remote_url, git_username, git_password, git_key_path FROM backup_repo WHERE name = ?").
		ExpectQuery().WithArgs(name).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "pull_interval", "local_path", "remote_url", "git_username", "git_password", "git_key_path"}).
				AddRow(name, 60, "/path/to/local", "https://remote.git", "username", "password", "/path/to/key"),
		)

	// Prepare the expectation for the SELECT statement to fetch the storages
	mock.ExpectPrepare("SELECT s.name, s.type, s.data FROM backup_repo_storage b JOIN storage s ON b.storage_name = s.name WHERE b.backup_repo_name = \\$1").
		ExpectQuery().WithArgs(name).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "type", "data"}).
				AddRow("storage1", "s3", `{"key": "value"}`).
				AddRow("storage2", "local", `{"key": "value"}`),
		)

	testBackupRepo := testdata.GetTestBackupRepo(t)
	storage1 := testdata.GetTestS3Storage(t)
	storage2 := testdata.GetTestS3Storage(t)

	testBackupRepo.Storages["storage1"] = &storage1
	testBackupRepo.Storages["storage2"] = &storage2

	// Create a Database instance with the mock BackupRepoProcessor and DB
	database := &database.Database{
		DB: sqlx.NewDb(db, "sqlmock"),
		BackupRepoProcessor: &MockBackupRepoProcessor{
			Result: &testBackupRepo,
		},
	}

	// Execute the function under test
	backupRepo, err := database.GetBackupRepoByName(name)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(backupRepo, &testBackupRepo) {
		t.Errorf("unexpected backupRepo values:\nexpected: %+v\ngot: %+v", testBackupRepo, backupRepo)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestGetAllBackupRepos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	// Prepare the expectation for the SELECT statement to fetch all backup repos
	mock.ExpectPrepare("SELECT name, pull_interval, local_path, remote_url, git_username, git_password, git_key_path FROM backup_repo").
		ExpectQuery().
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "pull_interval", "local_path", "remote_url", "git_username", "git_password", "git_key_path"}).
				AddRow("test-repo", 60, "/tmp", "https://github.com/example/test-repo.git", "username", "password", "keypath"),
		)

	// Prepare the expectation for the SELECT statement to fetch the storages for the test-repo
	mock.ExpectPrepare("SELECT s.name, s.type, s.data FROM backup_repo_storage b JOIN storage s ON b.storage_name = s.name WHERE b.backup_repo_name = \\$1").
		ExpectQuery().WithArgs("test-repo").
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "type", "data"}).
				AddRow("s3-storage", "s3", `{"endpoint": "test-endpoint", "region": "test-region", "access_key": "test-access-key", "secret_key": "test-secret-key", "bucket_name": "test-bucket"}`),
		)

	testBackupRepo := testdata.GetTestBackupRepo(t)

	// Create a Database instance with the mock BackupRepoProcessor and DB
	database := &database.Database{
		DB: sqlx.NewDb(db, "sqlmock"),
		BackupRepoProcessor: &MockBackupRepoProcessor{
			Result: &testBackupRepo,
		},
	}

	// Execute the function under test
	backupRepos, err := database.GetAllBackupRepos()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// Assert the number of returned backupRepos
	if len(backupRepos) != 1 {
		t.Errorf("unexpected number of backupRepos, expected 1, got %d", len(backupRepos))
	}

	if !reflect.DeepEqual(backupRepos[0], &testBackupRepo) {
		t.Errorf("unexpected backupRepo values:\nexpected: %+v\ngot: %+v", &testBackupRepo, backupRepos[0])
	}
}
