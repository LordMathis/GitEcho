package server

import (
	"encoding/json"
	"net/http"

	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/go-chi/chi/v5"
)

// @Summary Create storage
// @Description Create a new storage configuration
// @Tags storages
// @Accept json
// @Produce json
// @Param storage_type path string true "Storage type (s3)" Enums(s3)
// @Param storage body storage.S3Storage true "Storage configuration to create"
// @Success 200 {object} SuccessResponse "Success response"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /storage/{storage_type} [post]
func (a *APIHandler) HandleCreateStorage(w http.ResponseWriter, r *http.Request) {
	storageType := chi.URLParam(r, "storage_conf")

	var stor storage.Storage

	switch storageType {
	case "s3":
		var s3Storage storage.S3Storage
		err := json.NewDecoder(r.Body).Decode(&s3Storage)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		s3Storage.InitializeS3Storage()
		stor = &s3Storage
	}

	err := a.db.InsertOrUpdateStorage(stor)
	if err != nil {
		http.Error(w, "Failed to save storage configuration", http.StatusInternalServerError)
		return
	}

	a.storageManager.AddStorage(stor)

	response := SuccessResponse{
		Message: "Storage config created successfully",
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

// @Summary Get storage by name
// @Description Get the storage configuration by its name
// @Tags storages
// @Param storage_conf path string true "Name of the storage"
// @Produce json
// @Success 200 {object} storage.Storage "Storage configuration"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /storage/{storage_conf} [get]
func (a *APIHandler) HandleGetStorageByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "storage_conf")

	stor := a.storageManager.GetStorage(name)
	if stor == nil {
		http.Error(w, "Storage not found", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(stor)
	if err != nil {
		http.Error(w, "Failed to serialize storage config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// @Description Get all storage configurations
// @Tags storages
// @Produce json
// @Success 200 {array} storage.Storage "List of all storage configurations"
// @Failure 500 {string} string "Internal Server Error"
// @Router /storage/ [get]
func (a *APIHandler) HandleGetStorages(w http.ResponseWriter, r *http.Request) {

	stors := a.storageManager.GetAllStorages()

	response, err := json.Marshal(stors)
	if err != nil {
		http.Error(w, "Failed to serialize backup repositories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Delete storage
// @Description Delete the storage configuration by its name
// @Tags storages
// @Param storage_conf path string true "Name of the storage"
// @Produce json
// @Success 200 {object} SuccessResponse "Success response"
// @Failure 500 {string} string "Internal Server Error"
// @Router /storage/{storage_conf} [delete]
func (a *APIHandler) HandleDeleteStorage(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "storage_conf")

	err := a.db.DeleteStorage(name)
	if err != nil {
		http.Error(w, "Failed to delete backup repository", http.StatusInternalServerError)
		return
	}

	a.storageManager.DeleteStorage(name)

	for _, repo := range a.backupRepoManager.GetAllBackupRepos() {
		repo.RemoveStorage(name)
	}

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
