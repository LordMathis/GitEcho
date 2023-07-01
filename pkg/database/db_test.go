package database_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LordMathis/GitEcho/pkg/database"
	"github.com/LordMathis/GitEcho/pkg/encryption"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestInsertS3Storage(t *testing.T) {
	// Create a new mock DB and expect the necessary query and arguments
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	encryption.SetEncryptionKey([]byte("12345678901234567890123456789012"))

	mock.ExpectPrepare("INSERT INTO storage").
		ExpectQuery().
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Create a new instance of the Database and call the function
	database := &database.Database{DB: sqlxDB}
	s3Storage := &storage.S3Storage{
		Endpoint:   "test-endpoint",
		Region:     "test-region",
		AccessKey:  "test-access-key",
		SecretKey:  "test-secret-key",
		BucketName: "test-bucket",
	}

	storageID, err := database.InsertS3Storage(s3Storage)
	assert.NoError(t, err)
	assert.Equal(t, 1, storageID)

	// Assert that all the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
