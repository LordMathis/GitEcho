package backup

import (
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/go-git/go-git/v5"
)

// BackupAndUpload takes a BackupRepoConfig, pulls the changes, and uploads them to S3.
func BackupAndUpload(repoConfig backuprepo.BackupRepo) error {
	// Get the repository's worktree
	worktree, err := repoConfig.SrcRepo.Worktree()
	if err != nil {
		return err
	}

	// Pull the changes from the remote repository
	err = worktree.Pull(&git.PullOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	// Upload the local directory to S3
	err = repoConfig.Storage.UploadDirectory(repoConfig.LocalPath)
	if err != nil {
		return err
	}

	return nil
}
