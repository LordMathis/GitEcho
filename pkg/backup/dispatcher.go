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
	stopChannels map[string]chan struct{}
	addRepoChan  chan *backuprepo.BackupRepo
	wg           sync.WaitGroup
	RepositoryAdder
}

// NewBackupDispatcher creates a new BackupDispatcher instance.
func NewBackupDispatcher() *BackupDispatcher {
	return &BackupDispatcher{
		repositories: make(map[string]*backuprepo.BackupRepo),
		mutex:        sync.RWMutex{},
		stopChannels: make(map[string]chan struct{}),
		addRepoChan:  make(chan *backuprepo.BackupRepo),
	}
}

// AddRepository adds a new repository to the backup dispatcher.
func (d *BackupDispatcher) AddRepository(repo *backuprepo.BackupRepo) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if existingRepo, ok := d.repositories[repo.Name]; ok {
		// Repository with the same name already exists, unschedule it
		d.unscheduleBackup(existingRepo)
	}

	d.repositories[repo.Name] = repo
	d.scheduleBackup(repo)
}

func (d *BackupDispatcher) DeleteRepository(name string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if repo, ok := d.repositories[name]; ok {
		d.unscheduleBackup(repo)
		delete(d.repositories, name)
	}
}

// Start starts the backup dispatcher and runs the backup process for each repository at the specified intervals.
func (d *BackupDispatcher) Start() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for _, repo := range d.getRepositories() {
			d.scheduleBackup(repo)
		}

		<-d.stopChan

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
				err := BackupAndUpload(repo)
				if err != nil {
					log.Printf("Error backing up repository '%s': %v\n", repo.Name, err)
				}
			case <-d.stopChan:
				return
			case <-d.stopChannels[repo.Name]:
				return
			}
		}
	}()
}

func (d *BackupDispatcher) unscheduleBackup(repo *backuprepo.BackupRepo) {
	if stopChan, ok := d.stopChannels[repo.Name]; ok {
		// Send a signal to stop the backup process
		close(stopChan)
		delete(d.stopChannels, repo.Name)
	}
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
