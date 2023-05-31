package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/LordMathis/GitEcho/pkg/common"
)

type backupRepoRequest struct {
	Name                  string `json:"name"`
	RemoteUrl             string `json:"remote_url"`
	PullInterval          int    `json:"pull_interval"`
	S3url                 string `json:"s3_url"`
	S3bucket              string `json:"s3_bucket"`
	AWS_ACCESS_KEY_ID     string `json:"aws_access_key_id"`
	AWS_SECRET_ACCESS_KEY string `json:"aws_secret_access_key"`
}

func handleCreateBackupRepo(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var repoRequest backupRepoRequest

	err := json.NewDecoder(r.Body).Decode(&repoRequest)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	local_path := os.Getenv("GITECHO_DATA_PATH") + "/" + repoRequest.Name

	backup_repo, err := common.NewBackupRepo(
		repoRequest.Name,
		repoRequest.RemoteUrl,
		repoRequest.PullInterval,
		repoRequest.S3url,
		repoRequest.S3bucket,
		local_path,
	)

	if err != nil {
		log.Fatalln("There was an error creating the backup repo configuration")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backup repository config created successfully"}`))
}
