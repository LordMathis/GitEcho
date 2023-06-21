package main

import (
	"log"

	"github.com/LordMathis/GitEcho/pkg/app"
	"github.com/LordMathis/GitEcho/pkg/db"
	"github.com/LordMathis/GitEcho/pkg/http"
)

func main() {
	var err error
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	dispatcher := app.NewBackupDispatcher()
	dispatcher.Start()

	http.Start(dispatcher, database)
}
