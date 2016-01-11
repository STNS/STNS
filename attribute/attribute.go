package attribute

type User struct {
	GroupId   int      `toml:"group_id" json:"group_id"`
	Directory string   `toml:"directory" json:"directory"`
	Shell     string   `toml:"shell" json:"shell"`
	Gecos     string   `toml:"gecos" json:"gecos"`
	Keys      []string `toml:"keys" json:"keys"`
}

type All struct {
	Id int `toml:"id" json:"id"`
	*User
	// use group
	Users []string `toml:"users" json:"users"`
}
