package stns

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port       int    `toml:"port"`
	Include    string `toml:"include"`
	Salt       bool   `toml:"salt_enable"`
	Stretching int    `toml:"stretching_number"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	HashType   string `toml:"hash_type" json:"hash_type"`
	Backend    Backend
	Users      Attributes
	Groups     Attributes
	Sudoers    Attributes
}

type Backend struct {
	Driver   string `toml:"driver"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user" validate:"required"`
	Password string `toml:"password"`
}

var MinUserId, MinGroupId int

func LoadConfig(configFile string) (Config, error) {
	var config Config
	defaultConfig(&config)

	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	if config.Include != "" {
		if err := includeConfigFile(&config, config.Include); err != nil {
			fmt.Println(err)
			return Config{}, err
		}
	}
	setMinId(&MinUserId, config.Users)
	setMinId(&MinGroupId, config.Groups)
	mergeLinkAttribute("user", config.Users)
	mergeLinkAttribute("group", config.Groups)
	return config, nil
}

func defaultConfig(config *Config) {
	config.Port = 1104
	config.Salt = false
	config.Stretching = 0
	config.HashType = "sha256"
	config.Backend = Backend{
		Host: "localhost",
		Port: "3306",
	}
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

func setMinId(min *int, attrs Attributes) {
	*min = 0
	if len(attrs) > 0 {
		for _, a := range attrs {
			switch {
			case *min == 0:
				*min = a.Id
			case *min > a.Id:
				*min = a.Id
			}
		}
	}
}

func mergeLinkAttribute(rtype string, attr Attributes) {
	for k, v := range attr {
		mergeValue := []string{}
		linker := getLinker(rtype, v)

		if linker != nil && !reflect.ValueOf(linker).IsNil() &&
			linker.LinkParams() != nil && !reflect.ValueOf(linker.LinkParams()).IsNil() {
			for _, linkValue := range linker.LinkParams() {
				linkValues := map[string][]string{k: linker.LinkValue()}

				recursiveSetLinkValue(attr, rtype, linkValue, linkValues)
				for _, val := range linkValues {
					mergeValue = append(mergeValue, val...)
				}
				linker.SetLinkValue(RemoveDuplicates(mergeValue))
			}
		}
	}
}

func getLinker(rtype string, attr *Attribute) Linker {
	if attr != nil && !reflect.ValueOf(attr).IsNil() {
		if rtype == "user" {
			return attr.User
		} else if rtype == "group" {
			return attr.Group
		}
	}
	return nil
}

func recursiveSetLinkValue(attr Attributes, rtype, name string, result map[string][]string) {
	if result[name] != nil {
		return
	}

	linker := getLinker(rtype, attr[name])

	if linker != nil && !reflect.ValueOf(linker).IsNil() && len(linker.LinkValue()) > 0 {
		result[name] = linker.LinkValue()
		if linker.LinkParams() != nil || !reflect.ValueOf(linker.LinkParams()).IsNil() {
			for _, next_name := range linker.LinkParams() {
				recursiveSetLinkValue(attr, rtype, next_name, result)
			}
		}
	}
}

func member(n string, xs []string) bool {
	for _, x := range xs {
		if n == x {
			return true
		}
	}
	return false
}

func RemoveDuplicates(xs []string) []string {
	ys := make([]string, 0, len(xs))
	for _, x := range xs {
		if !member(x, ys) {
			ys = append(ys, x)
		}
	}
	return ys
}
