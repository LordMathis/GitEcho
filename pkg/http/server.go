package http

import (
	"log"
	"net/http"
)

func Start() {

	http.HandleFunc("/api/v1/createBackupRepo", handleCreateBackupRepo)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("There's an error with the server,", err)
	}
}
