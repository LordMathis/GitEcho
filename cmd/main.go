package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/config"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

func main() {

	configPath := flag.String("f", "config.yaml", "Path to the config file")
	generateKey := flag.Bool("g", false, "Generate encryption key and exit")
	restore := flag.Bool("r", false, "Restore from backup")

	flag.Parse()

	if *generateKey {
		key, err := encryption.GenerateEncryptionKey()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Generated encryption key:", key)
		return
	}

	config, err := config.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	if *restore {
		tail := flag.Args()

		if len(tail) != 3 {
			log.Fatalln("Invalid arguments")
		}

		repoName := tail[0]
		storageName := tail[1]
		localPath := tail[2]

		err := restoreRepository(config, repoName, storageName, localPath)
		if err != nil {
			log.Fatalln(err)
		}

		return
	}

	scheduler := backup.NewBackupScheduler()
	scheduler.Start()

	for _, repo := range config.Repositories {

		repo.Storages = make(map[string]storage.Storage, len(repo.StorageNames))

		for _, storageName := range repo.StorageNames {
			stor := config.Storages[storageName]
			repo.Storages[storageName] = stor.Config
		}

		scheduler.ScheduleBackup(repo)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		<-sigs
		scheduler.Stop()
		done <- true
	}()

	<-done
}

func restoreRepository(config *config.Config, repoName string, storageName string, localPath string) error {
	repo := config.Repositories[repoName]
	stor := repo.Storages[storageName]

	err := stor.DownloadDirectory(repo.Name, localPath)
	if err != nil {
		return err
	}

	return nil
}
