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

	apiHandler := handlers.NewAPIHandler(dispatcher, db, templatesDir)

	router := setupRouter(apiHandler)

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

	return filepath.Join(currentDir, "..", "..", "templates")
}

func setupRouter(apiHandler *handlers.APIHandler) *chi.Mux {
	router := chi.NewRouter()

	staticPath := filepath.Join(getTemplatesDirectory(), "static")
	staticURLPattern := "/static/*"

	// Set up middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}, // Add your allowed origins here
	}))

	router.Get(staticURLPattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))).ServeHTTP(w, r)
	}))

	router.Post("/api/v1/backupRepos", apiHandler.HandleCreateBackupRepo)
	router.Get("/api/v1/backupRepos", apiHandler.HandleGetBackupRepos)
	router.Get("/", apiHandler.HandleIndex)

	return router
}
