package stns

type Group struct {
	Users      []string `toml:"users" json:"users"`
	LinkGroups []string `toml:"link_groups" json:"link_groups"`
}

func (g *Group) LinkParams() []string {
	return g.LinkGroups
}

func (g *Group) LinkValue() []string {
	return g.Users
}

func (g *Group) SetLinkValue(v []string) {
	g.Users = v
}
