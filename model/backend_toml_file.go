package model

import (
	"fmt"
	"sort"

	validator "gopkg.in/go-playground/validator.v9"
)

type BackendTomlFile struct {
	users  *Users
	groups *Groups
}

const Highest = 0
const Lowest = 1

func NewBackendTomlFile(u *Users, g *Groups) (*BackendTomlFile, error) {
	if u != nil {
		ug := u.ToUserGroup()
		ensureName(ug)
		mergeLinkAttribute(ug, nil, nil, nil)

		if err := checkDuplicateID(ug); err != nil {
			return nil, err
		}
		if err := validateUserGroup(ug); err != nil {
			return nil, err
		}
	}

	if g != nil {
		gg := g.ToUserGroup()
		ensureName(gg)
		mergeLinkAttribute(gg, nil, nil, nil)
		if err := checkDuplicateID(gg); err != nil {
			return nil, err
		}
		if err := validateUserGroup(gg); err != nil {
			return nil, err
		}
	}

	return &BackendTomlFile{
		users:  u,
		groups: g,
	}, nil
}

func (t BackendTomlFile) FindUserByID(id int) (map[string]UserGroup, error) {
	r, e := tomlFileFindByID(id, t.users.ToUserGroup())
	return errorHandler(r, e, id, "user")
}

func (t BackendTomlFile) FindUserByName(name string) (map[string]UserGroup, error) {
	r, e := tomlFileFindByName(name, t.users.ToUserGroup())
	return errorHandler(r, e, name, "user")
}

func (t BackendTomlFile) Users() (map[string]UserGroup, error) {
	r := t.users.ToUserGroup()
	return errorHandler(r, nil, nil, "user")
}

func (t BackendTomlFile) FindGroupByID(id int) (map[string]UserGroup, error) {
	r, e := tomlFileFindByID(id, t.groups.ToUserGroup())
	return errorHandler(r, e, id, "group")
}

func (t BackendTomlFile) FindGroupByName(name string) (map[string]UserGroup, error) {
	r, e := tomlFileFindByName(name, t.groups.ToUserGroup())
	return errorHandler(r, e, name, "group")
}

func (t BackendTomlFile) Groups() (map[string]UserGroup, error) {
	r := t.groups.ToUserGroup()
	return errorHandler(r, nil, nil, "group")
}

func (t BackendTomlFile) HighestUserID() int {
	return tomlHighLowID(Highest, t.users.ToUserGroup())
}

func (t BackendTomlFile) LowestUserID() int {
	return tomlHighLowID(Lowest, t.users.ToUserGroup())
}

func (t BackendTomlFile) HighestGroupID() int {
	return tomlHighLowID(Highest, t.groups.ToUserGroup())
}

func (t BackendTomlFile) LowestGroupID() int {
	return tomlHighLowID(Lowest, t.groups.ToUserGroup())
}

func tomlFileFindByID(id int, list map[string]UserGroup) (map[string]UserGroup, error) {
	res := map[string]UserGroup{}
	if list != nil {
		for k, v := range list {
			if id == v.GetID() {
				res[k] = v
			}
		}
	}

	return res, nil
}

func tomlFileFindByName(name string, list map[string]UserGroup) (map[string]UserGroup, error) {
	res := map[string]UserGroup{}
	if list != nil {
		for k, v := range list {
			if name == v.GetName() {
				res[k] = v
			}
		}
	}
	return res, nil
}

func ensureName(list map[string]UserGroup) {
	if list != nil {
		for k, v := range list {
			v.setName(k)
		}
	}
}

type linkAttributers map[string]UserGroup

func (las linkAttributers) find(keys []string) linkAttributers {
	result := linkAttributers{}
	if las != nil {
		for _, key := range keys {
			for _, lv := range las {
				if lv.GetName() == key {
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
			for _, iv := range ls {
				if len(result[iv.GetName()]) == 0 {
					result[iv.GetName()] = append(result[iv.GetName()], iv.value()...)
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

// highest=0 lowest= 1
func tomlHighLowID(highorlow int, list map[string]UserGroup) int {
	current := 0
	if list != nil {
		for _, v := range list {
			if current == 0 || (highorlow == 0 && current < v.GetID()) || (highorlow == 1 && current > v.GetID()) {
				current = v.GetID()
			}
		}
	}
	return current
}

func checkDuplicateID(attr map[string]UserGroup) error {
	b := map[int]bool{}

	for _, a := range attr {
		if a.GetID() != 0 && b[a.GetID()] {
			return fmt.Errorf("Duplicate id is not allowed: %d", a.GetID())
		}
		b[a.GetID()] = true
	}
	return nil
}

func validateUserGroup(attr map[string]UserGroup) error {
	validate := validator.New()
	for _, a := range attr {
		err := validate.Struct(a)
		if err != nil {
			return err
		}
	}
	return nil
}
