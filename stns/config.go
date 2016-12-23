package stns

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/BurntSushi/toml"
)

// Config config object
type Config struct {
	Port     int    `toml:"port"`
	Include  string `toml:"include"`
	TLSCa    string `toml:"tls_ca"`
	TLSCert  string `toml:"tls_cert"`
	TLSKey   string `toml:"tls_key"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Users    Attributes
	Groups   Attributes
	Sudoers  Attributes
}

var minUserID int
var minGroupID int

// LoadConfig from /etc/stns/stns.conf
func LoadConfig(configFile string) (Config, error) {
	var config Config
	defaultConfig(&config)

	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	if config.Include != "" {
		if err := includeConfigFile(&config, config.Include); err != nil {
			return Config{}, err
		}
	}

	type setup struct {
		t     string
		attr  Attributes
		minID *int
	}

	s := []setup{
		setup{
			t:     "user",
			attr:  config.Users,
			minID: &minUserID,
		},
		setup{
			t:     "group",
			attr:  config.Groups,
			minID: &minGroupID,
		},
	}

	for _, r := range s {
		setMinID(r.minID, r.attr)
		mergeLinkAttribute(r.t, r.attr)
		if err := checkDuplicateId(r.t, r.attr); err != nil {
			return Config{}, err
		}
	}

	return config, nil
}

func defaultConfig(config *Config) {
	config.Port = 1104
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

func setMinID(min *int, attrs Attributes) {
	*min = 0
	if len(attrs) > 0 {
		for _, a := range attrs {
			switch {
			case *min == 0:
				*min = a.ID
			case *min > a.ID:
				*min = a.ID
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
				linker.SetLinkValue(removeDuplicates(mergeValue))
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
			for _, nextName := range linker.LinkParams() {
				recursiveSetLinkValue(attr, rtype, nextName, result)
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

func removeDuplicates(xs []string) []string {
	ys := make([]string, 0, len(xs))
	for _, x := range xs {
		if !member(x, ys) {
			ys = append(ys, x)
		}
	}
	return ys
}

func checkDuplicateId(t string, attr Attributes) error {
	b := map[int]bool{}

	for _, a := range attr {
		if a.ID != 0 && b[a.ID] {
			return fmt.Errorf("Duplicate id is not allowed %s_id:%d", t, a.ID)
		}
		b[a.ID] = true
	}
	return nil
}
