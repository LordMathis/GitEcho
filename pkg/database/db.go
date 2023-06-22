package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
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

func (db *Database) InsertBackupRepo(backup_repo backuprepo.BackupRepo) error {
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
	_, err = stmt.Exec(backup_repo.Name, backup_repo.RemoteUrl, backup_repo.PullInterval, backup_repo.S3URL, backup_repo.S3Bucket, backup_repo.LocalPath)
	if err != nil {
		return err
	}

	fmt.Println("Inserted BackupRepoConfig into the database!")

	return nil
}

func (db *Database) GetBackupRepoByName(name string) (*backuprepo.BackupRepo, error) {
	// Prepare the SELECT statement
	stmt, err := db.DB.Prepare(`
		SELECT name, remote_url, pull_interval, s3_url, s3_bucket, local_path
		FROM backup_repo
		WHERE name = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SELECT statement
	var backup_repo *backuprepo.BackupRepo
	err = stmt.QueryRow(name).Scan(
		&backup_repo.Name,
		&backup_repo.RemoteUrl,
		&backup_repo.PullInterval,
		&backup_repo.S3URL,
		&backup_repo.S3Bucket,
		&backup_repo.LocalPath,
	)
	if err != nil {
		return nil, err
	}

	backup_repo.InitializeRepo()

	return backup_repo, nil
}

// GetAllBackupRepoConfigs retrieves all stored BackupRepoConfig from the database.
func (db *Database) GetAllBackupRepos() ([]*backuprepo.BackupRepo, error) {
	query := "SELECT name, remote_url, pull_interval, s3_url, s3_bucket, local_path FROM backup_repo"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var backupRepos []*backuprepo.BackupRepo

	for rows.Next() {
		var backupRepo *backuprepo.BackupRepo
		err := rows.Scan(
			&backupRepo.Name,
			&backupRepo.RemoteUrl,
			&backupRepo.PullInterval,
			&backupRepo.S3URL,
			&backupRepo.S3Bucket,
			&backupRepo.LocalPath,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}

		backupRepo.InitializeRepo()

		backupRepos = append(backupRepos, backupRepo)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error occurred while iterating over rows: %v", err)
		return nil, err
	}

	return backupRepos, nil
}
