package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/stns"
	etcd "go.etcd.io/etcd/client/v2"

	"golang.org/x/net/context"
)

var ModuleName = "Etcd"

type BackendEtcd struct {
	config *stns.Config
	api    etcd.KeysAPI
}

func NewBackendEtcd(c *stns.Config) (model.Backend, error) {
	var endpoints []string
	var user, password string
	if c.Modules["etcd"].(map[string]interface{})["endpoints"] != nil {
		ep := c.Modules["etcd"].(map[string]interface{})["endpoints"].([]interface{})
		for _, e := range ep {
			endpoints = append(endpoints, e.(string))
		}
	}

	if c.Modules["etcd"].(map[string]interface{})["user"] != nil {
		user = c.Modules["etcd"].(map[string]interface{})["user"].(string)
	}

	if c.Modules["etcd"].(map[string]interface{})["password"] != nil {
		password = c.Modules["etcd"].(map[string]interface{})["password"].(string)
	}

	cfg := etcd.Config{
		Endpoints:               endpoints,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
		Username:                user,
		Password:                password,
	}
	cli, err := etcd.New(cfg)
	if err != nil {
		return nil, err
	}

	b := BackendEtcd{
		api:    etcd.NewKeysAPI(cli),
		config: c,
	}
	if c.Modules["etcd"].(map[string]interface{})["sync"] != nil && c.Modules["etcd"].(map[string]interface{})["sync"].(bool) {
		err := syncConfig(b, c)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (b BackendEtcd) FindUserByName(name string) (map[string]model.UserGroup, error) {
	users, err := b.Users()
	r := map[string]model.UserGroup{}

	if err != nil {
		return nil, err
	}

	if users != nil {
		for _, u := range users {
			if u.GetName() == name {
				r[u.GetName()] = u
				return r, nil
			}
		}
	}
	return nil, model.NewNotFoundError("user", nil)
}

func (b BackendEtcd) FindUserByID(id int) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/users/id/%d", id), nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("user", id)
		}
		return nil, err
	}

	users := model.Users{}
	user := new(model.User)
	if err := json.Unmarshal([]byte(r.Node.Value), user); err != nil {
		return nil, err
	}
	users[user.GetName()] = user
	return users.ToUserGroup(), nil
}

func (b BackendEtcd) Users() (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), "/users/id", nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("user", nil)
		}
		return nil, err
	}

	users := model.Users{}
	for _, n := range r.Node.Nodes {
		if n.Value == "null" {
			continue
		}
		user := new(model.User)
		if err := json.Unmarshal([]byte(n.Value), user); err != nil {
			return nil, err
		}
		users[user.GetName()] = user
	}
	return users.ToUserGroup(), nil
}

func (b BackendEtcd) FindGroupByName(name string) (map[string]model.UserGroup, error) {
	groups, err := b.Groups()
	r := map[string]model.UserGroup{}

	if err != nil {
		return nil, err
	}

	if groups != nil {
		for _, g := range groups {
			if g.GetName() == name {
				r[g.GetName()] = g
				return r, nil
			}
		}
	}
	return nil, model.NewNotFoundError("group", nil)
}

func (b BackendEtcd) FindGroupByID(id int) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/groups/id/%d", id), nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("group", id)
		}
		return nil, err
	}

	groups := model.Groups{}
	group := new(model.Group)
	if err := json.Unmarshal([]byte(r.Node.Value), group); err != nil {
		return nil, err
	}
	groups[group.GetName()] = group

	return groups.ToUserGroup(), nil
}

func (b BackendEtcd) Groups() (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), "/groups/id", nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("user", nil)
		}
		return nil, err
	}

	groups := model.Groups{}
	for _, n := range r.Node.Nodes {
		if n.Value == "null" {
			continue
		}
		group := new(model.Group)
		if err := json.Unmarshal([]byte(n.Value), group); err != nil {
			return nil, err
		}
		groups[group.GetName()] = group
	}
	return groups.ToUserGroup(), nil

}

func (b BackendEtcd) highlowUserID(high bool) int {
	ret := 0
	r, err := b.api.Get(context.Background(), "/users/id", nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return ret
		}
		return ret
	}

	for _, n := range r.Node.Nodes {
		user := new(model.User)
		if err := json.Unmarshal([]byte(n.Value), user); err != nil {
			return ret
		}
		if ret == 0 || (high && user.GetID() > ret) || (!high && user.GetID() < ret) {
			ret = user.GetID()
		}
	}
	return ret
}

func (b BackendEtcd) HighestUserID() int {
	return b.highlowUserID(true)
}

func (b BackendEtcd) LowestUserID() int {
	return b.highlowUserID(false)
}

func (b BackendEtcd) highlowGroupID(high bool) int {
	ret := 0
	r, err := b.api.Get(context.Background(), "/groups/id", nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return ret
		}
		return ret
	}

	for _, n := range r.Node.Nodes {
		group := new(model.Group)
		if err := json.Unmarshal([]byte(n.Value), group); err != nil {
			return ret
		}
		if ret == 0 || (high && group.GetID() > ret) || (!high && group.GetID() < ret) {
			ret = group.GetID()
		}
	}
	return ret
}

func (b BackendEtcd) HighestGroupID() int {
	return b.highlowGroupID(true)
}

func (b BackendEtcd) LowestGroupID() int {
	return b.highlowGroupID(false)
}

func (b BackendEtcd) CreateUser(v model.UserGroup) error {
	return b.create(fmt.Sprintf("/users/id/%d", v.GetID()), v)
}

func (b BackendEtcd) CreateGroup(v model.UserGroup) error {
	return b.create(fmt.Sprintf("/groups/id/%d", v.GetID()), v)
}

func (b BackendEtcd) create(path string, v model.UserGroup) error {
	bjson, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := b.api.Set(context.Background(), path, string(bjson), nil); err != nil {
		return err
	}
	return nil
}

func (b BackendEtcd) DeleteUser(id int) error {
	return b.delete(fmt.Sprintf("/users/id/%d", id))
}

func (b BackendEtcd) DeleteGroup(id int) error {
	return b.delete(fmt.Sprintf("/groups/id/%d", id))
}

func (b BackendEtcd) delete(path string) error {
	if _, err := b.api.Delete(context.Background(), path, nil); err != nil {
		return err
	}
	return nil
}

func (b BackendEtcd) UpdateUser(v model.UserGroup) error {
	return b.update(fmt.Sprintf("/users/id/%d", v.GetID()), v)
}

func (b BackendEtcd) UpdateGroup(v model.UserGroup) error {
	return b.update(fmt.Sprintf("/groups/id/%d", v.GetID()), v)
}

func (b BackendEtcd) update(path string, v model.UserGroup) error {
	bjson, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := b.api.Update(context.Background(), path, string(bjson)); err != nil {
		return err
	}
	return nil
}
