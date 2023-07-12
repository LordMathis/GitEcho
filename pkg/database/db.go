package database

import (
	"fmt"
	"os"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sqlx.DB
	BackupRepoProcessor backuprepo.BackupRepoProcessor
	StoragesInserter    StoragesInserter
}

func ConnectDB() (*Database, error) {

	var db *sqlx.DB
	var err error

	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case "postgres":

		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err = sqlx.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}

	case "sqlite":
		// SQLite connection string
		dbPath := os.Getenv("DB_PATH")

		db, err = sqlx.Open("sqlite3", dbPath)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{
		DB: db,
		BackupRepoProcessor: &backuprepo.BackupRepoProcessorImpl{
			StorageCreator: &storage.StorageCreatorImpl{},
		},
		StoragesInserter: &StoragesInserterImpl{},
	}, nil
}

func (db *Database) CloseDB() {
	db.Close()
}
