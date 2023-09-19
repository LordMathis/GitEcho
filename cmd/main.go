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

	flag.Parse()

	if *generateKey {
		key, err := encryption.GenerateEncryptionKey()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Generated encryption key:", string(key))
		return
	}

	config, err := config.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	scheduler := backup.NewBackupScheduler()
	scheduler.Start()

	for _, repo := range config.Repositories {

		repo.Storages = make([]storage.Storage, len(repo.StorageNames))

		for i, storageName := range repo.StorageNames {
			stor := config.Storages[storageName]
			repo.Storages[i] = stor.Config
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
