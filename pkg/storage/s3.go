package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	Session    *session.Session
	Endpoint   string `db:"endpoint"`
	Region     string `db:"region"`
	AccessKey  string `db:"access_key"`
	SecretKey  string `db:"secret_key"`
	BucketName string `db:"bucket_name"`
}

func getSession(endpoint, region, accessKey, secretKey string) (*session.Session, error) {
	config := &aws.Config{}

	if endpoint != "" {
		config.Endpoint = aws.String(endpoint)
	}

	if region != "" {
		config.Region = aws.String(region)
	}

	if accessKey != "" && secretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(accessKey, secretKey, "")
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func NewS3Storage(endpoint string, region string, accessKey string, secretKey string, bucketName string) (*S3Storage, error) {
	session, err := getSession(endpoint, region, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		Session:    session,
		Endpoint:   endpoint,
		Region:     region,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		BucketName: bucketName,
	}, nil
}

func NewS3StorageFromJson(storageData string) (*S3Storage, error) {
	var s3Storage S3Storage
	err := json.Unmarshal([]byte(storageData), &s3Storage)
	if err != nil {
		return nil, err
	}
	return &s3Storage, nil
}

func (s *S3Storage) DecryptKeys() error {
	// Decrypt the access key and secret key
	decryptedAccessKey, err := encryption.Decrypt([]byte(s.AccessKey))
	if err != nil {
		return err
	}

	decryptedSecretKey, err := encryption.Decrypt([]byte(s.SecretKey))
	if err != nil {
		return err
	}

	// Update the access key and secret key with the decrypted values
	s.AccessKey = string(decryptedAccessKey)
	s.SecretKey = string(decryptedSecretKey)

	return nil
}

// UploadDirectory uploads the files in the specified directory (including subdirectories) to an S3 storage bucket.
// If a file already exists in the remote storage, it will be overwritten.
func (s *S3Storage) UploadDirectory(directoryPath string) error {

	// Create a new S3 service client
	svc := s3.New(s.Session)

	// WalkDir through the directory recursively
	err := filepath.WalkDir(directoryPath, func(path string, d os.DirEntry, err error) error {
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
			Bucket: aws.String(s.BucketName),
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
