package model

type UserGroup interface {
	GetID() int
	GetName() string
	setName(string)
	linkValues() []string
	setLinkValues([]string)
	value() []string
}
type Backend interface {
	// Getter
	FindUserByID(int) (map[string]UserGroup, error)
	FindUserByName(string) (map[string]UserGroup, error)
	FindGroupByID(int) (map[string]UserGroup, error)
	FindGroupByName(string) (map[string]UserGroup, error)
	Users() (map[string]UserGroup, error)
	Groups() (map[string]UserGroup, error)
	HighestUserID() int
	LowestUserID() int
	HighestGroupID() int
	LowestGroupID() int

	// Setter
	CreateUser(UserGroup) error
	DeleteUser(int) error
	UpdateUser(UserGroup) error
	CreateGroup(UserGroup) error
	DeleteGroup(int) error
	UpdateGroup(UserGroup) error
}

func SyncConfig(resourceName string, b Backend, configResources, backendResources map[string]UserGroup) error {
	var backendResource UserGroup
	if configResources != nil {
		for _, cu := range configResources {
			found := false

			if backendResources != nil {
				for _, eu := range backendResources {
					if cu.GetID() == eu.GetID() {
						backendResource = eu
						found = true
						break
					}
				}
			}

			if found {
				if resourceName == "users" {
					// not overwrite password
					cu.(*User).Password = backendResource.(*User).Password

					if err := b.UpdateUser(cu); err != nil {
						return err
					}
				} else {
					if err := b.UpdateGroup(cu); err != nil {
						return err
					}
				}
			} else {
				if resourceName == "users" {
					if err := b.CreateUser(cu); err != nil {
						return err
					}

				} else {
					if err := b.CreateGroup(cu); err != nil {
						return err
					}
				}
			}
		}
	}

	if backendResources != nil {
		for _, eu := range backendResources {
			found := false
			if configResources != nil {
				for _, cu := range configResources {
					if cu.GetID() == eu.GetID() {
						found = true
						break
					}
				}
			}
			if !found {
				if resourceName == "users" {
					if err := b.DeleteUser(eu.GetID()); err != nil {
						return err
					}

				} else {
					if err := b.DeleteGroup(eu.GetID()); err != nil {
						return err
					}
				}
			}
		}

	}
	return nil
}
