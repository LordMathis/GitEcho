package storage

import "fmt"

type Storage interface {
	UploadDirectory(directoryPath string) error
}

type StorageCreator interface {
	CreateStorage(storageType, storageData string) (Storage, error)
}

type StorageCreatorImpl struct {
}

func (c *StorageCreatorImpl) CreateStorage(storageType, storageData string) (Storage, error) {
	var storageInstance Storage
	var err error

	switch storageType {
	case "s3":
		storageInstance, err = NewS3StorageFromJson(storageData)
		if err != nil {
			return nil, err
		}

		s3storage, ok := storageInstance.(*S3Storage)
		if !ok {
			return nil, fmt.Errorf("failed to retype storageInstance to S3Storage")
		}

		err = s3storage.DecryptKeys()
		if err != nil {
			return nil, err
		}

		return s3storage, nil

	default:
		return nil, fmt.Errorf("unknown storage type")
	}
}
