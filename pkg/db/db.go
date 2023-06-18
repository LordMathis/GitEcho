package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/LordMathis/GitEcho/pkg/common"
)

type Database struct {
	*sql.DB
}

func ConnectDB() (*Database, error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create the database if it doesn't exist
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		return nil, err
	}

	// Create the backup_repo table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS backup_repo (
			name TEXTPRIMARY KEY,
			remote_url TEXT,
			pull_interval INT,
			s3_url TEXT,
			s3_bucket TEXT,
			local_path TEXT
		)
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (db *Database) CloseDB() {
	db.Close()
}

func (db *Database) InsertBackupRepo(backup_repo common.BackupRepo) error {
	// Prepare the INSERT statement
	stmt, err := db.DB.Prepare(`
		INSERT INTO backup_repo (name, remote_url, pull_interval, s3_url, s3_bucket, local_path)
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the INSERT statement
	_, err = stmt.Exec(backup_repo.Name, config.RemoteUrl, config.PullInterval, config.S3url, config.S3bucket, config.LocalPath)
	if err != nil {
		return err
	}

	fmt.Println("Inserted BackupRepoConfig into the database!")

	return nil
}

func (db *Database) GetBackupRepoConfigByName(name string) (common.BackupRepo, error) {
	// Prepare the SELECT statement
	stmt, err := db.DB.Prepare(`
		SELECT remote_url, pull_interval, s3_url, s3_bucket, local_path
		FROM backup_repo
		WHERE name = $1
	`)
	if err != nil {
		return BackupRepo{}, err
	}
	defer stmt.Close()

	// Execute the SELECT statement
	var config BackupRepo
	err = stmt.QueryRow(name).Scan(
		&config.RemoteUrl,
		&config.PullInterval,
		&config.S3url,
		&config.S3bucket,
		&config.LocalPath,
	)
	if err != nil {
		return BackupRepo{}, err
	}

	return config, nil
}
