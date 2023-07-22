package database

func (db *Database) DeleteBackupRepo(name string) error {
	// Begin a transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback if there's an error, otherwise commit

	// Delete the backup repo from the backup_repo table
	_, err = tx.Exec("DELETE FROM backup_repo WHERE name = $1", name)
	if err != nil {
		return err
	}

	// Delete the associated rows from backup_repo_storage table
	_, err = tx.Exec("DELETE FROM backup_repo_storage WHERE backup_repo_name = $1", name)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteStorage(name string) error {
	// Begin a transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback if there's an error, otherwise commit

	// Delete the storage from the storage table
	_, err = tx.Exec("DELETE FROM storage WHERE name = $1", name)
	if err != nil {
		return err
	}

	// Delete the associated rows from backup_repo_storage table
	_, err = tx.Exec("DELETE FROM backup_repo_storage WHERE storage_name = $1", name)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
