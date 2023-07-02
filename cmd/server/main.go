package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/handlers"
)

func main() {

	generateKey := flag.Bool("generate-key", false, "Generate encryption key and exit")
	flag.Parse()

	if *generateKey {
		key, err := encryption.GenerateEncryptionKey()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Generated encryption key:", key)
		return
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	// Check if the encryption key is provided
	err := validateEncryptionKey(encryptionKey)
	if err != nil {
		log.Fatalln(err)
	}

	db := initializeDatabase()
	defer db.CloseDB()

	dispatcher := initializeBackupDispatcher(db)
	dispatcher.Start()

	templatesDir := getTemplatesDirectory()

	apiHandler := handlers.NewAPIHandler(dispatcher, db, templatesDir)

	router := setupRouter(apiHandler)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("There's an error with the server:", err)
	}
}

func validateEncryptionKey(encryptionKey string) error {
	if encryptionKey == "" {
		return fmt.Errorf("encryption key not set, please set the ENCRYPTION_KEY environment variable")
	}

	// Check if the encryption key has the correct size
	keySize := len(encryptionKey)
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return fmt.Errorf("invalid encryption key size, encryption key must be 16, 24, or 32 bytes in length")
	}

	return nil
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

	return filepath.Join(currentDir, "..", "..", "templates")
}

func setupRouter(apiHandler *handlers.APIHandler) *chi.Mux {
	router := chi.NewRouter()

	// Set up middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}, // Add your allowed origins here
	}))

	router.Post("/api/v1/backupRepos", apiHandler.HandleCreateBackupRepo)
	router.Get("/api/v1/backupRepos", apiHandler.HandleGetBackupRepos)
	router.Get("/", apiHandler.HandleIndex)

	return router
}
