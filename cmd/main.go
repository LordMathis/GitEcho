package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/server"
	"github.com/LordMathis/GitEcho/pkg/storage"
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

	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.CloseDB()

	storageManager := initializeStorageManager(db)
	backupRepoManager := initializeBackupRepoManager(db, storageManager)
	scheduler := backup.NewBackupScheduler(backupRepoManager)

	templatesDir := getTemplatesDirectory()

	apiHandler := server.NewAPIHandler(db, backupRepoManager, storageManager, scheduler, templatesDir)
	router := server.SetupRouter(apiHandler)

	scheduler.Start()

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

func getTemplatesDirectory() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(currentDir, "..", "templates")
}

func initializeStorageManager(db *database.Database) *storage.StorageManager {
	storageManager := storage.NewStorageManager()

	stoarages, err := db.GetAllStorages()
	if err != nil {
		log.Fatalln(err)
	}

	for _, storage := range stoarages {
		storageManager.AddStorage(storage)
	}

	return storageManager
}

func initializeBackupRepoManager(db *database.Database, sm *storage.StorageManager) *backuprepo.BackupRepoManager {
	backupRepoManager := backuprepo.NewBackupRepoManager()

	repos, err := db.GetAllBackupRepos()
	if err != nil {
		log.Fatalln(err)
	}

	for _, repo := range repos {
		backupRepoManager.AddBackupRepo(repo)
		storageNames, err := db.GetBackupRepoStorageNames(repo.Name)
		if err != nil {
			log.Fatalln(err)
		}

		for _, storageName := range storageNames {
			stor := sm.GetStorage(storageName)
			if stor != nil {
				backupRepoManager.AddStorage(repo.Name, stor)
			}
		}
	}

	return backupRepoManager
}
