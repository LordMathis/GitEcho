package config

import (
	"os"

	"github.com/LordMathis/GitEcho/pkg/backuprepo"
	"github.com/LordMathis/GitEcho/pkg/storage"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DataPath     string                            `yaml:"data_path"`
	Repositories map[string]*backuprepo.BackupRepo `yaml:"repositories"`
	Storages     map[string]*storage.BaseStorage   `yaml:"storages"`
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	var t struct {
		DataPath     string                   `yaml:"data_path"`
		Repositories []*backuprepo.BackupRepo `yaml:"repositories"`
		Storages     []*storage.BaseStorage   `yaml:"storages"`
	}

	err := value.Decode(&t)
	if err != nil {
		return err
	}

	c.DataPath = t.DataPath
	c.Repositories = make(map[string]*backuprepo.BackupRepo)
	c.Storages = make(map[string]*storage.BaseStorage)

	for _, repo := range t.Repositories {
		c.Repositories[repo.Name] = repo
	}

	for _, storage := range t.Storages {
		c.Storages[storage.Name] = storage
	}

	return nil
}

func ReadConfig(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	yaml.Unmarshal(data, &config)

	return config, nil
}
