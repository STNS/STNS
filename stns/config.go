package stns

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/STNS/STNS/model"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-yaml/yaml"
)

type decoder interface {
	decode(string, *Config) error
}

type tomlDecoder struct{}
type yamlDecoder struct{}

func (t *tomlDecoder) decode(path string, conf *Config) error {
	_, err := toml.DecodeFile(path, conf)
	return err
}

func (y *yamlDecoder) decode(path string, conf *Config) error {
	_, err := toml.DecodeFile(path, conf)

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, conf)
	return err
}

func decode(path string, conf *Config) error {
	var d decoder
	d = new(tomlDecoder)
	if strings.HasSuffix(path, "yaml") || strings.HasSuffix(path, "yml") {
		d = new(yamlDecoder)
	}
	return d.decode(path, conf)
}

func NewConfig(confPath string) (Config, error) {
	var conf Config

	if strings.HasPrefix(confPath, "s3:") {
		u, err := url.Parse(confPath)
		if err != nil {
			return conf, err
		}

		if u.Host == "" || u.Path == "" {
			return conf, errors.New("Bucket name and path are required to use S3")
		}

		client := s3.New(session.New(), nil)
		res, err := client.GetObject(&s3.GetObjectInput{
			Bucket: &u.Host,
			Key:    &u.Path,
		})
		defer res.Body.Close()
		if err != nil {
			return conf, err
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return conf, err
		}

		tmpFile, err := ioutil.TempFile(os.TempDir(), "stns-")
		if err != nil {
			return conf, err
		}
		defer os.Remove(tmpFile.Name())

		if _, err = tmpFile.Write(body); err != nil {
			return conf, err
		}
		confPath = tmpFile.Name()
	} else {
		conf.dir = filepath.Dir(confPath)
	}
	defaultConfig(&conf)

	if err := decode(confPath, &conf); err != nil {
		return conf, err
	}

	if conf.Include != "" {
		if err := includeConfigFile(&conf, conf.Include); err != nil {
			return Config{}, err
		}
	}

	if conf.dir != "" && !strings.HasPrefix(conf.ModulePath, "/") {
		conf.ModulePath = filepath.Join(conf.dir, conf.ModulePath)
	}

	return conf, nil
}

type Config struct {
	dir       string
	Port      int        `toml:"port"`
	BasicAuth *BasicAuth `toml:"basic_auth" yaml:"basic_auth"`
	TokenAuth *TokenAuth `toml:"token_auth" yaml:"token_auth"`

	UseServerStarter bool
	Users            *model.Users
	Groups           *model.Groups
	Include          string   `toml:"include"`
	ModulePath       string   `toml:"module_path" yaml:"module_path"`
	LoadModules      []string `toml:"load_modules" yaml:"load_modules"`
	Modules          map[string]interface{}
	TLS              *TLS
	LDAP             *LDAP
}

type LDAP struct {
	BaseDN string
}
type TLS struct {
	CA   string
	Cert string
	Key  string
}
type BasicAuth struct {
	User     string
	Password string
}
type TokenAuth struct {
	Tokens []string
}

func defaultConfig(c *Config) {
	c.Port = 1104
	c.ModulePath = "/usr/local/stns/modules.d"
	c.LDAP = &LDAP{
		BaseDN: "dc=stns,dc=local",
	}
}

func includeConfigFile(config *Config, include string) error {
	if !strings.HasPrefix(include, "/") {
		include = filepath.Join(config.dir, include)
	}

	files, err := filepath.Glob(include)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := decode(file, config); err != nil {
			return fmt.Errorf("while loading included config file %s: %s", file, err)
		}
	}
	return nil
}
