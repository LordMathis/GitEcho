package server

import (
	"encoding/json"
	"net/http"

	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/go-chi/chi/v5"
)

func (a *APIHandler) HandleCreateStorage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var baseStorage storage.BaseStorage
	err := json.NewDecoder(r.Body).Decode(&baseStorage)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stor, err := storage.CreateStorage(baseStorage)
	if err != nil {
		http.Error(w, "Failed to create storage", http.StatusInternalServerError)
		return
	}

	err = a.db.InsertOrUpdateStorage(stor)
	if err != nil {
		http.Error(w, "Failed to save storage configuration", http.StatusInternalServerError)
		return
	}

	a.storageManager.AddStorage(stor)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Storage config created successfully"}`))
}

func (a *APIHandler) HandleGetStorageByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "storage_name")

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

func (a *APIHandler) HandleDeleteStorage(w http.ResponseWriter, r *http.Request) {
	// Get the repository name from the URL/query parameters
	name := chi.URLParam(r, "storage_name")

	// Delete the backup repository from the database
	err := a.db.DeleteStorage(name)
	if err != nil {
		// Handle the error (e.g., return appropriate HTTP response)
		http.Error(w, "Failed to delete backup repository", http.StatusInternalServerError)
		return
	}

	a.storageManager.DeleteStorage(name)

	for _, repo := range a.backupRepoManager.GetAllBackupRepos() {
		repo.RemoveStorage(name)
	}

	// Delete the backup repository from the dispatcher
	response := map[string]string{
		"message": "Remote storage deleted successfully",
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
