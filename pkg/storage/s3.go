package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	once    sync.Once
	sess    *session.Session
	sessErr error
)

func getSession() (*session.Session, error) {
	once.Do(func() {
		var sessConfig *aws.Config

		// Check if a custom region is specified
		region := os.Getenv("AWS_REGION")
		if region != "" {
			sessConfig = &aws.Config{
				Region: aws.String(region),
			}
		}

		sess, sessErr = session.NewSession(sessConfig)
	})
	return sess, sessErr
}

// UploadDirectory uploads the files in the specified directory (including subdirectories) to an S3 storage bucket.
// If a file already exists in the remote storage, it will be overwritten.
func UploadDirectory(bucketName, directoryPath string) error {
	// Get the session
	sess, err := getSession()
	if err != nil {
		return err
	}

	// Create a new S3 service client
	svc := s3.New(sess)

	// WalkDir through the directory recursively
	err = filepath.WalkDir(directoryPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Open the file
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// Prepare the S3 object key by preserving the directory structure
		relPath, err := filepath.Rel(directoryPath, path)
		if err != nil {
			return err
		}
		key := filepath.Join(relPath)

		// Create the input parameters for the S3 PutObject operation
		input := &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   f,
		}

		// Perform the S3 PutObject operation
		_, err = svc.PutObject(input)
		if err != nil {
			return err
		}

		fmt.Printf("Uploaded file %s\n", key)
		return nil
	})

	return err
}
