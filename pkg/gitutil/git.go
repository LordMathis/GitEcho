package gitutil

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type GitClient struct {
	AuthMethod transport.AuthMethod
}

func NewGitClient(username, password, keyPath string) *GitClient {
	authMethod := GetAuth(username, password, keyPath)
	return &GitClient{
		AuthMethod: authMethod,
	}
}

func (g *GitClient) CloneRepository(remoteURL, localPath string) (*git.Repository, error) {
	fmt.Printf("Cloning repository from %s to %s\n", remoteURL, localPath)

	repo, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      remoteURL,
		Progress: os.Stdout,
		Auth:     g.AuthMethod,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %v", err)
	}

	return repo, nil
}

func (g *GitClient) OpenRepository(localPath string) (*git.Repository, error) {
	fmt.Printf("Opening repository at %s\n", localPath)

	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %v", err)
	}

	return repo, nil
}

func (g *GitClient) PullChanges(repo *git.Repository) error {
	fmt.Println("Pulling latest changes")

	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       g.AuthMethod,
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull repository: %v", err)
	}

	return nil
}

func GetAuth(username, password, keyPath string) transport.AuthMethod {

	if keyPath != "" {
		auth, _ := ssh.NewPublicKeysFromFile("git", keyPath, "")
		return auth
	}

	if username != "" && password != "" {
		return &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}

	return nil
}
