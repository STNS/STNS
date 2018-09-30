package model

type UserGroup interface {
	GetID() int
	GetName() string
	setName(string)
	linkValues() []string
	setLinkValues([]string)
	value() []string
}

type Backend interface {
	GetterBackend
	SetterBackend
}

type GetterBackends []GetterBackend
type GetterBackend interface {
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

type SetterBackend interface {
	Create(string, UserGroup) error
}

func mergeUserGroup(m1, m2 map[string]UserGroup) map[string]UserGroup {
	ans := map[string]UserGroup{}

	for k, v := range m1 {
		ans[k] = v
	}
	for k, v := range m2 {
		ans[k] = v
	}
	return (ans)
}
