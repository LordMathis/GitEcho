package main

import (
	"log"

	"github.com/LordMathis/GitEcho/pkg/app"
	"github.com/LordMathis/GitEcho/pkg/db"
	"github.com/LordMathis/GitEcho/pkg/http"
)

func main() {
	var err error
	db.DB, err = db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.CloseDB()

	dispatcher := app.NewBackupDispatcher()
	dispatcher.Start()

	http.Start()
}
