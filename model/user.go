package model

type User struct {
	Base
	Password      string   `toml:"password" json:"password"`
	GroupID       int      `toml:"group_id" json:"group_id"`
	Directory     string   `toml:"directory" json:"directory"`
	Shell         string   `toml:"shell" json:"shell"`
	Gecos         string   `toml:"gecos" json:"gecos"`
	Keys          []string `toml:"keys" json:"keys"`
	LinkUsers     []string `toml:"link_users" json:"-"`
	SetupCommands []string `toml:"setup_commands" json:"setup_commands"`
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

func (u *User) name() string {
	return u.Base.Name
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
