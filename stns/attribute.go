package stns

import "strconv"

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
