package server

import (
	"path/filepath"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type APIHandler struct {
	Dispatcher           *backup.BackupDispatcher
	Db                   *database.Database
	StorageManager       *storage.StorageManager
	RepositoryAdder      backup.RepositoryAdder
	BackupRepoNameGetter database.BackupRepoNameGetter
	BackupReposGetter    database.BackupReposGetter
	BackupRepoInserter   database.BackupRepoInserter
	BackupRepoProcessor  backuprepo.BackupRepoProcessor
	TemplatesDir         string
	StaticDir            string
}

func NewAPIHandler(dispatcher *backup.BackupDispatcher, db *database.Database, templatesDir string) *APIHandler {
	return &APIHandler{
		Dispatcher:           dispatcher,
		Db:                   db,
		BackupRepoNameGetter: db,
		BackupReposGetter:    db,
		BackupRepoInserter:   db,
		RepositoryAdder:      dispatcher,
		BackupRepoProcessor:  &backuprepo.BackupRepoProcessorImpl{},
		TemplatesDir:         templatesDir,
		StaticDir:            filepath.Join(templatesDir, "static"),
	}
}
