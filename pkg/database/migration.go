package database

import (
	"os"
	"path/filepath"

	migrate "github.com/rubenv/sql-migrate"
)

func (db *Database) MigrateDB() error {
	// Obtain *sql.DB from *sqlx.DB
	sqlDB := db.DB.DB

	dbType := os.Getenv("DB_TYPE")

	// Set up the migration source
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construct the absolute path to the migrations directory
	migrationsDir := filepath.Join(currentDir, "..", "pkg", "database", "migrations")

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	// Apply migrations
	_, err = migrate.Exec(sqlDB, dbType, migrations, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}
