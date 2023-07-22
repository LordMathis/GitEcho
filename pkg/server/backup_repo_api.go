package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/go-chi/chi/v5"
)

func (a *APIHandler) HandleCreateBackupRepo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var backupRepo *backuprepo.BackupRepo
	err := json.NewDecoder(r.Body).Decode(&backupRepo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	localPath := os.Getenv("GITECHO_DATA_PATH") + "/" + backupRepo.Name
	backupRepo.LocalPath = localPath

	err = backupRepo.InitializeRepo()
	if err != nil {
		http.Error(w, "Failed to create backup repository", http.StatusInternalServerError)
		return
	}

	err = a.BackupRepoInserter.InsertOrUpdateBackupRepo(backupRepo)
	if err != nil {
		http.Error(w, "Failed to store backup repository", http.StatusInternalServerError)
		return
	}

	a.RepositoryAdder.AddRepository(backupRepo)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backup repository config created successfully"}`))
}

func (a *APIHandler) HandleGetBackupRepoByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	// Handle request for a specific backup repo
	backupRepo, err := a.BackupRepoNameGetter.GetBackupRepoByName(name)
	if err != nil {
		// Handle error
		http.Error(w, "Failed to retrieve backup repository", http.StatusInternalServerError)
		return
	}

	// Convert backup repo to JSON and send response
	response, err := json.Marshal(backupRepo)
	if err != nil {
		http.Error(w, "Failed to serialize backup repositories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (a *APIHandler) HandleGetBackupRepos(w http.ResponseWriter, r *http.Request) {

	// Retrieve all backup repos from the database
	backupRepos, err := a.BackupReposGetter.GetAllBackupRepos()
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

func (a *APIHandler) HandleDeleteBackupRepo(w http.ResponseWriter, r *http.Request) {
	// Get the repository name from the URL/query parameters
	name := chi.URLParam(r, "name")
	// Alternatively, if using query parameters: name := r.URL.Query().Get("name")

	// Delete the backup repository from the database
	err := a.Db.DeleteBackupRepo(name)
	if err != nil {
		// Handle the error (e.g., return appropriate HTTP response)
		http.Error(w, "Failed to delete backup repository", http.StatusInternalServerError)
		return
	}

	// Delete the backup repository from the dispatcher
	a.Dispatcher.DeleteRepository(name)

	response := map[string]string{
		"message": "Backup repository deleted successfully",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
