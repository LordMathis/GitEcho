package database

import (
	"fmt"
	"os"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type StorageInserter interface {
	InsertStorage(s *storage.Storage) (int, error)
}

type BackupRepoProcessor interface {
	ProcessBackupRepo(backupRepoData *backuprepo.BackupRepoData) (*backuprepo.BackupRepo, error)
}

type Database struct {
	*sqlx.DB
	StorageInserter
	BackupRepoProcessor
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

	return &Database{DB: db}, nil
}

func (db *Database) CloseDB() {
	db.Close()
}
