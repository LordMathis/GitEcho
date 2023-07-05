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

func GetTestBackupRepo(t *testing.T, s *storage.S3Storage) backuprepo.BackupRepo {
	return backuprepo.BackupRepo{
		Name:         "test-repo",
		PullInterval: 60,
		RemoteURL:    "https://example.com",
		LocalPath:    "/tmp",
		Storage:      s,
		Credentials: backuprepo.Credentials{
			GitUsername: "test-username",
			GitPassword: "test-password",
			GitKeyPath:  "test-keypath",
		},
	}
}