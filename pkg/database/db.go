package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type Database struct {
	*sqlx.DB
}

func ConnectDB() (*Database, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (db *Database) CloseDB() {
	db.Close()
}

func (db *Database) MigrateDB() error {
	// Obtain *sql.DB from *sqlx.DB
	sqlDB := db.DB.DB

	// Set up the migration source
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construct the absolute path to the migrations directory
	migrationsDir := filepath.Join(currentDir, "..", "..", "pkg", "database", "migrations")

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	// Apply migrations
	_, err = migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertBackupRepo(backupRepo backuprepo.BackupRepo) error {
	// Determine the storage type
	var storageID int
	var err error

	switch backupRepo.Storage.(type) {
	case *storage.S3Storage:
		storageID, err = db.InsertS3Storage(backupRepo.Storage.(*storage.S3Storage))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported storage type")
	}

	// Prepare the INSERT statement for backup_repo
	stmtBackupRepo, err := db.DB.PrepareNamed(`
		INSERT INTO backup_repo (name, pull_interval, storage_id, local_path)
		VALUES (:name, :pull_interval, :storage_id, :local_path)
	`)
	if err != nil {
		return err
	}
	defer stmtBackupRepo.Close()

	// Set the storage ID in the backupRepo struct
	backupRepo.StorageID = storageID

	// Execute the INSERT statement for backup_repo
	_, err = stmtBackupRepo.Exec(backupRepo)
	if err != nil {
		return err
	}

	fmt.Println("Inserted BackupRepo and Storage into the database!")

	return nil
}

// InsertS3Storage inserts S3 storage into the database and returns the storage ID or an error
func (db *Database) InsertS3Storage(s3Storage *storage.S3Storage) (int, error) {
	stmt, err := db.DB.PrepareNamed(`
		INSERT INTO storage (type, data)
		VALUES ('s3', '{"endpoint": :endpoint, "region": :region, "access_key": :access_key, "secret_key": :secret_key, "bucket_name": :bucket_name}')
		RETURNING id
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var storageID int
	err = stmt.Get(&storageID, s3Storage)
	if err != nil {
		return 0, err
	}

	return storageID, nil
}

func (db *Database) GetBackupRepoByName(name string) (*backuprepo.BackupRepo, error) {
	// Prepare the SELECT statement
	stmt, err := db.DB.Preparex(`
		SELECT backup_repo.name, backup_repo.remote_url, backup_repo.pull_interval, storage.type, storage.data, backup_repo.local_path
		FROM backup_repo
		INNER JOIN storage ON backup_repo.storage_id = storage.id
		WHERE backup_repo.name = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SELECT statement
	var backupRepoData backuprepo.BackupRepoData

	err = stmt.Get(&backupRepoData.BackupRepo, &backupRepoData.StorageType, &backupRepoData.StorageData, name)
	if err != nil {
		return nil, err
	}

	var storageInstance storage.Storage

	// Based on the storage type, unmarshal into the appropriate storage struct
	switch backupRepoData.StorageType {
	case "s3":
		storageInstance, err = storage.NewS3StorageFromJson(string(backupRepoData.StorageData))
		if err != nil {
			return nil, err
		}
	}

	backupRepo := backupRepoData.BackupRepo

	backupRepo.Storage = storageInstance
	backupRepo.InitializeRepo()

	return backupRepo, nil
}

// GetAllBackupRepoConfigs retrieves all stored BackupRepoConfig from the database.
func (db *Database) GetAllBackupRepos() ([]*backuprepo.BackupRepo, error) {

	query := `
		SELECT backup_repo.*, storage.type, storage.data
		FROM backup_repo
		INNER JOIN storage ON backup_repo.storage_id = storage.id
	`
	var backupRepoData []*backuprepo.BackupRepoData
	err := db.Select(&backupRepoData, query)
	if err != nil {
		return nil, err
	}

	backupRepos := make([]*backuprepo.BackupRepo, len(backupRepoData))
	for i, data := range backupRepoData {
		var storageInstance storage.Storage
		// Based on the storage type, unmarshal the data into the appropriate storage struct
		switch data.StorageType {
		case "s3":
			storageInstance, err = storage.NewS3StorageFromJson(data.StorageData)
			if err != nil {
				return nil, err
			}
		}

		backupRepo := data.BackupRepo
		backupRepo.Storage = storageInstance

		backupRepo.InitializeRepo()
		backupRepos[i] = backupRepo
	}

	return backupRepos, nil
}
