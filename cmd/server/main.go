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
	var err error
	database, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	err = database.MigrateDB()
	if err != nil {
		log.Fatal(err)
	}

	defer database.CloseDB()

	dispatcher := backup.NewBackupDispatcher()

	backupRepos, err := database.GetAllBackupRepos()

	if err != nil {
		log.Fatal(err)
	}

	for _, backupRepo := range backupRepos {
		dispatcher.AddRepository(*backupRepo)
	}

	dispatcher.Start()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the absolute path to the templates directory
	templatesDir := filepath.Join(currentDir, "..", "..", "templates")

	apiHandler := handlers.APIHandler{
		Dispatcher:   dispatcher,
		Db:           database,
		TemplatesDir: templatesDir,
	}

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

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("There's an error with the server,", err)
	}
}
