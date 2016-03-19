package attribute

import (
	"reflect"
	"strconv"
)

type All struct {
	Id int `toml:"id" json:"id"`
	*User
	*Group
}

type User struct {
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

type UserGroups map[string]*All

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
