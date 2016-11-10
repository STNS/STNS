package stns

// Group group object
type Group struct {
	Users      []string `toml:"users" json:"users"`
	LinkGroups []string `toml:"link_groups" json:"link_groups"`
}

// LinkParams return link group name
func (g *Group) LinkParams() []string {
	return g.LinkGroups
}

// LinkValue return group users
func (g *Group) LinkValue() []string {
	return g.Users
}

// SetLinkValue set group users
func (g *Group) SetLinkValue(v []string) {
	g.Users = v
}
