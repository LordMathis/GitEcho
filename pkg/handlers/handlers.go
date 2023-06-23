package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
)

type APIHandler struct {
	Dispatcher *backup.BackupDispatcher
	Db         *database.Database
}

func (a *APIHandler) HandleCreateBackupRepo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var backupRepo backuprepo.BackupRepo
	err := json.NewDecoder(r.Body).Decode(&backupRepo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	localPath := os.Getenv("GITECHO_DATA_PATH") + "/" + backupRepo.Name
	backupRepo.LocalPath = localPath

	err = backupRepo.InitializeRepo()
	if err != nil {
		http.Error(w, "Failed to create backup repository configuration", http.StatusInternalServerError)
		return
	}

	err = a.Db.InsertBackupRepo(backupRepo)
	if err != nil {
		http.Error(w, "Failed to create backup repository configuration", http.StatusInternalServerError)
		return
	}

	a.Dispatcher.AddRepository(backupRepo)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backup repository config created successfully"}`))
}
