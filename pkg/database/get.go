package database

import (
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type BackupRepoNameGetter interface {
	GetBackupRepoByName(name string) (*backuprepo.BackupRepo, error)
}

type BackupReposGetter interface {
	GetAllBackupRepos() ([]*backuprepo.BackupRepo, error)
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

	err = stmt.Get(&backupRepoData, name)
	if err != nil {
		return nil, err
	}

	backupRepo, err := db.BackupRepoProcessor.ProcessBackupRepo(&backupRepoData)
	if err != nil {
		return nil, err
	}

	if s3Storage, ok := backupRepo.Storage.(*storage.S3Storage); ok {
		s3Storage.DecryptKeys()
		backupRepo.Storage = s3Storage
	}

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
		backupRepo, err := db.BackupRepoProcessor.ProcessBackupRepo(data)
		if err != nil {
			return nil, err
		}

		if s3Storage, ok := backupRepo.Storage.(*storage.S3Storage); ok {
			s3Storage.DecryptKeys()
			backupRepo.Storage = s3Storage
		}

		backupRepos[i] = backupRepo
	}

	return backupRepos, nil
}
