package model

type BackendTomlFile struct {
	users  *Users
	groups *Groups
}

func NewBackendTomlFile(u *Users, g *Groups) *BackendTomlFile {
	if u != nil {
		ensureName(u.ToUserGroup())
	}

	if g != nil {
		ensureName(g.ToUserGroup())
	}

	return &BackendTomlFile{
		users:  u,
		groups: g,
	}
}

func (t BackendTomlFile) FindUserByID(id int) map[string]UserGroup {
	return tomlFileFindByID(id, t.users.ToUserGroup())
}

func (t BackendTomlFile) FindUserByName(name string) map[string]UserGroup {
	return tomlFileFindByName(name, t.users.ToUserGroup())
}

func tomlFileFindByID(id int, list map[string]UserGroup) map[string]UserGroup {
	res := map[string]UserGroup{}
	if list != nil {
		for k, v := range list {
			if id == v.id() {
				res[k] = v
			}
		}
	}
	return res
}

func tomlFileFindByName(name string, list map[string]UserGroup) map[string]UserGroup {
	res := map[string]UserGroup{}
	if list != nil {
		for k, v := range list {
			if name == v.name() {
				res[k] = v
			}
		}
	}
	return res
}

func ensureName(list map[string]UserGroup) {
	if list != nil {
		for k, v := range list {
			v.setName(k)
		}
	}
}
