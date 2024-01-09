package vendors

import (
	"fmt"
	"net/http"

	"github.com/LordMathis/GitEcho/pkg/repository"
	"github.com/LordMathis/GitEcho/pkg/webhooks"
	"github.com/go-playground/webhooks/v6/github"
)

func NewGitHubHandler(config *webhooks.WebhookConfig, repo *repository.BackupRepo) func(http.ResponseWriter, *http.Request) {

	hook, _ := github.New(github.Options.Secret(config.Secret))

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn't one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", release)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)
		}
	}

	return handleFunc
}
