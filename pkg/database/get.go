package database

import (
	"encoding/json"

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
	var parsedRepo backuprepo.ParsedJSONRepo
	err = stmtBackupRepo.Get(&parsedRepo, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	// Prepare the SELECT statement to fetch the storages
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
	var storageData []storage.BaseStorage
	err = stmtStorages.Select(&storageData, parsedRepo.Name)
	if err != nil {
		return nil, err
	}

	for _, stor := range storageData {

		var dataMap map[string]interface{}
		err = json.Unmarshal([]byte(stor.Data), &dataMap)
		if err != nil {
			return nil, err
		}

		// Add the additional attribute to the data map
		dataMap["type"] = stor.Type // Add your desired value here

		// Encode the modified data map back to JSON
		updatedDataJSON, err := json.Marshal(dataMap)
		if err != nil {
			return nil, err
		}

		if parsedRepo.Storages == nil {
			parsedRepo.Storages = make(map[string]json.RawMessage)
		}
		parsedRepo.Storages[stor.Name] = updatedDataJSON
	}

	// Process the backup repo and storages
	backupRepo, err := db.BackupRepoProcessor.ProcessBackupRepo(&parsedRepo)
	if err != nil {
		return nil, err
	}

	return backupRepo, nil
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
	var parsedRepos []*backuprepo.ParsedJSONRepo
	err = stmtBackupRepos.Select(&parsedRepos)
	if err != nil {
		return nil, err
	}

	// If there are no backup repos, return an empty slice
	if len(parsedRepos) == 0 {
		return []*backuprepo.BackupRepo{}, nil
	}

	var backupRepos []*backuprepo.BackupRepo
	for _, parsedRepo := range parsedRepos {
		// Prepare the SELECT statement to fetch the storages for each backup repo
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
		var storageData []storage.BaseStorage
		err = stmtStorages.Select(&storageData, parsedRepo.Name)
		if err != nil {
			return nil, err
		}

		for _, storage := range storageData {
			var dataMap map[string]interface{}
			err = json.Unmarshal([]byte(storage.Data), &dataMap)
			if err != nil {
				return nil, err
			}

			// Add the additional attribute to the data map
			dataMap["type"] = storage.Type // Add your desired value here

			// Encode the modified data map back to JSON
			updatedDataJSON, err := json.Marshal(dataMap)
			if err != nil {
				return nil, err
			}

			if parsedRepo.Storages == nil {
				parsedRepo.Storages = make(map[string]json.RawMessage)
			}
			parsedRepo.Storages[storage.Name] = updatedDataJSON
		}

		// Process each backup repo and storages
		backupRepo, err := db.BackupRepoProcessor.ProcessBackupRepo(parsedRepo)
		if err != nil {
			return nil, err
		}

		backupRepos = append(backupRepos, backupRepo)
	}

	return backupRepos, nil
}
