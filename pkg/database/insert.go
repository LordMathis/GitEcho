package database

import (
	"encoding/json"
	"fmt"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type StorageInserter interface {
	InsertStorage(s *storage.Storage) (int, error)
}

type BackupRepoInserter interface {
	InsertBackupRepo(backupRepo *backuprepo.BackupRepo) error
}

func (db *Database) InsertBackupRepo(backupRepo *backuprepo.BackupRepo) error {

	storageID, err := db.StorageInserter.InsertStorage(&backupRepo.Storage)
	if err != nil {
		return err
	}

	// Encrypt the password
	password := backupRepo.GitPassword
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

	backupRepo.GitPassword = string(encryptedPassword)
	backupRepo.StorageID = storageID

	_, err = stmtBackupRepo.Exec(backupRepo)
	if err != nil {
		return err
	}

	// Put plaintext password back to struct
	backupRepo.GitPassword = password

	fmt.Println("Inserted BackupRepo and Storage into the database!")

	return nil
}

func (db *Database) InsertStorage(s *storage.Storage) (int, error) {
	switch s := (*s).(type) {
	case *storage.S3Storage:
		return db.InsertS3Storage(s)
	default:
		return 0, fmt.Errorf("unsupported storage type")
	}
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
