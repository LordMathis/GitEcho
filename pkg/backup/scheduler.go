package backup

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"

	"github.com/go-co-op/gocron"
)

// BackupDispatcher is responsible for managing the backup process for multiple repositories.
type BackupScheduler struct {
	bm           *backuprepo.BackupRepoManager
	mutex        sync.RWMutex
	stopChan     chan struct{}
	stopChannels map[string]chan chan<- bool
	addRepoChan  chan *backuprepo.BackupRepo
	wg           sync.WaitGroup
	cron         *gocron.Scheduler
}

// NewBackupDispatcher creates a new BackupDispatcher instance.
func NewBackupScheduler(bm *backuprepo.BackupRepoManager) *BackupScheduler {
	return &BackupScheduler{
		bm:           bm,
		mutex:        sync.RWMutex{},
		stopChannels: make(map[string]chan chan<- bool),
		addRepoChan:  make(chan *backuprepo.BackupRepo),
		cron:         gocron.NewScheduler(time.UTC),
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

		d.cron.StartAsync()
		<-d.stopChan
		d.cron.Stop()

	}()
}

// Stop stops the backup dispatcher and waits for the backup process to complete.
func (d *BackupScheduler) Stop() {
	close(d.stopChan)
	d.wg.Wait()
}

func (d *BackupScheduler) ScheduleBackup(repo *backuprepo.BackupRepo) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Clear existing schedule for this repo
	if stopChan, ok := d.stopChannels[repo.Name]; ok {
		stop := make(chan bool)
		stopChan <- stop
		close(stopChan)
		delete(d.stopChannels, repo.Name)
	}

	stopChan := make(chan chan<- bool)
	d.stopChannels[repo.Name] = stopChan

	if repo.Schedule != "" {
		if interval, err := strconv.Atoi(repo.Schedule); err == nil && interval > 0 {
			// Schedule as minutes interval
			d.cron.Every(uint64(interval)).Minutes().Do(func() {
				repo.BackupAndUpload()
			})
		} else {
			// Treat Schedule as a cron expression
			_, err := d.cron.Cron(repo.Schedule).Do(func() {
				repo.BackupAndUpload()
			})
			if err != nil {
				log.Printf("Error scheduling backup for repo '%s': %v\n", repo.Name, err)
			}
		}
	}
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
