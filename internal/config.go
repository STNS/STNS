package palk

import "github.com/BurntSushi/toml"

const filePath = "/etc/palk.conf"

var AllConfig *Config

type Config struct {
	Users  map[string]*Attr
	Groups map[string]*Attr
}

type UserAttr struct {
	Group_Id  int
	Directory string
	Shell     string
}
type GroupAttr struct {
	Users []string
}
type Attr struct {
	Id   int
	Name string
	*UserAttr
	*GroupAttr
}

func LoadConfig() {
	_, err := toml.DecodeFile(filePath, &AllConfig)
	if err != nil {
		panic(err)
	}
}
