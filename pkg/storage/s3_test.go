package storage_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
)

// MockS3Client is a mock implementation of the s3iface.S3API interface
type MockS3Client struct {
	s3iface.S3API
	MockListObjectsV2Pages func(*s3.ListObjectsV2Input, func(*s3.ListObjectsV2Output, bool) bool) error
	MockGetObject          func(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	// Implement your desired behavior here
	return nil, nil
}

func (m *MockS3Client) ListObjectsV2Pages(params *s3.ListObjectsV2Input, callback func(page *s3.ListObjectsV2Output, lastPage bool) bool) error {
	if m.MockListObjectsV2Pages != nil {
		return m.MockListObjectsV2Pages(params, callback)
	}
	return nil
}

func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if m.MockGetObject != nil {
		return m.MockGetObject(input)
	}
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

func TestDownloadDirectory(t *testing.T) {
	// Create a temporary directory for the downloaded files
	tempDir, err := os.MkdirTemp("", "test-download")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock S3 client
	mockS3Client := &MockS3Client{}

	// Set up the mock behavior for ListObjectsV2Pages
	mockObjects := []*s3.Object{
		{
			Key: aws.String("path/to/file1.txt"),
		},
		{
			Key: aws.String("path/to/file2.txt"),
		},
	}
	mockS3Client.MockListObjectsV2Pages = func(params *s3.ListObjectsV2Input, callback func(page *s3.ListObjectsV2Output, lastPage bool) bool) error {
		callback(&s3.ListObjectsV2Output{
			Contents: mockObjects,
		}, true)
		return nil
	}

	// Set up the mock behavior for GetObject
	mockS3Client.MockGetObject = func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		object := &s3.GetObjectOutput{
			Body: io.NopCloser(strings.NewReader("test data")),
		}
		return object, nil
	}

	s3Storage := &storage.S3Storage{
		Session:    nil,
		BucketName: "test-bucket",
		S3Client:   mockS3Client,
	}

	// Call the function under test
	err = s3Storage.DownloadDirectory("path/to", tempDir)
	if err != nil {
		t.Fatalf("DownloadDirectory returned an error: %v", err)
	}

	// Verify the downloaded files
	file1Path := filepath.Join(tempDir, "file1.txt")
	file2Path := filepath.Join(tempDir, "file2.txt")
	if _, err := os.Stat(file1Path); os.IsNotExist(err) {
		t.Fatalf("File %s does not exist", file1Path)
	}
	if _, err := os.Stat(file2Path); os.IsNotExist(err) {
		t.Fatalf("File %s does not exist", file2Path)
	}
}
