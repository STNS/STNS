package model

type UserGroup interface {
	id() int
	name() string
	setName(string)
}

type Backend interface {
	FindUserByID(int) map[string]UserGroup
	FindUserByName(string) map[string]UserGroup
}
