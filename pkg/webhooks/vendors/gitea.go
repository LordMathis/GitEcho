package vendors

import (
	"context"
	"net/http"

	"github.com/LordMathis/GitEcho/pkg/repository"
	"github.com/LordMathis/GitEcho/pkg/webhooks"
	"github.com/go-playground/webhooks/v6/gitea"
)

func NewGiteaHandler(config *webhooks.WebhookConfig, repo *repository.BackupRepo) func(http.ResponseWriter, *http.Request) {

	hook, _ := gitea.New(gitea.Options.Secret(config.Secret))
	var hookEvents []gitea.Event

	for _, event := range config.Events {
		switch event {
		case "create":
			hookEvents = append(hookEvents, gitea.CreateEvent)
		case "push":
			hookEvents = append(hookEvents, gitea.PushEvent)
		case "pull_request":
			hookEvents = append(hookEvents, gitea.PullRequestEvent)
		case "release":
			hookEvents = append(hookEvents, gitea.ReleaseEvent)
		}
	}

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, hookEvents...)
		if err != nil {
			if err == gitea.ErrEventNotFound {
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
