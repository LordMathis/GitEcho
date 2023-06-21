package http

import (
	"log"
	"net/http"

	"github.com/LordMathis/GitEcho/pkg/app"
	"github.com/LordMathis/GitEcho/pkg/db"
)

func Start(dispatcher *app.BackupDispatcher, db *db.Database) {

	apiHandler := APIHandler{
		dispatcher: dispatcher,
		db:         db,
	}

	http.HandleFunc("/api/v1/createBackupRepo", apiHandler.handleCreateBackupRepo)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("There's an error with the server,", err)
	}
}
