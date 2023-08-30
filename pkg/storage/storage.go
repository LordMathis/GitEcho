package storage

import "gopkg.in/yaml.v3"

type Storage interface {
	UploadDirectory(directoryPath string) error
	DownloadDirectory(remotePath, localPath string) error
}

type BaseStorage struct {
	Storage
	Name   string  `yaml:"name"`
	Type   string  `yaml:"type"`
	Config Storage `yaml:"config"`
}

func (b *BaseStorage) UnmarshalYAML(value *yaml.Node) error {

	var t struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
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
			Config *S3StorageConfig `yaml:"config"`
		}

		err := value.Decode(&c)
		if err != nil {
			return err
		}

		c.Config.InitializeS3Storage()
		b.Config = c.Config
	}

	return nil
}
