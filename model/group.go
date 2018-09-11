package model

type Group struct {
	Base       `yaml:",inline"`
	Users      []string `toml:"users" json:"users"`
	LinkGroups []string `toml:"link_groups" yaml:"link_groups" json:"-"`
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

func (u *Group) setLinkValues(ks []string) {
	u.Users = uniqStrings(append(u.Users, ks...))
}

func (u *Group) linkValues() []string {
	return u.LinkGroups
}

func (u *Group) value() []string {
	return u.Users
}

func (u *Group) name() string {
	return u.Base.Name
}
