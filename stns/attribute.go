package stns

import (
	"reflect"
	"strconv"
)

type Attribute struct {
	Id int `toml:"id" json:"id"`
	*User
	*Group
}

type User struct {
	Password  string   `toml:"password" json:"password"`
	HashType  string   `toml:"hash_type" json:"hash_type"`
	GroupId   int      `toml:"group_id" json:"group_id"`
	Directory string   `toml:"directory" json:"directory"`
	Shell     string   `toml:"shell" json:"shell"`
	Gecos     string   `toml:"gecos" json:"gecos"`
	Keys      []string `toml:"keys" json:"keys"`
	LinkUsers []string `toml:"link_users" json:"link_users"`
}

type Group struct {
	Users      []string `toml:"users" json:"users"`
	LinkGroups []string `toml:"link_groups" json:"link_groups"`
}

type Attributes map[string]*Attribute

type Linker interface {
	LinkTargetColumnValue() []string
	LinkValues() []string
	SetLinkValue([]string)
}

func (u *User) LinkTargetColumnValue() []string {
	return u.LinkUsers
}

func (u *User) LinkValues() []string {
	return u.Keys
}

func (u *User) SetLinkValue(v []string) {
	u.Keys = v
}

func (g *Group) LinkTargetColumnValue() []string {
	return g.LinkGroups
}

func (g *Group) LinkValues() []string {
	return g.Users
}

func (g *Group) SetLinkValue(v []string) {
	g.Users = v
}

func (u Attributes) GetByName(name string) Attributes {
	attr := u[name]
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		return nil
	}
	return Attributes{
		name: attr,
	}
}

func (u Attributes) GetById(_id string) Attributes {
	id, _ := strconv.Atoi(_id)
	for k, u := range u {
		if u.Id == id {
			return Attributes{
				k: u,
			}
		}
	}
	return nil
}

func (u Attributes) Merge(m1 Attributes) {
	for i, v := range m1 {
		u[i] = v
	}
}
