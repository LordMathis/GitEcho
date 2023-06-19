package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/LordMathis/GitEcho/pkg/common"
	"github.com/LordMathis/GitEcho/pkg/db"
)

func handleCreateBackupRepo(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var backup_repo common.BackupRepo

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

	err = db.DB.InsertBackupRepo(backup_repo)
	if err != nil {
		log.Fatalln("There was an error creating the backup repo configuration")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backup repository config created successfully"}`))
}
