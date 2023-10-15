package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/config"
)

func main() {

	configPath := flag.String("f", "config.yaml", "Path to the config file")
	restore := flag.Bool("r", false, "Restore from backup")
	help := flag.Bool("h", false, "Print help and exit")

	flag.Parse()

	if *help {
		flag.Usage = func() {
			w := flag.CommandLine.Output()
			fmt.Fprintf(w, "Usage:\n")

			flag.VisitAll(func(f *flag.Flag) {
				switch f.Name {
				case "f":
					fmt.Fprintf(w, "  -f <path> Path to the config file \n")
				case "r":
					fmt.Fprintf(w, "  -r <repository_name> <storage_name> <local_path> Restore repository from storage backup to local path\n")
				default:
					fmt.Fprintf(w, "  -%v %v\n", f.Name, f.Usage)
				}
			})

		}
		flag.Usage()
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

	err := stor.DownloadDirectory(context.Background(), repo.Name, localPath)
	if err != nil {
		return err
	}

	return nil
}
