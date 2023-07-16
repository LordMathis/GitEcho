package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/server"
)

func main() {

	generateKey := flag.Bool("g", false, "Generate encryption key and exit")
	flag.Parse()

	if *generateKey {
		key, err := encryption.GenerateEncryptionKey()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Generated encryption key:", string(key))
		return
	}

	// Check if the encryption key is provided
	key, err := encryption.ValidateEncryptionKey()
	if err != nil {
		log.Fatalln(err)
	}
	encryption.SetEncryptionKey(key)

	db := initializeDatabase()
	defer db.CloseDB()

	dispatcher := initializeBackupDispatcher(db)
	dispatcher.Start()

	templatesDir := getTemplatesDirectory()

	apiHandler := server.NewAPIHandler(dispatcher, db, templatesDir)
	router := server.SetupRouter(apiHandler)

	port := os.Getenv("GITECHO_PORT")
	if port == "" {
		// Use a default port if the environment variable is not set
		port = "8080"
	}

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalln("There's an error with the server:", err)
	}
}

func initializeDatabase() *database.Database {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.MigrateDB()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func initializeBackupDispatcher(db *database.Database) *backup.BackupDispatcher {
	dispatcher := backup.NewBackupDispatcher()

	backupRepos, err := db.GetAllBackupRepos()
	if err != nil {
		log.Fatal(err)
	}

	for _, backupRepo := range backupRepos {
		dispatcher.AddRepository(backupRepo)
	}

	return dispatcher
}

func getTemplatesDirectory() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(currentDir, "..", "templates")
}