package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LordMathis/GitEcho/pkg/backup"
	"github.com/LordMathis/GitEcho/pkg/config"
	"github.com/LordMathis/GitEcho/pkg/webhooks"
	"github.com/LordMathis/GitEcho/pkg/webhooks/vendors"
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

	webhookServer := webhooks.NewWebhookServer(":8080")

	for _, repo := range config.Repositories {
		if repo.Schedule != "" {
			scheduler.ScheduleBackup(repo)
			log.Printf("Scheduled backup for repo '%s' with schedule '%s'\n", repo.Name, repo.Schedule)
		}

		if repo.WebhookConfig != nil {
			log.Printf("Registering webhook handler for repo '%s'\n", repo.Name)

			switch repo.WebhookConfig.Vendor {
			case "github":
				webhookServer.RegisterWebhookHandler(repo.Name, vendors.NewGitHubHandler(repo.WebhookConfig, repo))
			case "gitea":
				webhookServer.RegisterWebhookHandler(repo.Name, vendors.NewGiteaHandler(repo.WebhookConfig, repo))
			case "gitlab":
				webhookServer.RegisterWebhookHandler(repo.Name, vendors.NewGitLabHandler(repo.WebhookConfig, repo))
			default:
				log.Printf("Unknown webhook vendor '%s'\n", repo.WebhookConfig.Vendor)
			}
		}
	}

	go func() {
		if err := webhookServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("Error during server ListenAndServe:", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		<-sigs
		scheduler.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := webhookServer.Shutdown(ctx); err != nil {
			log.Println("Error during server shutdown:", err)
		}

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
