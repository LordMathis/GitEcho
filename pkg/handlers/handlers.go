package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type APIHandler struct {
	Dispatcher   *backup.BackupDispatcher
	Db           *database.Database
	TemplatesDir string
}

func NewAPIHandler(dispatcher *backup.BackupDispatcher, db *database.Database, templatesDir string) *APIHandler {
	return &APIHandler{
		Dispatcher:   dispatcher,
		Db:           db,
		TemplatesDir: templatesDir,
	}
}

func (a *APIHandler) HandleCreateBackupRepo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var backupRepoData backuprepo.BackupRepoData
	err := json.NewDecoder(r.Body).Decode(&backupRepoData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	localPath := os.Getenv("GITECHO_DATA_PATH") + "/" + backupRepoData.Name

	var storageInstance storage.Storage
	switch backupRepoData.StorageType {
	case "s3":
		storageInstance, err = storage.NewS3StorageFromJson(backupRepoData.StorageData)
		if err != nil {
			http.Error(w, "Failed to create storage instance", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unknown storage type", http.StatusBadRequest)
		return
	}

	backupRepo := &backuprepo.BackupRepo{
		Name:         backupRepoData.Name,
		PullInterval: backupRepoData.PullInterval,
		Storage:      storageInstance,
		LocalPath:    localPath,
	}

	err = backupRepo.InitializeRepo()
	if err != nil {
		http.Error(w, "Failed to create backup repository configuration", http.StatusInternalServerError)
		return
	}

	err = a.Db.InsertBackupRepo(*backupRepo)
	if err != nil {
		http.Error(w, "Failed to create backup repository configuration", http.StatusInternalServerError)
		return
	}

	a.Dispatcher.AddRepository(*backupRepo)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backup repository config created successfully"}`))
}

func (api *APIHandler) HandleGetBackupRepos(w http.ResponseWriter, r *http.Request) {
	// Retrieve all backup repos from the database
	backupRepos, err := api.Db.GetAllBackupRepos()
	if err != nil {
		// Handle the error appropriately, e.g., return an error response
		http.Error(w, "Failed to retrieve backup repositories", http.StatusInternalServerError)
		return
	}

	// Serialize the backup repos to JSON
	backupReposJSON, err := json.Marshal(backupRepos)
	if err != nil {
		// Handle the error appropriately, e.g., return an error response
		http.Error(w, "Failed to serialize backup repositories", http.StatusInternalServerError)
		return
	}

	// Set the response content type header
	w.Header().Set("Content-Type", "application/json")

	// Write the backup repos JSON as the response body
	w.Write(backupReposJSON)
}

func (a *APIHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {

	templatePah := filepath.Join(a.TemplatesDir, "index.html")
	tmpl, err := template.ParseFiles(templatePah)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
