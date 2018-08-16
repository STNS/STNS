package model

type User struct {
	Base
	Password      string   `toml:"password" json:"password"`
	GroupID       int      `toml:"group_id" json:"group_id"`
	Directory     string   `toml:"directory" json:"directory"`
	Shell         string   `toml:"shell" json:"shell"`
	Gecos         string   `toml:"gecos" json:"gecos"`
	Keys          []string `toml:"keys" json:"keys"`
	LinkUsers     []string `toml:"link_users" json:"link_users"`
	SetupCommands []string `toml:"setup_commands" json:"setup_commands"`
}

type Users map[string]*User

func (us *Users) ToInterfaces() map[string]interface{} {
	if us != nil {
		iusers := make(map[string]interface{}, len(*us))
		for k, v := range *us {
			iusers[k] = v
		}
		return iusers
	}
	return nil
}
