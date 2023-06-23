package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
)

type Database struct {
	*sqlx.DB
}

func ConnectDB() (*Database, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (db *Database) CloseDB() {
	db.Close()
}

func (db *Database) MigrateDB() error {
	// Obtain *sql.DB from *sqlx.DB
	sqlDB := db.DB.DB

	// Set up the migration source
	migrations := &migrate.FileMigrationSource{
		Dir: "path/to/migrations",
	}

	// Apply migrations
	_, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertBackupRepo(backupRepo backuprepo.BackupRepo) error {
	// Prepare the INSERT statement
	stmt, err := db.DB.PrepareNamed(`
		INSERT INTO backup_repo (name, remote_url, pull_interval, s3_url, s3_bucket, local_path)
		VALUES (:name, :remote_url, :pull_interval, :s3_url, :s3_bucket, :local_path)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the INSERT statement
	_, err = stmt.Exec(backupRepo)
	if err != nil {
		return err
	}

	fmt.Println("Inserted BackupRepoConfig into the database!")

	return nil
}

func (db *Database) GetBackupRepoByName(name string) (*backuprepo.BackupRepo, error) {
	// Prepare the SELECT statement
	stmt, err := db.DB.Preparex(`
		SELECT *
		FROM backup_repo
		WHERE name = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SELECT statement
	var backupRepo backuprepo.BackupRepo
	err = stmt.Get(&backupRepo, name)
	if err != nil {
		return nil, err
	}

	backupRepo.InitializeRepo()

	return &backupRepo, nil
}

// GetAllBackupRepoConfigs retrieves all stored BackupRepoConfig from the database.
func (db *Database) GetAllBackupRepos() ([]*backuprepo.BackupRepo, error) {
	query := "SELECT * FROM backup_repo_config"
	backup_repos := []*backuprepo.BackupRepo{}
	err := db.Select(&backup_repos, query)
	if err != nil {
		return nil, err
	}

	for _, backup_repo := range backup_repos {
		backup_repo.InitializeRepo()
	}

	return backup_repos, nil
}
