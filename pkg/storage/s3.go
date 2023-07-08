package storage

import (
	"encoding/json"
	"fmt"
	"io"
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
	} else {
		config.Region = aws.String("us-east-1")
	}

	if accessKey != "" && secretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(accessKey, secretKey, "")
	}

	config.DisableSSL = aws.Bool(true)
	config.S3ForcePathStyle = aws.Bool(true)

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
	basePath := os.Getenv("GITECHO_DATA_PATH")
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
		relPath, err := filepath.Rel(basePath, path)
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

func (s *S3Storage) DownloadDirectory(s3Path, localPath string) error {
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
		Prefix: aws.String(s3Path),
	}

	err := s.S3Client.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				// Construct the local file path
				relPath, err := filepath.Rel(s3Path, *obj.Key)
				if err != nil {
					fmt.Printf("Failed to determine the relative path: %v\n", err)
					continue
				}
				filePath := filepath.Join(localPath, relPath)

				// Ensure the directory path exists
				err = os.MkdirAll(filepath.Dir(filePath), 0755)
				if err != nil {
					fmt.Printf("Failed to create directory: %v\n", err)
					continue
				}

				// Download the file
				err = s.downloadFile(*obj.Key, filePath)
				if err != nil {
					fmt.Printf("Failed to download file: %v\n", err)
				}
			}
			return true
		})

	if err != nil {
		return fmt.Errorf("failed to list objects: %v", err)
	}

	return nil
}

func (s *S3Storage) downloadFile(s3Key, filePath string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(s3Key),
	}

	output, err := s.S3Client.GetObject(input)
	if err != nil {
		return fmt.Errorf("failed to get S3 object: %v", err)
	}
	defer output.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, output.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
