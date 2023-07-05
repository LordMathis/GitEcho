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
	Dispatcher          *backup.BackupDispatcher
	Db                  *database.Database
	BackupRepoProcessor backuprepo.BackupRepoProcessor
	TemplatesDir        string
}

func NewAPIHandler(dispatcher *backup.BackupDispatcher, db *database.Database, templatesDir string) *APIHandler {
	return &APIHandler{
		Dispatcher: dispatcher,
		Db:         db,
		BackupRepoProcessor: &backuprepo.BackupRepoProcessorImpl{
			StorageCreator: &storage.StorageCreatorImpl{},
		},
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
	backupRepoData.LocalPath = localPath

	backupRepo, err := a.BackupRepoProcessor.ProcessBackupRepo(&backupRepoData)
	if err != nil {
		http.Error(w, "Failed to create backup repository", http.StatusInternalServerError)
		return
	}

	err = a.Db.BackupRepoInserter.InsertBackupRepo(backupRepo)
	if err != nil {
		http.Error(w, "Failed to create backup repository configuration", http.StatusInternalServerError)
		return
	}

	a.Dispatcher.RepositoryAdder.AddRepository(backupRepo)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backup repository config created successfully"}`))
}

func (api *APIHandler) HandleGetBackupRepos(w http.ResponseWriter, r *http.Request) {
	// Retrieve all backup repos from the database
	backupRepos, err := api.Db.BackupReposGetter.GetAllBackupRepos()
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
