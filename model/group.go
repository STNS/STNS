package model

type Group struct {
	Base
	Users      []string `toml:"users" json:"users"`
	LinkGroups []string `toml:"link_groups" json:"link_groups"`
}
type Groups map[string]*Group

func (gs *Groups) ToUserGroup() map[string]UserGroup {
	if gs != nil {
		igroups := make(map[string]UserGroup, len(*gs))
		for k, v := range *gs {
			igroups[k] = v
		}
		return igroups
	}
	return nil
}
