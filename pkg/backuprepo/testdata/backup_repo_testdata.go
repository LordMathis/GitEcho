package testdata

import (
	"testing"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
)

func GetTestS3Storage(t *testing.T) storage.S3Storage {
	return storage.S3Storage{
		Endpoint:   "test-endpoint",
		Region:     "test-region",
		AccessKey:  "test-access-key",
		SecretKey:  "test-secret-key",
		BucketName: "test-bucket",
	}
}

func GetTestBackupRepo(t *testing.T) backuprepo.BackupRepo {
	repo := backuprepo.BackupRepo{
		Name:         "test-repo",
		RemoteURL:    "https://github.com/example/test-repo.git",
		PullInterval: 60,
		LocalPath:    "/tmp",
		Storages:     make(map[string]storage.Storage),
		Credentials: backuprepo.Credentials{
			GitUsername: "username",
			GitPassword: "password",
			GitKeyPath:  "keypath",
		},
	}

	return repo
}
