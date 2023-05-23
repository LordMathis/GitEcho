package git

import (
	git "github.com/go-git/go-git/v5"
)

func Clone(url string, cloneOptions *git.CloneOptions, dir string) (*git.Repository, error) {
	repo, err := git.PlainClone(dir, false, cloneOptions)
	return repo, err
}

func Pull(repo *git.Repository) error {
	worktree, err := repo.Worktree()

	if err != nil {
		return err
	}

	return worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
}
