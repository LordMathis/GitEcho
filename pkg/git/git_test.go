package git

import (
	"os"
	"testing"

	git "github.com/go-git/go-git/v5"
)

func TestClone(t *testing.T) {
	tempDir := t.TempDir()

	cloneOptions := &git.CloneOptions{
		URL:      "https://github.com/go-git/go-git.git",
		Progress: os.Stdout,
	}

	repo, err := Clone(cloneOptions, tempDir)

	if err != nil {
		t.Fatal(err)
	}

	if repo == nil {
		t.Fatal("repo is nil")
	}
}

func TestPull(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := Clone(&git.CloneOptions{
		URL:      "https://github.com/go-git/go-git.git",
		Progress: os.Stdout,
	}, tempDir)

	if err != nil {
		t.Fatal(err)
	}

	worktree, err := repo.Worktree()

	if err != nil {
		t.Fatal(err)
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		t.Fatal(err)
	}

}
