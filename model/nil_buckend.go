package model

type BackendNil struct {
}

func NewBackendNil() (Backend, error) {
	return BackendNil{}, nil
}

func (b BackendNil) FindUserByID(id int) (map[string]UserGroup, error) {
	return nil, NewNotFoundError("user", "dummy")
}

func (b BackendNil) FindUserByName(name string) (map[string]UserGroup, error) {
	return nil, NewNotFoundError("user", "dummy")
}

func (b BackendNil) Users() (map[string]UserGroup, error) {
	return nil, NewNotFoundError("user", "dummy")
}

func (b BackendNil) FindGroupByID(id int) (map[string]UserGroup, error) {
	return nil, NewNotFoundError("group", "dummy")
}

func (b BackendNil) FindGroupByName(name string) (map[string]UserGroup, error) {
	return nil, NewNotFoundError("group", "dummy")
}

func (b BackendNil) Groups() (map[string]UserGroup, error) {
	return nil, NewNotFoundError("group", "dummy")
}

func (b BackendNil) highlowUserID(high bool) int {
	return 0
}

func (b BackendNil) HighestUserID() int {
	return 0
}

func (b BackendNil) LowestUserID() int {
	return 0
}

func (b BackendNil) highlowGroupID(high bool) int {
	return 0
}

func (b BackendNil) HighestGroupID() int {
	return 0
}

func (b BackendNil) LowestGroupID() int {
	return 0
}

func (b BackendNil) CreateUser(v UserGroup) error {
	return nil
}

func (b BackendNil) CreateGroup(v UserGroup) error {
	return nil
}

func (b BackendNil) create(path string, v UserGroup) error {
	return nil
}

func (b BackendNil) DeleteUser(id int) error {
	return nil
}

func (b BackendNil) DeleteGroup(id int) error {
	return nil
}

func (b BackendNil) delete(path string) error {
	return nil
}

func (b BackendNil) UpdateUser(v UserGroup) error {
	return nil
}

func (b BackendNil) UpdateGroup(v UserGroup) error {
	return nil
}

func (b BackendNil) update(path string, v UserGroup) error {
	return nil
}
