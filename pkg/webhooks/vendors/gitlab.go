package vendors

import (
	"context"
	"net/http"

	"github.com/LordMathis/GitEcho/pkg/repository"
	"github.com/LordMathis/GitEcho/pkg/webhooks"
	"github.com/go-playground/webhooks/v6/gitlab"
)

func NewGitLabHandler(config *webhooks.WebhookConfig, repo *repository.BackupRepo) func(http.ResponseWriter, *http.Request) {

	hook, _ := gitlab.New(gitlab.Options.Secret(config.Secret))
	var hookEvents []gitlab.Event

	for _, event := range config.Events {
		switch event {
		case "Push Hook":
			hookEvents = append(hookEvents, gitlab.PushEvents)
		case "Tag Push Hook":
			hookEvents = append(hookEvents, gitlab.TagEvents)
		case "Merge Request Hook":
			hookEvents = append(hookEvents, gitlab.MergeRequestEvents)
		}
	}

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, hookEvents...)
		if err != nil {
			if err == gitlab.ErrEventNotFound {
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
