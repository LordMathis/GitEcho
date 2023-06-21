package app

import (
	"log"
	"sync"
	"time"

	"github.com/LordMathis/GitEcho/pkg/common"
)

// BackupDispatcher is responsible for managing the backup process for multiple repositories.
type BackupDispatcher struct {
	repositories []common.BackupRepo
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewBackupDispatcher creates a new BackupDispatcher instance.
func NewBackupDispatcher() *BackupDispatcher {
	return &BackupDispatcher{
		repositories: make([]common.BackupRepo, 0),
		stopChan:     make(chan struct{}),
	}
}

// AddRepository adds a new repository to the backup dispatcher.
func (d *BackupDispatcher) AddRepository(repo common.BackupRepo) {
	d.repositories = append(d.repositories, repo)
}

// Start starts the backup dispatcher and runs the backup process for each repository at the specified intervals.
func (d *BackupDispatcher) Start() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for _, repo := range d.repositories {
			d.scheduleBackup(&repo)
		}
		for {
			select {
			case <-time.After(time.Second): // Adjust the interval as needed
				for _, repo := range d.repositories {
					d.scheduleBackup(&repo)
				}
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
func (d *BackupDispatcher) scheduleBackup(repo *common.BackupRepo) {
	go func() {
		for {
			select {
			case <-time.After(time.Duration(repo.PullInterval) * time.Minute): // Adjust the interval as needed
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
