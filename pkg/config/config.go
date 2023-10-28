package config

import (
	"os"

	"github.com/LordMathis/GitEcho/pkg/repository"
	"gopkg.in/yaml.v3"

	"github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
)

type Config struct {
	DataPath         string                            `yaml:"data_path"`
	RcloneConfigPath string                            `yaml:"rclone_config_path"`
	Repositories     map[string]*repository.BackupRepo `yaml:"repositories"`
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	var t struct {
		DataPath         string                   `yaml:"data_path"`
		RcloneConfigPath string                   `yaml:"rclone_config_path"`
		Repositories     []*repository.BackupRepo `yaml:"repositories"`
	}

	err := value.Decode(&t)
	if err != nil {
		return err
	}

	c.DataPath = t.DataPath
	c.Repositories = make(map[string]*repository.BackupRepo)

	if t.RcloneConfigPath != "" {
		config.SetConfigPath(t.RcloneConfigPath)
		configfile.Install()
	}

	for _, repo := range t.Repositories {
		repo.LocalPath = c.DataPath + "/" + repo.Name
		repo.Initialize()
		c.Repositories[repo.Name] = repo
		for _, stor := range repo.Storages {
			err = stor.InitializeStorage()
			if err != nil {
				return err
			}
		}
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
