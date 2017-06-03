package stns

import (
	"sort"
	"strconv"
)

// Attribute attribute object
type Attribute struct {
	ID     int `toml:"id" json:"id"`
	PrevID int `json:"-"`
	NextID int `json:"-"`
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

func (u Attributes) SortByID() SortAttributes {
	pl := make(SortAttributes, len(u))
	i := 0
	for k, v := range u {
		pl[i] = SortAttributePair{k, v.ID}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl

}

type SortAttributePair struct {
	Name string
	ID   int
}

type SortAttributes []SortAttributePair

func (p SortAttributes) Len() int           { return len(p) }
func (p SortAttributes) Less(i, j int) bool { return p[i].ID > p[j].ID }
func (p SortAttributes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p SortAttributes) PrevID(id int) int {
	for k, v := range p {
		if v.ID == id {
			if k > 0 {
				return p[k-1].ID
			} else {
				return 0
			}
		}
	}
	return 0
}

func (p SortAttributes) NextID(id int) int {
	for k, v := range p {
		if v.ID == id {
			if len(p)-1 > k {
				return p[k+1].ID
			} else {
				return 0
			}
		}
	}
	return 0
}
