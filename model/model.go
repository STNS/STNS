package model

type UserGroup interface {
	id() int
	name() string
	setName(string)
}

func FindByID(id int, list map[string]interface{}) map[string]UserGroup {
	res := map[string]UserGroup{}
	if list != nil {
		for k, v := range list {
			l, ok := v.(UserGroup)
			if ok {
				if id == l.id() {
					res[k] = l
				}
			}
		}
	}
	return res
}

func EnsureName(list map[string]interface{}) {
	if list != nil {
		for k, v := range list {
			l, ok := v.(UserGroup)
			if ok {
				l.setName(k)
			}
		}
	}
}
