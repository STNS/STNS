package stns

import (
	"sort"
	"strconv"
)

// Attribute attribute object
type Attribute struct {
	ID int `toml:"id" json:"id"`
	*User
	*Group
}

// Attributes attributes object
type Attributes map[string]*Attribute

// Linker linker interface
type Linker interface {
	LinkParams() []string
	LinkValue() []string
	SetLinkValue([]string)
}

// GetByName is find attribute by name
func (u Attributes) GetByName(name string) Attributes {
	attr := u[name]
	if attr == nil ||
		(attr.User == nil && attr.Group == nil && attr.ID == 0) {
		return nil
	}
	return Attributes{
		name: attr,
	}
}

// GetByID is find attribute by id
func (u Attributes) GetByID(_id string) Attributes {
	id, _ := strconv.Atoi(_id)
	for k, u := range u {
		if u.ID == id {
			return Attributes{
				k: u,
			}
		}
	}
	return nil
}

func (u Attributes) appendWithSortByID(id int) []int {
	r := []int{id}
	for _, v := range u {
		r = append(r, v.ID)
	}
	nodup := removeDupInts(r)
	sort.Sort(sort.IntSlice(nodup))
	return nodup
}

func (u Attributes) PrevID(id int) int {
	list := u.appendWithSortByID(id)
	for i, v := range list {
		if v == id {
			if i > 0 {
				return list[i-1]
			} else {
				return id
			}
		}
	}
	return 0
}

func (u Attributes) NextID(id int) int {
	list := u.appendWithSortByID(id)
	for i, v := range list {
		if v == id {
			if i < len(list)-1 {
				return list[i+1]
			} else if i == len(list)-1 {
				return list[i]
			} else {
				return -1
			}
		}
	}
	return 0
}
func iMember(n int, xs []int) bool {
	for _, x := range xs {
		if n == x {
			return true
		}
	}
	return false
}

func removeDupInts(xs []int) []int {
	ys := make([]int, 0, len(xs))
	for _, x := range xs {
		if !iMember(x, ys) {
			ys = append(ys, x)
		}
	}
	return ys
}
