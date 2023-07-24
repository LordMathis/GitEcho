package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/go-chi/chi/v5"
)

// @Summary Create a new backup repository configuration
// @Description Create a new backup repository configuration with the given data
// @Tags backup-repositories
// @Accept json
// @Produce json
// @Param backupRepo body backuprepo.BackupRepo true "Backup repository data"
// @Success 200 {object} SuccessResponse "Success response"“
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/repository [post]
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

	response := SuccessResponse{
		Message: "Backup repository config created successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// @Summary Get a backup repository by name
// @Description Get the backup repository with the given name
// @Tags backup-repositories
// @Accept json
// @Produce json
// @Param repo_name path string true "Name of the backup repository to retrieve"
// @Success 200 {object} backuprepo.BackupRepo "Backup repository data"
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Backup repository not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/repository/{repo_name} [get]
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

// @Summary Get all backup repositories
// @Description Get a list of all backup repositories
// @Tags backup-repositories
// @Accept json
// @Produce json
// @Success 200 {array} backuprepo.BackupRepo "List of backup repositories"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/repository [get]
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

// @Summary Delete a backup repository
// @Description Delete a backup repository by its name
// @Tags backup-repositories
// @Param repo_name path string true "Name of the backup repository to delete"
// @Success 200 {object} SuccessResponse "Success response"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/repository/{repo_name} [delete]
func (a *APIHandler) HandleDeleteBackupRepo(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "repo_name")

	err := a.db.DeleteBackupRepo(name)
	if err != nil {
		http.Error(w, "Failed to delete backup repository", http.StatusInternalServerError)
		return
	}

	a.backupRepoManager.DeleteBackupRepo(name)

	response := SuccessResponse{
		Message: "Backup repository deleted successfully",
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

// @Summary Get backup repository storages
// @Description Get all storages associated with a backup repository by its name
// @Tags backup-repositories
// @Param repo_name path string true "Name of the backup repository"
// @Success 200 {array} storage.Storage "List of storages associated with the backup repository"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/repository/{repo_name}/storage/ [get]
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

// @Summary Add storage to backup repository
// @Description Associate a storage with a backup repository by their names
// @Tags backup-repositories
// @Param repo_name path string true "Name of the backup repository"
// @Param storage_name path string true "Name of the storage"
// @Success 200 {object} SuccessResponse "Success response"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/repository/{repo_name}/storage/{storage_name} [post]
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

	response := SuccessResponse{
		Message: "Storage added to backup repository successfully",
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

// @Summary Remove storage from backup repository
// @Description Remove the association between a storage and a backup repository by their names
// @Tags backup-repositories
// @Param repo_name path string true "Name of the backup repository"
// @Param storage_name path string true "Name of the storage"
// @Success 200 {object} SuccessResponse "Success response"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/repository/{repo_name}/storage/{storage_name} [delete]
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

	response := SuccessResponse{
		Message: "Storage removed from backup repository successfully",
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
