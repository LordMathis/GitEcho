package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/handlers"
)

func main() {
	db := initializeDatabase()
	defer db.CloseDB()

	dispatcher := initializeBackupDispatcher(db)
	dispatcher.Start()

	templatesDir := getTemplatesDirectory()

	apiHandler := handlers.NewAPIHandler(dispatcher, db, templatesDir)

	router := setupRouter(apiHandler)

	err := http.ListenAndServe(":8080", router)
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
		dispatcher.AddRepository(*backupRepo)
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

	router.Post("/api/v1/createBackupRepo", apiHandler.HandleCreateBackupRepo)
	router.Get("/api/v1/getBackupRepos", apiHandler.HandleGetBackupRepos)
	router.Get("/", apiHandler.HandleIndex)

	return router
}
