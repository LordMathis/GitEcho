package storage_test

import (
	"testing"

	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
)

// MockS3Client is a mock implementation of the s3iface.S3API interface
type MockS3Client struct {
	s3iface.S3API
}

func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	// Implement your desired behavior here
	return nil, nil
}

func TestUploadDirectory(t *testing.T) {
	// Create a mock S3 client
	mockS3Client := &MockS3Client{}

	// Create the S3Storage instance with the mock S3 client
	sess, err := session.NewSession()
	assert.NoError(t, err)

	s3Storage := &storage.S3Storage{
		Session:    sess,
		BucketName: "test-bucket",
		S3Client:   mockS3Client,
	}

	// Call the function under test
	err = s3Storage.UploadDirectory("./testdata")
	assert.NoError(t, err)
}
