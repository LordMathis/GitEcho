package handlers

import (
	"encoding/json"
	"log"
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

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var backup_repo backuprepo.BackupRepo

	err := json.NewDecoder(r.Body).Decode(&backup_repo)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	local_path := os.Getenv("GITECHO_DATA_PATH") + "/" + backup_repo.Name
	backup_repo.LocalPath = local_path

	err = backup_repo.InitializeRepo()
	if err != nil {
		log.Fatalln("There was an error creating the backup repo configuration")
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = a.Db.InsertBackupRepo(backup_repo)
	if err != nil {
		log.Fatalln("There was an error creating the backup repo configuration")
		w.WriteHeader(http.StatusInternalServerError)
	}

	a.Dispatcher.AddRepository(backup_repo)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backup repository config created successfully"}`))
}
