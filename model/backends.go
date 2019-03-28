package model

import "golang.org/x/sync/errgroup"

func (gb Backends) FindUserByID(v int) (map[string]UserGroup, error) {
	r := map[string]UserGroup{}
	var notfound error
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			lr, err := b.FindUserByID(v)
			if err != nil {
				switch err.(type) {
				case NotFoundError:
					notfound = err
				default:
					return err
				}
			}
			r = mergeUserGroup(r, lr)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// record notfound
	if len(r) == 0 {
		return nil, notfound
	}

	return r, nil
}
func (gb Backends) FindUserByName(v string) (map[string]UserGroup, error) {
	r := map[string]UserGroup{}
	var notfound error
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			lr, err := b.FindUserByName(v)
			if err != nil {
				switch err.(type) {
				case NotFoundError:
					notfound = err
				default:
					return err
				}
			}
			r = mergeUserGroup(r, lr)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// record notfound
	if len(r) == 0 {
		return nil, notfound
	}

	return r, nil
}
func (gb Backends) FindGroupByID(v int) (map[string]UserGroup, error) {
	r := map[string]UserGroup{}
	var notfound error
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			lr, err := b.FindGroupByID(v)
			if err != nil {
				switch err.(type) {
				case NotFoundError:
					notfound = err
				default:
					return err
				}
			}
			r = mergeUserGroup(r, lr)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// record notfound
	if len(r) == 0 {
		return nil, notfound
	}

	return r, nil
}
func (gb Backends) FindGroupByName(v string) (map[string]UserGroup, error) {
	r := map[string]UserGroup{}
	var notfound error
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			lr, err := b.FindGroupByName(v)
			if err != nil {
				switch err.(type) {
				case NotFoundError:
					notfound = err
				default:
					return err
				}
			}
			r = mergeUserGroup(r, lr)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// record notfound
	if len(r) == 0 {
		return nil, notfound
	}

	return r, nil
}
func (gb Backends) Users() (map[string]UserGroup, error) {
	r := map[string]UserGroup{}
	var notfound error
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			lr, err := b.Users()
			if err != nil {
				switch err.(type) {
				case NotFoundError:
					notfound = err
				default:
					return err
				}
			}
			r = mergeUserGroup(r, lr)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// record notfound
	if len(r) == 0 {
		return nil, notfound
	}

	return r, nil
}
func (gb Backends) Groups() (map[string]UserGroup, error) {
	r := map[string]UserGroup{}
	var notfound error
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			lr, err := b.Groups()
			if err != nil {
				switch err.(type) {
				case NotFoundError:
					notfound = err
				default:
					return err
				}
			}
			r = mergeUserGroup(r, lr)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// record notfound
	if len(r) == 0 {
		return nil, notfound
	}

	return r, nil
}
func (gb Backends) CreateUser(v UserGroup) error {
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			err := b.CreateUser(v)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil
	}

	return nil
}
func (gb Backends) CreateGroup(v UserGroup) error {
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			err := b.CreateGroup(v)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil
	}

	return nil
}
func (gb Backends) DeleteUser(v int) error {
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			err := b.DeleteUser(v)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil
	}

	return nil
}
func (gb Backends) DeleteGroup(v int) error {
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			err := b.DeleteGroup(v)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil
	}

	return nil
}
func (gb Backends) UpdateUser(v UserGroup) error {
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			err := b.UpdateUser(v)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil
	}

	return nil
}
func (gb Backends) UpdateGroup(v UserGroup) error {
	eg := errgroup.Group{}
	for _, b := range gb {
		eg.Go(func() error {
			err := b.UpdateGroup(v)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil
	}

	return nil
}
func (gb Backends) HighestUserID() int {
	r := 0
	for _, b := range gb {
		lr := b.HighestUserID()
		if lr != 0 {
			r = lr
		}
	}
	return r
}
func (gb Backends) LowestUserID() int {
	r := 0
	for _, b := range gb {
		lr := b.LowestUserID()
		if lr != 0 {
			r = lr
		}
	}
	return r
}
func (gb Backends) HighestGroupID() int {
	r := 0
	for _, b := range gb {
		lr := b.HighestGroupID()
		if lr != 0 {
			r = lr
		}
	}
	return r
}
func (gb Backends) LowestGroupID() int {
	r := 0
	for _, b := range gb {
		lr := b.LowestGroupID()
		if lr != 0 {
			r = lr
		}
	}
	return r
}
