package backup

import (
	"log"
	"sync"
	"time"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
)

// BackupDispatcher is responsible for managing the backup process for multiple repositories.
type BackupScheduler struct {
	bm           *backuprepo.BackupRepoManager
	mutex        sync.RWMutex
	stopChan     chan struct{}
	stopChannels map[string]chan struct{}
	addRepoChan  chan *backuprepo.BackupRepo
	wg           sync.WaitGroup
}

// NewBackupDispatcher creates a new BackupDispatcher instance.
func NewBackupScheduler(bm *backuprepo.BackupRepoManager) *BackupScheduler {
	return &BackupScheduler{
		bm:           bm,
		mutex:        sync.RWMutex{},
		stopChannels: make(map[string]chan struct{}),
		addRepoChan:  make(chan *backuprepo.BackupRepo),
	}
}

// Start starts the backup dispatcher and runs the backup process for each repository at the specified intervals.
func (d *BackupScheduler) Start() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for _, repo := range d.bm.GetAllBackupRepos() {
			d.ScheduleBackup(repo)
		}

		<-d.stopChan

	}()
}

// Stop stops the backup dispatcher and waits for the backup process to complete.
func (d *BackupScheduler) Stop() {
	close(d.stopChan)
	d.wg.Wait()
}

// scheduleBackup schedules the backup process for a single repository.
func (d *BackupScheduler) ScheduleBackup(repo *backuprepo.BackupRepo) {
	d.stopChannels[repo.Name] = make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(repo.PullInterval) * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := repo.BackupAndUpload()
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

func (d *BackupScheduler) UnscheduleBackup(repo *backuprepo.BackupRepo) {
	if stopChan, ok := d.stopChannels[repo.Name]; ok {
		// Send a signal to stop the backup process
		close(stopChan)
		delete(d.stopChannels, repo.Name)
	}
}

func (d *BackupScheduler) RescheduleBackup(repo *backuprepo.BackupRepo) {
	d.UnscheduleBackup(repo)
	d.ScheduleBackup(repo)
}
