package model

type User struct {
	Base      `yaml:",inline"`
	Password  string   `toml:"password" yaml:"password" json:"password"`
	GroupID   int      `toml:"group_id" yaml:"group_id" json:"group_id"`
	Directory string   `toml:"directory" yaml:"directory" json:"directory"`
	Shell     string   `toml:"shell" yaml:"shell" json:"shell"`
	Gecos     string   `toml:"gecos" yaml:"gecos" json:"gecos"`
	Keys      []string `toml:"keys" yaml:"keys" json:"keys"`
	LinkUsers []string `toml:"link_users" yaml:"link_users" json:"-"`
}

type Users map[string]*User

func (u *User) linkValues() []string {
	return u.LinkUsers
}
func (u *User) setLinkValues(ks []string) {
	u.Keys = uniqStrings(append(u.Keys, ks...))
}

func (u *User) value() []string {
	return u.Keys
}

func (us *Users) ToUserGroup() map[string]UserGroup {
	if us != nil {
		iusers := make(map[string]UserGroup, len(*us))
		for k, v := range *us {
			iusers[k] = v
		}
		return iusers
	}
	return nil
}
