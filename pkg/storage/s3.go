package storage

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/LordMathis/GitEcho/pkg/common"
)

type s3Client struct {
	client *s3.Client
	bucket string
}

func newS3Client(bucket string) (*s3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &s3Client{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *s3Client) Push(repo *common.RepositoryBackupConfig) error {
	uploader := manager.NewUploader(s.client)

	err := filepath.WalkDir(repo.LocalPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				return nil
			}

			uploadFile, err := os.Open(path)
			if err != nil {
				return err
			}

			result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
				Bucket: aws.String("my-bucket"),
				Key:    aws.String("my-object-key"),
				Body:   uploadFile,
			})

			fmt.Println(result)

			return nil
		})

	if err != nil {
		log.Println(err)
	}
}
