package backup

import (
	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/gitutil"
)

// BackupAndUpload takes a BackupRepoConfig, pulls the changes, and uploads them to S3.
func BackupAndUpload(repo *backuprepo.BackupRepo) error {
	gitclient := gitutil.NewGitClient(repo.Credentials.GitUsername, repo.Credentials.GitPassword, repo.Credentials.GitKeyPath)
	err := gitclient.PullChanges(repo.SrcRepo)
	if err != nil {
		return err
	}

	// Upload the local directory to S3
	for _, storage := range repo.Storages {
		err := storage.UploadDirectory(repo.LocalPath)
		if err != nil {
			return err
		}
	}

	return nil
}
