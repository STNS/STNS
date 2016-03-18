package attribute

import (
	"reflect"
	"strconv"
)

type All struct {
	Id int `toml:"id" json:"id"`
	// use user
	GroupId   int      `toml:"group_id" json:"group_id"`
	Directory string   `toml:"directory" json:"directory"`
	Shell     string   `toml:"shell" json:"shell"`
	Gecos     string   `toml:"gecos" json:"gecos"`
	Keys      []string `toml:"keys" json:"keys"`
	LinkUsers []string `toml:"link_users" json:"link_users"`
	// use group
	Users []string `toml:"users" json:"users"`
}

type UserGroups map[string]*All

func (u UserGroups) GetByName(name string) UserGroups {
	attr := u[name]
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		return nil
	}
	return UserGroups{
		name: attr,
	}
}

func (u UserGroups) GetById(_id string) UserGroups {
	id, _ := strconv.Atoi(_id)
	for k, u := range u {
		if u.Id == id {
			return UserGroups{
				k: u,
			}
		}
	}
	return nil
}

func (u UserGroups) Merge(m1 UserGroups) {
	for i, v := range m1 {
		u[i] = v
	}
}
