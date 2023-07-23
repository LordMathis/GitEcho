package storage

import "fmt"

type Storage interface {
	UploadDirectory(directoryPath string) error
	DownloadDirectory(remotePath, localPath string) error
	GetName() string
	GetType() StorageType
}

type StorageType string

type BaseStorage struct {
	Name string      `json:"name" db:"name"`
	Type StorageType `json:"type" db:"type"`
	Data string      `json:"data" db:"data"`
}

const S3StorageType StorageType = "s3"

func CreateStorage(baseStorage BaseStorage) (Storage, error) {
	var storageInstance Storage
	var err error

	switch baseStorage.Type {
	case S3StorageType:
		storageInstance, err = NewS3StorageFromBase(baseStorage)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown storage type")
	}

	return storageInstance, nil
}
