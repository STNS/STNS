package stns

// User user object
type User struct {
	Password  string   `toml:"password" json:"password"`
	GroupID   int      `toml:"group_id" json:"group_id"`
	Directory string   `toml:"directory" json:"directory"`
	Shell     string   `toml:"shell" json:"shell"`
	Gecos     string   `toml:"gecos" json:"gecos"`
	Keys      []string `toml:"keys" json:"keys"`
	LinkUsers []string `toml:"link_users" json:"link_users"`
}

// LinkParams return link users name
func (u *User) LinkParams() []string {
	return u.LinkUsers
}

// LinkValue return ssh keys
func (u *User) LinkValue() []string {
	return u.Keys
}

// SetLinkValue set ssh keys
func (u *User) SetLinkValue(v []string) {
	u.Keys = v
}
