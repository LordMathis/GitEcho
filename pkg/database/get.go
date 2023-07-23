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
	// Prepare the SELECT statement to fetch the backup repo
	stmtBackupRepo, err := db.DB.PrepareNamed(`
		SELECT name, pull_interval, local_path, remote_url, git_username, git_password, git_key_path
		FROM backup_repo
		WHERE name = :name
	`)
	if err != nil {
		return nil, err
	}
	defer stmtBackupRepo.Close()

	// Execute the SELECT statement to fetch the backup repo
	var backup_repo backuprepo.BackupRepo
	err = stmtBackupRepo.Get(&backup_repo, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	// Process the backup repo and storages
	processedRepo, err := backuprepo.ProcessBackupRepo(&backup_repo)
	if err != nil {
		return nil, err
	}

	return processedRepo, nil
}

// GetAllBackupRepoConfigs retrieves all stored BackupRepoConfig from the database.
func (db *Database) GetAllBackupRepos() ([]*backuprepo.BackupRepo, error) {
	// Prepare the SELECT statement to fetch all backup repos
	stmtBackupRepos, err := db.DB.Preparex(`
		SELECT name, pull_interval, local_path, remote_url, git_username, git_password, git_key_path
		FROM backup_repo
	`)
	if err != nil {
		return nil, err
	}
	defer stmtBackupRepos.Close()

	// Execute the SELECT statement to fetch all backup repos
	var backupRepos []*backuprepo.BackupRepo
	err = stmtBackupRepos.Select(&backupRepos)
	if err != nil {
		return nil, err
	}

	retBackupRepos := []*backuprepo.BackupRepo{}

	for _, backupRepo := range backupRepos {
		// Process each backup repo and storages
		processedBackupRepo, err := backuprepo.ProcessBackupRepo(backupRepo)
		if err != nil {
			return nil, err
		}

		retBackupRepos = append(retBackupRepos, processedBackupRepo)
	}

	return retBackupRepos, nil
}

func (db *Database) GetStorageByName(name string) (storage.Storage, error) {
	// Prepare the SELECT statement to fetch the storage
	stmtStorage, err := db.DB.PrepareNamed(`
		SELECT name, type, data
		FROM storage
		WHERE name = :name
	`)
	if err != nil {
		return nil, err
	}
	defer stmtStorage.Close()

	// Execute the SELECT statement to fetch the storage
	var baseStorage storage.BaseStorage
	err = stmtStorage.Get(&baseStorage, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	// Create the appropriate storage instance based on the base storage type
	storageInstance, err := storage.CreateStorage(baseStorage)
	if err != nil {
		return nil, err
	}

	return storageInstance, nil
}

func (db *Database) GetAllStorages() ([]storage.Storage, error) {
	// Prepare the SELECT statement to fetch all storages
	stmtStorages, err := db.DB.Preparex(`
		SELECT name, type, data
		FROM storage
	`)
	if err != nil {
		return nil, err
	}
	defer stmtStorages.Close()

	// Execute the SELECT statement to fetch all storages
	var baseStorages []storage.BaseStorage
	err = stmtStorages.Select(&baseStorages)
	if err != nil {
		return nil, err
	}

	var storages []storage.Storage
	for _, baseStorage := range baseStorages {
		// Create the appropriate storage instance based on the base storage type
		storageInstance, err := storage.CreateStorage(baseStorage)
		if err != nil {
			return nil, err
		}
		storages = append(storages, storageInstance)
	}

	return storages, nil
}

func (db *Database) GetBackupRepoStorages(repoName string) ([]storage.Storage, error) {
	// Prepare the SELECT statement to fetch the storages associated with the backup repo
	stmtStorages, err := db.DB.Preparex(`
		SELECT s.name, s.type, s.data
		FROM backup_repo_storage b JOIN storage s ON b.storage_name = s.name
		WHERE b.backup_repo_name = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmtStorages.Close()

	// Execute the SELECT statement to fetch the storages
	var baseStorages []*storage.BaseStorage
	err = stmtStorages.Select(&baseStorages, repoName)
	if err != nil {
		return nil, err
	}

	// Convert BaseStorage to concrete storage types using CreateStorage function
	var storages []storage.Storage
	for _, baseStorage := range baseStorages {
		storageInstance, err := storage.CreateStorage(*baseStorage)
		if err != nil {
			return nil, err
		}
		storages = append(storages, storageInstance)
	}

	return storages, nil
}
