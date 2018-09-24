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
	FindUserByID(int) (map[string]UserGroup, error)
	FindUserByName(string) (map[string]UserGroup, error)
	FindGroupByID(int) (map[string]UserGroup, error)
	FindGroupByName(string) (map[string]UserGroup, error)
	Users() (map[string]UserGroup, error)
	Groups() (map[string]UserGroup, error)
	HighestUserID() int
	LowestUserID() int
	HighestGroupID() int
	LowestGroupID() int
}
