package database

func (db *Database) DeleteBackupRepo(name string) error {
	// Start a database transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Delete the backup repository from the backup_repo table
	_, err = tx.Exec("DELETE FROM backup_repo WHERE name = $1", name)
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
