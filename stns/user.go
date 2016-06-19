package stns

type User struct {
	Password  string   `toml:"password" json:"password"`
	HashType  string   `toml:"hash_type" json:"hash_type"`
	GroupId   int      `toml:"group_id" json:"group_id"`
	Directory string   `toml:"directory" json:"directory"`
	Shell     string   `toml:"shell" json:"shell"`
	Gecos     string   `toml:"gecos" json:"gecos"`
	Keys      []string `toml:"keys" json:"keys"`
	LinkUsers []string `toml:"link_users" json:"link_users"`
}

func (u *User) LinkParams() []string {
	return u.LinkUsers
}

func (u *User) LinkValue() []string {
	return u.Keys
}

func (u *User) SetLinkValue(v []string) {
	u.Keys = v
}
