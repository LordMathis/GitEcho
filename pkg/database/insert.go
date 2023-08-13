package database

import (
	"fmt"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

func (db *Database) InsertOrUpdateBackupRepo(backupRepo *backuprepo.BackupRepo) error {

	// Prepare the INSERT statement for backup_repo
	stmtBackupRepo, err := db.PrepareNamed(`
		INSERT INTO backup_repo (name, schedule, local_path, remote_url, git_username, git_password, git_key_path)
		VALUES (:name, :schedule, :local_path, :remote_url, :git_username, :git_password, :git_key_path)
		ON CONFLICT (name) DO UPDATE SET
			schedule = EXCLUDED.schedule,
			local_path = EXCLUDED.local_path,
			remote_url = EXCLUDED.remote_url,
			git_username = EXCLUDED.git_username,
			git_password = EXCLUDED.git_password,
			git_key_path = EXCLUDED.git_key_path
	`)
	if err != nil {
		return err
	}
	defer stmtBackupRepo.Close()

	// Execute the INSERT statement for backup_repo
	_, err = stmtBackupRepo.Exec(backupRepo)
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

		dataJSON, err := storage.MarshalS3Storage(s)
		if err != nil {
			return err
		}

		_, err = stmtStorage.Exec(&storage.BaseStorage{
			Name: s.Name,
			Type: "s3",
			Data: dataJSON,
		})
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported storage type: %T", stor)
	}

	return nil
}

func (db *Database) InsertBackupRepoStorage(repoName, storageName string) error {
	// Prepare the INSERT statement to associate the backup repo with the storage
	stmtInsert, err := db.DB.Prepare(`
		INSERT INTO backup_repo_storage (backup_repo_name, storage_name)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmtInsert.Close()

	// Execute the INSERT statement
	_, err = stmtInsert.Exec(repoName, storageName)
	if err != nil {
		return err
	}

	return nil
}
