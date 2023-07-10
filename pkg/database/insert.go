package database

import (
	"fmt"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/jmoiron/sqlx"
)

type BackupRepoInserter interface {
	InsertOrUpdateBackupRepo(backupRepo *backuprepo.BackupRepo) error
}

type StoragesInserter interface {
	InsertOrUpdateStorages(tx *sqlx.Tx, backupRepoName string, storages map[string]storage.Storage) error
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

	// Execute the INSERT statement for backup_repo
	_, err = stmtBackupRepo.Exec(backupRepo)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = db.StoragesInserter.InsertOrUpdateStorages(tx, backupRepo.Name, backupRepo.Storages)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *StoragesInserterImpl) InsertOrUpdateStorages(tx *sqlx.Tx, backupRepoName string, storages map[string]storage.Storage) error {
	stmtStorage, err := tx.PrepareNamed(`
		INSERT INTO storage (name, type, data)
		VALUES (:name, :type, :data)
		ON CONFLICT (name) DO UPDATE SET type = EXCLUDED.type, data = EXCLUDED.data
	`)
	if err != nil {
		return err
	}
	defer stmtStorage.Close()

	stmtBackupRepoStorage, err := tx.Prepare(`
		INSERT INTO backup_repo_storage (backup_repo_name, storage_name)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmtBackupRepoStorage.Close()

	for name, stor := range storages {

		switch s := stor.(type) {
		case *storage.S3Storage:

			dataJSON, err := s.S3StorageMarshaler.MarshalS3Storage(s)
			if err != nil {
				return err
			}

			_, err = stmtStorage.Exec(&storage.BaseStorage{
				Name: name,
				Type: "s3",
				Data: string(dataJSON),
			})
			if err != nil {
				return err
			}

			_, err = stmtBackupRepoStorage.Exec(backupRepoName, name)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported storage type: %T", stor)
		}
	}

	return nil
}
