package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3StorageConfig struct {
	Session        *session.Session `yaml:"-"`
	S3Client       s3iface.S3API    `yaml:"-"`
	Endpoint       string           `yaml:"endpoint"`
	Region         string           `yaml:"region"`
	AccessKey      string           `yaml:"access_key"`
	SecretKey      string           `yaml:"secret_key"`
	BucketName     string           `yaml:"bucket_name"`
	DisableSSL     bool             `yaml:"disable_ssl"`
	ForcePathStyle bool             `yaml:"force_path_style"`
	Encryption     EncryptionConfig `yaml:"encryption"`
}

type EncryptionConfig struct {
	Enabled bool   `yaml:"enabled"`
	Key     string `yaml:"key"`
}

func getSession(endpoint, region, accessKey, secretKey string, disableSSL, forcePathStyle bool) (*session.Session, error) {
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

	config.DisableSSL = aws.Bool(disableSSL)
	config.S3ForcePathStyle = aws.Bool(forcePathStyle)

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *S3StorageConfig) Initialize() error {
	session, err := getSession(s.Endpoint, s.Region, s.AccessKey, s.SecretKey, s.DisableSSL, s.ForcePathStyle)
	if err != nil {
		return err
	}

	svc := s3.New(session)

	s.Session = session
	s.S3Client = svc

	return nil
}

func (s *S3StorageConfig) UploadDirectory(directoryPath string) error {

	parent := filepath.Dir(directoryPath)

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
		relPath, err := filepath.Rel(parent, path)
		if err != nil {
			return err
		}
		key := filepath.Join(relPath)

		input := &s3.PutObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(key),
			Body:   f,
		}

		if s.Encryption.Enabled {
			// Perform client-side encryption before uploading
			encryptedData, err := encryption.EncryptData(f, []byte(s.Encryption.Key))
			if err != nil {
				return err
			}

			keyParts := strings.SplitN(key, "/", 2)

			encryptedKey, err := encryption.ScrambleString(keyParts[1], []byte(s.Encryption.Key))
			if err != nil {
				return err
			}

			input.Body = aws.ReadSeekCloser(encryptedData)
			input.Key = aws.String(fmt.Sprintf("%s/%s", keyParts[0], encryptedKey))
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

func (s *S3StorageConfig) DownloadDirectory(remotePath, localPath string) error {
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
		Prefix: aws.String(remotePath + "/"),
	}

	err := s.S3Client.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				// Construct the local file path
				remoteRelPath := *obj.Key

				if s.Encryption.Enabled {

					keyParts := strings.SplitN(*obj.Key, "/", 2)

					decryptedKey, err := encryption.UnscrambleString(keyParts[1], []byte(s.Encryption.Key))
					if err != nil {
						fmt.Printf("Failed to decrypt object key: %v\n", err)
						continue
					}

					remoteRelPath = filepath.Join(keyParts[0], decryptedKey)
				}

				relPath, err := filepath.Rel(remotePath, remoteRelPath)
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

func (s *S3StorageConfig) downloadFile(s3Key, filePath string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(s3Key),
	}

	output, err := s.S3Client.GetObject(input)
	if err != nil {
		return fmt.Errorf("failed to get S3 object: %v", err)
	}
	defer output.Body.Close()

	var reader io.Reader = output.Body
	if s.Encryption.Enabled {
		decryptedData, err := encryption.DecryptData(output.Body, []byte(s.Encryption.Key))
		if err != nil {
			return fmt.Errorf("failed to decrypt file: %v", err)
		}
		reader = decryptedData
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
