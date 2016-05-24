package stns

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port       int    `toml:"port"`
	Include    string `toml:"include"`
	Salt       bool   `toml:"salt_enable"`
	Stretching int    `toml:"stretching_count"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	HashType   string `toml:"hash_type" json:"hash_type"`
	Users      Attributes
	Groups     Attributes
	Sudoers    Attributes
}

func LoadConfig(configFile string) (*Config, error) {
	var config Config
	defaultConfig(&config)

	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, err
	}

	if config.Include != "" {
		if err := includeConfigFile(&config, config.Include); err != nil {
			return nil, err
		}
	}
	return &config, nil
}

func defaultConfig(config *Config) {
	config.Port = 1104
	config.Salt = false
	config.Stretching = 0
	config.HashType = "sha256"
}

func includeConfigFile(config *Config, include string) error {
	files, err := filepath.Glob(include)
	if err != nil {
		return err
	}

	for _, file := range files {
		_, err := toml.DecodeFile(file, &config)
		if err != nil {
			return fmt.Errorf("while loading included config file %s: %s", file, err)
		}
	}
	return nil
}
