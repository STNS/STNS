package model

type UserGroup interface {
	id() int
	name() string
	setName(string)
	linkValues() []string
	setLinkValues([]string)
	value() []string
}

type Backend interface {
	FindUserByID(int) map[string]UserGroup
	FindUserByName(string) map[string]UserGroup
	FindGroupByID(int) map[string]UserGroup
	FindGroupByName(string) map[string]UserGroup
	Users() map[string]UserGroup
	Groups() map[string]UserGroup
	HighestUserID() int
	LowestUserID() int
	HighestGroupID() int
	LowestGroupID() int
}
