package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/encryption"
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

	// Encrypt the password
	password := backupRepo.Credentials.Password
	encryptedPassword, err := encryption.Encrypt([]byte(password))
	if err != nil {
		return err
	}

	// Prepare the INSERT statement for backup_repo
	stmtBackupRepo, err := db.DB.PrepareNamed(`
		INSERT INTO backup_repo (name, pull_interval, storage_id, local_path, git_username, git_password, git_key_path, remote_url)
		VALUES (:name, :pull_interval, :storage_id, :local_path, :git_username, :git_password, :git_key_path, :remote_url)
	`)
	if err != nil {
		return err
	}
	defer stmtBackupRepo.Close()

	backupRepo.Credentials.Password = string(encryptedPassword)
	backupRepo.StorageID = storageID

	_, err = stmtBackupRepo.Exec(backupRepo)
	if err != nil {
		return err
	}

	// Put plaintext password back to struct
	backupRepo.Credentials.Password = password

	fmt.Println("Inserted BackupRepo and Storage into the database!")

	return nil
}

// InsertS3Storage inserts S3 storage into the database and returns the storage ID or an error
func (db *Database) InsertS3Storage(s3Storage *storage.S3Storage) (int, error) {

	// Encrypt the access key and secret key
	encryptedAccessKey, err := encryption.Encrypt([]byte(s3Storage.AccessKey))
	if err != nil {
		return 0, err
	}

	encryptedSecretKey, err := encryption.Encrypt([]byte(s3Storage.SecretKey))
	if err != nil {
		return 0, err
	}

	// Create a new instance of S3Storage with encrypted keys
	encryptedS3Storage := &storage.S3Storage{
		Endpoint:   s3Storage.Endpoint,
		Region:     s3Storage.Region,
		AccessKey:  string(encryptedAccessKey),
		SecretKey:  string(encryptedSecretKey),
		BucketName: s3Storage.BucketName,
	}

	// Encode the struct fields as JSON
	dataJSON, err := json.Marshal(encryptedS3Storage)
	if err != nil {
		return 0, err
	}

	// Prepare the SQL statement with positional parameters
	stmt, err := db.DB.Prepare(`
		INSERT INTO storage (type, data)
		VALUES ('s3', $1::jsonb)
		RETURNING id
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var storageID int
	err = stmt.QueryRow(dataJSON).Scan(&storageID)
	if err != nil {
		return 0, err
	}

	return storageID, nil
}

func (db *Database) GetBackupRepoByName(name string) (*backuprepo.BackupRepo, error) {
	// Prepare the SELECT statement
	stmt, err := db.DB.Preparex(`
		SELECT backup_repo.*, storage.type, storage.data
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

	err = stmt.Get(&backupRepoData)
	if err != nil {
		return nil, err
	}

	return processBackupRepo(backupRepoData)
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
		backupRepo, err := processBackupRepo(*data)
		if err != nil {
			return nil, err
		}

		backupRepos[i] = backupRepo
	}

	return backupRepos, nil
}

func processBackupRepo(backupRepoData backuprepo.BackupRepoData) (*backuprepo.BackupRepo, error) {
	var storageInstance storage.Storage

	switch backupRepoData.StorageType {
	case "s3":
		storageInstance, err := storage.NewS3StorageFromJson(string(backupRepoData.StorageData))
		if err != nil {
			return nil, err
		}

		err = storageInstance.DecryptKeys()
		if err != nil {
			return nil, err
		}
	}

	password := backupRepoData.BackupRepo.Credentials.Password
	if password != "" {
		decryptedPassword, err := encryption.Decrypt([]byte(password))
		if err != nil {
			return nil, err
		}
		backupRepoData.BackupRepo.Credentials.Password = string(decryptedPassword)
	}

	backupRepo := backupRepoData.BackupRepo
	backupRepo.Storage = storageInstance

	backupRepo.InitializeRepo()

	return backupRepo, nil
}
