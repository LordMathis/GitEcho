package config

import (
	"os"

	"github.com/LordMathis/GitEcho/pkg/repository"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DataPath     string                            `yaml:"data_path"`
	Repositories map[string]*repository.BackupRepo `yaml:"repositories"`
	Storages     map[string]*storage.Storage       `yaml:"storages"`
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	var t struct {
		DataPath     string                   `yaml:"data_path"`
		Repositories []*repository.BackupRepo `yaml:"repositories"`
		Storages     []*storage.Storage       `yaml:"storages"`
	}

	err := value.Decode(&t)
	if err != nil {
		return err
	}

	c.DataPath = t.DataPath
	c.Repositories = make(map[string]*repository.BackupRepo)
	c.Storages = make(map[string]*storage.Storage)

	for _, repo := range t.Repositories {
		repo.LocalPath = c.DataPath + "/" + repo.Name
		repo.Initialize()
		c.Repositories[repo.Name] = repo
	}

	for _, stor := range t.Storages {
		err = stor.InitializeStorage()
		if err != nil {
			return err
		}
		c.Storages[stor.RemoteName] = stor
	}

	return nil
}

func ReadConfig(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfigFile(data)
}

func ParseConfigFile(data []byte) (*Config, error) {
	config := &Config{}
	yaml.Unmarshal(data, &config)

	return config, nil
}
