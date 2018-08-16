package model

type UserGroup interface {
	id() int
	name() string
	setName(string)
}

func FindByID(id int, list map[string]UserGroup) map[string]UserGroup {
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

func EnsureName(list map[string]UserGroup) {
	if list != nil {
		for k, v := range list {
			v.setName(k)
		}
	}
}
