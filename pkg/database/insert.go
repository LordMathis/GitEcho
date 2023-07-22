package database

import (
	"fmt"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type BackupRepoInserter interface {
	InsertOrUpdateBackupRepo(backupRepo *backuprepo.BackupRepo) error
}

type StoragesInserter interface {
	InsertOrUpdateStorage(storage *storage.Storage) error
}

type StoragesInserterImpl struct{}

func (db *Database) InsertOrUpdateBackupRepo(backupRepo *backuprepo.BackupRepo) error {
	// Start a database transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Prepare the INSERT statement for backup_repo
	stmtBackupRepo, err := tx.PrepareNamed(`
		INSERT INTO backup_repo (name, pull_interval, local_path, remote_url, git_username, git_password, git_key_path)
		VALUES (:name, :pull_interval, :local_path, :remote_url, :git_username, :git_password, :git_key_path)
		ON CONFLICT (name) DO UPDATE SET
			pull_interval = EXCLUDED.pull_interval,
			local_path = EXCLUDED.local_path,
			remote_url = EXCLUDED.remote_url,
			git_username = EXCLUDED.git_username,
			git_password = EXCLUDED.git_password,
			git_key_path = EXCLUDED.git_key_path
	`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmtBackupRepo.Close()

	stmtBackupRepoStorage, err := tx.Prepare(`
		INSERT INTO backup_repo_storage (backup_repo_name, storage_name)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmtBackupRepoStorage.Close()

	// Execute the INSERT statement for backup_repo
	_, err = stmtBackupRepo.Exec(backupRepo)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, storageName := range backupRepo.StorageNames {
		_, err = stmtBackupRepoStorage.Exec(backupRepo.Name, storageName)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertOrUpdateStorage(stor storage.Storage) error {
	stmtStorage, err := db.PrepareNamed(`
		INSERT INTO storage (name, type, data)
		VALUES (:name, :type, :data)
		ON CONFLICT (name) DO UPDATE SET type = EXCLUDED.type, data = EXCLUDED.data
	`)
	if err != nil {
		return err
	}
	defer stmtStorage.Close()

	switch s := stor.(type) {
	case *storage.S3Storage:

		dataJSON, err := s.S3StorageMarshaler.MarshalS3Storage(s)
		if err != nil {
			return err
		}

		_, err = stmtStorage.Exec(&storage.BaseStorage{
			Name: s.Name,
			Type: "s3",
			Data: string(dataJSON),
		})
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported storage type: %T", stor)
	}

	return nil
}
