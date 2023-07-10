package storage

import (
	"encoding/json"
	"fmt"
)

type Storage interface {
	UploadDirectory(directoryPath string) error
	DownloadDirectory(remotePath, localPath string) error
}

type BaseStorage struct {
	Name string `json:"name" db:"name"`
	Type string `json:"type" db:"type"`
	Data string `json:"data" db:"data"`
}

type StorageCreator interface {
	CreateStorage(storageData json.RawMessage) (Storage, error)
}

type StorageCreatorImpl struct {
}

func (c *StorageCreatorImpl) CreateStorage(storage json.RawMessage) (Storage, error) {
	var storageInstance Storage
	var err error

	var baseStorage BaseStorage
	err = json.Unmarshal(storage, &baseStorage)
	if err != nil {
		return nil, err
	}

	switch baseStorage.Type {
	case "s3":
		storageInstance, err = NewS3StorageFromJson(storage)
		if err != nil {
			return nil, err
		}

		return storageInstance, nil

	default:
		return nil, fmt.Errorf("unknown storage type")
	}
}
