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
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3Storage struct {
	Session    *session.Session `json:"-"`
	S3Client   s3iface.S3API    `json:"-"`
	Endpoint   string           `json:"endpoint" db:"endpoint"`
	Region     string           `json:"region" db:"region"`
	AccessKey  string           `json:"access_key" db:"access_key"`
	SecretKey  string           `json:"secret_key" db:"secret_key"`
	BucketName string           `json:"bucket_name" db:"bucket_name"`
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

func NewS3StorageFromJson(storageData string) (*S3Storage, error) {
	var s3Storage S3Storage
	err := json.Unmarshal([]byte(storageData), &s3Storage)
	if err != nil {
		return nil, err
	}

	session, err := getSession(s3Storage.Endpoint, s3Storage.Region, s3Storage.AccessKey, s3Storage.SecretKey)
	if err != nil {
		return nil, err
	}

	svc := s3.New(session)

	s3Storage.Session = session
	s3Storage.S3Client = svc

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
		_, err = s.S3Client.PutObject(input)
		if err != nil {
			return err
		}

		fmt.Printf("Uploaded file %s\n", key)
		return nil
	})

	return err
}
