package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
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

	backupRepo.InitializeStorages()

	err = a.db.InsertOrUpdateBackupRepo(backupRepo)
	if err != nil {
		http.Error(w, "Failed to store backup repository", http.StatusInternalServerError)
		return
	}

	a.backupRepoManager.AddBackupRepo(backupRepo)
	a.scheduler.RescheduleBackup(backupRepo)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backup repository config created successfully"}`))
}

func (a *APIHandler) HandleGetBackupRepoByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "repo_name")

	backupRepo := a.backupRepoManager.GetBackupRepo(name)

	response, err := json.Marshal(backupRepo)
	if err != nil {
		http.Error(w, "Failed to serialize backup repositories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (a *APIHandler) HandleGetBackupRepos(w http.ResponseWriter, r *http.Request) {

	backupRepos := a.backupRepoManager.GetAllBackupRepos()

	backupReposJSON, err := json.Marshal(backupRepos)
	if err != nil {
		http.Error(w, "Failed to serialize backup repositories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(backupReposJSON)
}

func (a *APIHandler) HandleDeleteBackupRepo(w http.ResponseWriter, r *http.Request) {
	// Get the repository name from the URL/query parameters
	name := chi.URLParam(r, "repo_name")

	// Delete the backup repository from the database
	err := a.db.DeleteBackupRepo(name)
	if err != nil {
		// Handle the error (e.g., return appropriate HTTP response)
		http.Error(w, "Failed to delete backup repository", http.StatusInternalServerError)
		return
	}

	// Delete the backup repository from the dispatcher
	a.backupRepoManager.DeleteBackupRepo(name)

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

func (a *APIHandler) HandleGetBackupRepoStorages(w http.ResponseWriter, r *http.Request) {
	// Get the repository name from the URL/query parameters
	name := chi.URLParam(r, "repo_name")

	repo := a.backupRepoManager.GetBackupRepo(name)
	var storages []storage.Storage

	for _, stor := range repo.Storages {
		storages = append(storages, stor)
	}

	// Serialize the storages to JSON
	jsonResponse, err := json.Marshal(storages)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (a *APIHandler) HandleAddBackupRepoStorage(w http.ResponseWriter, r *http.Request) {
	// Get the repository name and storage name from the URL/query parameters
	repoName := chi.URLParam(r, "repo_name")
	storageName := chi.URLParam(r, "storage_name")

	// Check if the backup repo exists in the manager
	repo := a.backupRepoManager.GetBackupRepo(repoName)
	if repo == nil {
		http.Error(w, "Backup repository not found", http.StatusNotFound)
		return
	}

	// Check if the storage exists in the manager
	storage := a.storageManager.GetStorage(storageName)
	if storage == nil {
		http.Error(w, "Storage not found", http.StatusNotFound)
		return
	}

	// Associate the backup repo with the storage in the database
	err := a.db.InsertBackupRepoStorage(repoName, storageName)
	if err != nil {
		http.Error(w, "Failed to add storage to backup repository", http.StatusInternalServerError)
		return
	}

	repo.AddStorage(storage)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Storage added to backup repository successfully"}`))
}

func (a *APIHandler) HandleRemoveBackupRepoStorage(w http.ResponseWriter, r *http.Request) {
	// Get the repository name and storage name from the URL/query parameters
	repoName := chi.URLParam(r, "repo_name")
	storageName := chi.URLParam(r, "storage_name")

	repo := a.backupRepoManager.GetBackupRepo(repoName)
	if repo == nil {
		http.Error(w, "Backup repository not found", http.StatusNotFound)
		return
	}

	err := a.db.DeleteBackupRepoStorage(repoName, storageName)
	if err != nil {
		http.Error(w, "Failed to remove storage from backup repository", http.StatusInternalServerError)
		return
	}

	repo.RemoveStorage(storageName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Storage removed from backup repository successfully"}`))
}
