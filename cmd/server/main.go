package main

import (
	"log"
	"net/http"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/handlers"
)

func main() {
	var err error
	database, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	dispatcher := backup.NewBackupDispatcher()
	dispatcher.Start()

	apiHandler := handlers.APIHandler{
		Dispatcher: dispatcher,
		Db:         database,
	}

	http.HandleFunc("/api/v1/createBackupRepo", apiHandler.HandleCreateBackupRepo)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("There's an error with the server,", err)
	}
}
