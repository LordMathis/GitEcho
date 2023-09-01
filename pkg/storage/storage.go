package storage

import (
	"gopkg.in/yaml.v3"
)

type Uploader interface {
	UploadDirectory(directoryPath string) error
}

type Downloader interface {
	DownloadDirectory(remotePath, localPath string) error
}

type Initializer interface {
	Initialize() error
}

type Storage interface {
	Uploader
	Downloader
	Initializer
}

type BaseStorage struct {
	Name   string  `yaml:"name"`
	Type   string  `yaml:"type"`
	Config Storage `yaml:"config"`
}

func (b *BaseStorage) UnmarshalYAML(value *yaml.Node) error {

	var t struct {
		Name   string    `yaml:"name"`
		Type   string    `yaml:"type"`
		Config yaml.Node `yaml:"config"`
	}

	err := value.Decode(&t)
	if err != nil {
		return err
	}

	b.Name = t.Name
	b.Type = t.Type

	switch b.Type {
	case "s3":
		var c struct {
			Config S3StorageConfig `yaml:"config"`
		}

		err := value.Decode(&c)
		if err != nil {
			return err
		}

		b.Config = &c.Config
	}

	return nil
}
