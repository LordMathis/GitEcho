package vendors

import (
	"context"
	"net/http"

	"github.com/LordMathis/GitEcho/pkg/repository"
	"github.com/LordMathis/GitEcho/pkg/webhooks"
	"github.com/go-playground/webhooks/v6/github"
)

func NewGitHubHandler(config *webhooks.WebhookConfig, repo *repository.BackupRepo) func(http.ResponseWriter, *http.Request) {

	hook, _ := github.New(github.Options.Secret(config.Secret))
	var hookEvents []github.Event

	for _, event := range config.Events {
		switch event {
		case "create":
			hookEvents = append(hookEvents, github.CreateEvent)
		case "push":
			hookEvents = append(hookEvents, github.PushEvent)
		case "pull_request":
			hookEvents = append(hookEvents, github.PullRequestEvent)
		case "release":
			hookEvents = append(hookEvents, github.ReleaseEvent)
		}
	}

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, hookEvents...)
		if err != nil {
			if err == github.ErrEventNotFound {
				return
			}
		}

		if payload != nil {
			go func() {
				repo.BackupAndUpload(context.Background())
			}()
		}

	}

	return handleFunc
}
