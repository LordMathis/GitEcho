package backup

import (
	"log"
	"sync"
	"time"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
)

type RepositoryAdder interface {
	AddRepository(repo *backuprepo.BackupRepo)
}

// BackupDispatcher is responsible for managing the backup process for multiple repositories.
type BackupDispatcher struct {
	repositories map[string]*backuprepo.BackupRepo
	mutex        sync.RWMutex
	stopChan     chan struct{}
	addRepoChan  chan *backuprepo.BackupRepo
	wg           sync.WaitGroup
	RepositoryAdder
}

// NewBackupDispatcher creates a new BackupDispatcher instance.
func NewBackupDispatcher() *BackupDispatcher {
	return &BackupDispatcher{
		repositories: make(map[string]*backuprepo.BackupRepo),
		mutex:        sync.RWMutex{},
		stopChan:     make(chan struct{}),
		addRepoChan:  make(chan *backuprepo.BackupRepo),
	}
}

// AddRepository adds a new repository to the backup dispatcher.
func (d *BackupDispatcher) AddRepository(repo *backuprepo.BackupRepo) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.repositories[repo.Name] = repo
	d.addRepoChan <- repo
}

// Start starts the backup dispatcher and runs the backup process for each repository at the specified intervals.
func (d *BackupDispatcher) Start() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for _, repo := range d.getRepositories() {
			d.scheduleBackup(repo)
		}
		for {
			select {
			case repo := <-d.addRepoChan:
				d.scheduleBackup(repo)
			case <-d.stopChan:
				return
			}
		}
	}()
}

// Stop stops the backup dispatcher and waits for the backup process to complete.
func (d *BackupDispatcher) Stop() {
	close(d.stopChan)
	d.wg.Wait()
}

// scheduleBackup schedules the backup process for a single repository.
func (d *BackupDispatcher) scheduleBackup(repo *backuprepo.BackupRepo) {
	go func() {
		ticker := time.NewTicker(time.Duration(repo.PullInterval) * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := BackupAndUpload(*repo)
				if err != nil {
					log.Printf("Error backing up repository '%s': %v\n", repo.Name, err)
				}
			case <-d.stopChan:
				return
			}
		}
	}()
}

func (d *BackupDispatcher) getRepositories() []*backuprepo.BackupRepo {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	repos := make([]*backuprepo.BackupRepo, 0, len(d.repositories))
	for _, repo := range d.repositories {
		repos = append(repos, repo)
	}

	return repos
}
