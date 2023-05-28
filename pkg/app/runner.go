package app

import (
	"os"

	"github.com/LordMathis/GitEcho/pkg/common"
	gitclient "github.com/LordMathis/GitEcho/pkg/git"
	s3storage "github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/go-git/go-git/v5"
)

func Run(repo *common.RepositoryBackupConfig) error {

	if _, err := os.Stat(repo.LocalPath); os.IsNotExist(err) {
		return err
	}

	err := gitclient.Pull(repo.SrcRepo)

	if err != nil && err == git.NoErrAlreadyUpToDate {
		return nil
	}

	if err != nil {
		return err
	}

	s3Client, err := s3storage.NewS3Client(repo.S3bucket)

	if err != nil {
		return err
	}

	s3Client.Push(repo)

	return nil
}
