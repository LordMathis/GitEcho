package backuprepo

import (
	"sync"

	"github.com/LordMathis/GitEcho/pkg/storage"
)

type BackupRepoManager struct {
	mutex        sync.RWMutex
	repositories map[string]*BackupRepo
}

func NewBackupRepoManager() *BackupRepoManager {
	return &BackupRepoManager{
		mutex:        sync.RWMutex{},
		repositories: make(map[string]*BackupRepo),
	}
}

func (bm *BackupRepoManager) AddBackupRepo(r *BackupRepo) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	bm.repositories[r.Name] = r
}

func (bm *BackupRepoManager) GetBackupRepo(name string) *BackupRepo {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	return bm.repositories[name]
}

func (bm *BackupRepoManager) DeleteBackupRepo(name string) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	delete(bm.repositories, name)
}

func (bm *BackupRepoManager) GetAllBackupRepos() []*BackupRepo {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	repos := make([]*BackupRepo, 0, len(bm.repositories))
	for _, repo := range bm.repositories {
		repos = append(repos, repo)
	}

	return repos
}

func (bm *BackupRepoManager) AddStorage(repo_name string, storage storage.Storage) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	bm.repositories[repo_name].AddStorage(storage)
}

func (bm *BackupRepoManager) RemoveStorage(repo_name, storage_name string) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	bm.repositories[repo_name].RemoveStorage(storage_name)
}
