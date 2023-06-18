package git

import (
	"os"

	git "github.com/go-git/go-git/v5"
)

func Pull(repo *git.Repository) error {
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

func Open(localPath string) (*git.Repository, error) {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
