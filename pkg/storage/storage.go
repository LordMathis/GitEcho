package storage

type Storage interface {
	UploadDirectory(directoryPath string) error
}
