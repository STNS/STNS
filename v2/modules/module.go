package main

import (
	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/stns"
)

func syncConfig(b model.Backend, config *stns.Config) error {
	users, err := b.Users()
	if err != nil {
		switch err.(type) {
		case model.NotFoundError:
			users = map[string]model.UserGroup{}
		default:
			return err
		}
	}

	if err := model.SyncConfig("users", b, config.Users.ToUserGroup(), users); err != nil {
		return err
	}

	groups, err := b.Groups()
	if err != nil {
		switch err.(type) {
		case model.NotFoundError:
			groups = map[string]model.UserGroup{}
		default:
			return err
		}
	}

	if err := model.SyncConfig("groups", b, config.Groups.ToUserGroup(), groups); err != nil {
		return err
	}
	return nil
}
