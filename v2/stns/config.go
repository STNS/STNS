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
	"github.com/STNS/STNS/v2/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func downloadFromS3(path, key string) (*os.File, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	if u.Host == "" || u.Path == "" {
		return nil, errors.New("bucket name and path are required to use S3")
	}
	tmpDir := os.TempDir()
	tmpFile, err := ioutil.TempFile(tmpDir, "stns-")
	if err != nil {
		return nil, err
	}

	if key == "" {
		key = u.Path
	}
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(tmpFile, &s3.GetObjectInput{
		Bucket: aws.String(u.Host),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}

func NewConfig(confPath string) (Config, error) {
	var conf Config
	loadPath := confPath

	if strings.HasPrefix(confPath, "s3:") {
		f, err := downloadFromS3(confPath, "")
		if err != nil {
			return conf, err
		}
		defer os.Remove(f.Name())

		loadPath = f.Name()
	}
	conf.dir = filepath.Dir(loadPath)
	defaultConfig(&conf)

	if err := decode(loadPath, &conf); err != nil {
		return conf, err
	}

	if conf.Include != "" {
		if err := includeConfigFile(&conf, confPath, conf.Include); err != nil {
			return Config{}, err
		}
	}
	overwrideConfig(&conf)
	return conf, nil
}

type Config struct {
	dir       string
	Port      int        `toml:"port"`
	BasicAuth *BasicAuth `toml:"basic_auth" yaml:"basic_auth"`
	TokenAuth *TokenAuth `toml:"token_auth" yaml:"token_auth"`

	AllowIPs         []string `toml:"allow_ips" yaml:"allow_ips"`
	UseServerStarter bool
	Users            *model.Users
	Groups           *model.Groups
	Include          string `toml:"include"`
	ModulePath       string `toml:"module_path" yaml:"module_path"`
	LoadModule       string `toml:"load_module" yaml:"load_module"`
	Modules          map[string]interface{}
	TLS              *TLS
	LDAP             *LDAP
	Redis            *redis `toml:"redis"`
}

type LDAP struct {
	BaseDN string
}

type redis struct {
	Host     string
	Port     int
	Password string
	TTL      int
	DB       int
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

func overwrideConfig(c *Config) {
	if os.Getenv("STNS_BASIC_AUTH_USER") != "" || os.Getenv("STNS_BASIC_AUTH_PASSWORD") != "" {
		if c.BasicAuth == nil {
			c.BasicAuth = &BasicAuth{}
		}
		if os.Getenv("STNS_BASIC_AUTH_USER") != "" {
			c.BasicAuth.User = os.Getenv("STNS_BASIC_AUTH_USER")
		}

		if os.Getenv("STNS_BASIC_AUTH_PASSWORD") != "" {
			c.BasicAuth.Password = os.Getenv("STNS_BASIC_AUTH_PASSWORD")
		}
	}

	if c.Redis != nil {
		if os.Getenv("STNS_REDIS_PASSWORD") != "" {
			c.Redis.Password = os.Getenv("STNS_REDIS_PASSWORD")
		}
	}
	if os.Getenv("STNS_AUTH_TOKEN") != "" {
		if c.TokenAuth == nil {
			c.TokenAuth = &TokenAuth{}
		}
		if os.Getenv("STNS_AUTH_TOKEN") != "" {
			c.TokenAuth.Tokens = strings.Split(os.Getenv("STNS_AUTH_TOKEN"), ",")
		}
	}

	if os.Getenv("STNS_ETCD_PASSWORD") != "" && c.Modules["etcd"] != nil {
		c.Modules["etcd"].(map[string]interface{})["password"] = os.Getenv("STNS_ETCD_PASSWORD")
	}
}

func includeConfigFile(config *Config, confPath, include string) error {

	if strings.HasPrefix(confPath, "s3:") {
		if strings.HasPrefix(include, "/") {
			return errors.New("absolute path can not be used when using S3")
		}

		f, err := downloadFromS3(confPath, include)
		if err != nil {
			return err
		}
		defer os.Remove(f.Name())
		if err := decode(f.Name(), config); err != nil {
			return fmt.Errorf("while loading included config file %s: %s", f.Name(), err)
		}
	} else {
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
	}
	return nil
}
