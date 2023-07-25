package server

import (
	"path/filepath"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

type APIHandler struct {
	db                *database.Database
	backupRepoManager *backuprepo.BackupRepoManager
	storageManager    *storage.StorageManager
	scheduler         *backup.BackupScheduler
	templatesDir      string
	staticDir         string
}

func NewAPIHandler(db *database.Database, bm *backuprepo.BackupRepoManager, sm *storage.StorageManager, scheduler *backup.BackupScheduler, templatesDir string) *APIHandler {
	return &APIHandler{
		db:                db,
		backupRepoManager: bm,
		storageManager:    sm,
		scheduler:         scheduler,
		templatesDir:      templatesDir,
		staticDir:         filepath.Join(templatesDir, "static"),
	}
}

type SuccessResponse struct {
	Message string `json:"message"`
}
