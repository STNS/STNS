package model

type BackendDummy struct {
}

func NewBackendDummy() (Backend, error) {
	return BackendDummy{}, nil
}

func (b BackendDummy) FindUserByID(id int) (map[string]UserGroup, error) {
	if id == 1 {
		return map[string]UserGroup{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "test",
					},
					Password: "$6$/C5VdIWEaQVD4Y9D$CQz5Qc99yKucuwvVWIrc2cgnLCOgTbq/QXvKGCXa3f3gYx3xc0/EOhyHAUehS92J9iy8IUqhpnGXpaKYVMoZK1",
				},
			},
			nil
	}

	return nil, NewNotFoundError("user", "dummy")
}

func (b BackendDummy) FindUserByName(name string) (map[string]UserGroup, error) {
	if name == "test" {
		return map[string]UserGroup{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "test",
					},
					Password: "$6$/C5VdIWEaQVD4Y9D$CQz5Qc99yKucuwvVWIrc2cgnLCOgTbq/QXvKGCXa3f3gYx3xc0/EOhyHAUehS92J9iy8IUqhpnGXpaKYVMoZK1",
				},
			},
			nil
	}

	return nil, NewNotFoundError("user", "dummy")
}

func (b BackendDummy) Users() (map[string]UserGroup, error) {
	return nil, nil
}

func (b BackendDummy) FindGroupByID(id int) (map[string]UserGroup, error) {
	return nil, nil
}

func (b BackendDummy) FindGroupByName(name string) (map[string]UserGroup, error) {
	return nil, nil
}

func (b BackendDummy) Groups() (map[string]UserGroup, error) {
	return nil, nil
}

func (b BackendDummy) highlowUserID(high bool) int {
	return 0
}

func (b BackendDummy) HighestUserID() int {
	return 0
}

func (b BackendDummy) LowestUserID() int {
	return 0
}

func (b BackendDummy) highlowGroupID(high bool) int {
	return 0
}

func (b BackendDummy) HighestGroupID() int {
	return 0
}

func (b BackendDummy) LowestGroupID() int {
	return 0
}

func (b BackendDummy) CreateUser(v UserGroup) error {
	return nil
}

func (b BackendDummy) CreateGroup(v UserGroup) error {
	return nil
}

func (b BackendDummy) create(path string, v UserGroup) error {
	return nil
}

func (b BackendDummy) DeleteUser(id int) error {
	return nil
}

func (b BackendDummy) DeleteGroup(id int) error {
	return nil
}

func (b BackendDummy) delete(path string) error {
	return nil
}

func (b BackendDummy) UpdateUser(v UserGroup) error {
	return nil
}

func (b BackendDummy) UpdateGroup(v UserGroup) error {
	return nil
}

func (b BackendDummy) update(path string, v UserGroup) error {
	return nil
}
