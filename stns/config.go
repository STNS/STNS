package stns

import "github.com/BurntSushi/toml"

func NewConfig(confPath string) (Config, error) {
	var conf Config
	defaultConfig(&conf)

	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

type Config struct {
	UseServerStarter bool
}

func defaultConfig(c *Config) {
}
