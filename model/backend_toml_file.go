package model

import "sort"

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

func (t BackendTomlFile) Users() map[string]UserGroup {
	return t.users.ToUserGroup()
}

func (t BackendTomlFile) FindGroupByID(id int) map[string]UserGroup {
	return tomlFileFindByID(id, t.groups.ToUserGroup())
}

func (t BackendTomlFile) FindGroupByName(name string) map[string]UserGroup {
	return tomlFileFindByName(name, t.groups.ToUserGroup())
}

func (t BackendTomlFile) Groups() map[string]UserGroup {
	return t.groups.ToUserGroup()
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

// user or group
type linkAttributer interface {
	linkValues() []string
	setLinkValues([]string)
	value() []string
	name() string
}

type linkAttributers map[string]linkAttributer

func (las linkAttributers) find(keys []string) linkAttributers {
	result := linkAttributers{}
	if las != nil {
		for _, key := range keys {
			for _, lv := range las {
				if lv.name() == key {
					result[key] = lv
				}
			}
		}
	}
	return result
}

// Userは公開鍵をlink_usersから取得して、自身の鍵としてマージする
// Groupはグループのメンバーをlink_groupから取得して、自身のメンバーとしてマージする
func mergeLinkAttribute(master, current linkAttributers, result map[string][]string, nest *int) map[string][]string {
	if current == nil {
		current = master
	}

	if result == nil {
		result = map[string][]string{}
	}

	if nest == nil {
		i := 0
		nest = &i
	}

	for _, v := range current {
		links := v.linkValues()
		if len(links) > 0 {
			ls := master.find(links)
			// user3をrangeする
			for _, iv := range ls {
				if len(result[iv.name()]) == 0 {
					result[iv.name()] = append(result[iv.name()], iv.value()...)
					*nest++
					result = mergeLinkAttribute(master, ls, result, nest)
					*nest--
				}
			}
			if *nest == 0 {
				for _, rv := range result {
					v.setLinkValues(rv)
				}
				result = map[string][]string{}
			}
		}
	}
	return result
}

func mapSliceToSlice(m map[string][]string) []string {
	var result []string
	for _, v := range m {
		result = append(result, v...)
	}
	return result
}

func isStringsExist(n string, xs []string) bool {
	for _, x := range xs {
		if n == x {
			return true
		}
	}
	return false
}

func uniqStrings(xs []string) []string {
	ys := make([]string, 0, len(xs))
	for _, x := range xs {
		if !isStringsExist(x, ys) {
			ys = append(ys, x)
		}
	}
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
	return ys
}
