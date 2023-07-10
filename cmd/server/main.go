package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/handlers"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

func main() {

	TestInsertOrUpdateStorages()

	// generateKey := flag.Bool("g", false, "Generate encryption key and exit")
	// flag.Parse()

	// if *generateKey {
	// 	key, err := encryption.GenerateEncryptionKey()
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}

	// 	fmt.Println("Generated encryption key:", string(key))
	// 	return
	// }

	// // Check if the encryption key is provided
	// key, err := encryption.ValidateEncryptionKey()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// encryption.SetEncryptionKey(key)

	// db := initializeDatabase()
	// defer db.CloseDB()

	// dispatcher := initializeBackupDispatcher(db)
	// dispatcher.Start()

	// templatesDir := getTemplatesDirectory()

	// apiHandler := handlers.NewAPIHandler(dispatcher, db, templatesDir)

	// router := setupRouter(apiHandler)

	// port := os.Getenv("GITECHO_PORT")
	// if port == "" {
	// 	// Use a default port if the environment variable is not set
	// 	port = "8080"
	// }

	// err = http.ListenAndServe(":"+port, router)
	// if err != nil {
	// 	log.Fatalln("There's an error with the server:", err)
	// }
}

func TestInsertOrUpdateStorages() {

	encryption.SetEncryptionKey([]byte("12345678901234567890123456789012"))

	db := initializeDatabase()
	defer db.CloseDB()

	stor := &storage.S3Storage{
		Endpoint:   "http://example.com",
		Region:     "us-west-1",
		AccessKey:  "access_key",
		SecretKey:  "secret_key",
		BucketName: "my-bucket",
	}

	backupRepo := &backuprepo.BackupRepo{
		Name:         "test-repo",
		RemoteURL:    "https://github.com/example/test-repo.git",
		PullInterval: 60,
		LocalPath:    "/tmp",
		Credentials: backuprepo.Credentials{
			GitUsername: "username",
			GitPassword: "password",
			GitKeyPath:  "keypath",
		},
	}

	backupRepo.Storages = make(map[string]storage.Storage)
	backupRepo.Storages["test"] = stor

	err := db.InsertOrUpdateBackupRepo(backupRepo)
	if err != nil {
		panic(err)
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
