package model

type Group struct {
	Base
	Users      []string `toml:"users" json:"users"`
	LinkGroups []string `toml:"link_groups" json:"link_groups"`
}
type Groups map[string]*Group

func (gs *Groups) ToInterfaces() map[string]interface{} {
	if gs != nil {
		igroups := make(map[string]interface{}, len(*gs))
		for k, v := range *gs {
			igroups[k] = v
		}
		return igroups
	}
	return nil
}
