package model

type Group struct {
	ID         int      `toml:"id" json:"id"`
	Name       string   `toml:"name" json:"name"`
	Users      []string `toml:"users" json:"users"`
	LinkGroups []string `toml:"link_groups" json:"link_groups"`
}
type Groups map[string]*Group

func (gs *Groups) EnsureName() {
	if gs != nil {
		for k, v := range *gs {
			v.Name = k
		}
	}
}
