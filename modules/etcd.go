package modules

import (
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"time"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
	etcd "go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

type BackendEtcd struct {
	config *stns.Config
	api    etcd.KeysAPI
	block  cipher.Block
}

func NewBackendEtcd(c *stns.Config) (model.Backend, error) {
	cfg := etcd.Config{
		Endpoints:               c.Etcd.Endpoints,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
		Username:                c.Etcd.User,
		Password:                c.Etcd.Password,
	}
	cli, err := etcd.New(cfg)
	if err != nil {
		return nil, err
	}

	return BackendEtcd{
		api:    etcd.NewKeysAPI(cli),
		config: c,
	}, nil
}

func (b BackendEtcd) FindUserByID(id int) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/users/id/%d", id), nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("user", id)
		}
		return nil, err
	}

	users := new(model.Users)
	if err := json.Unmarshal([]byte(r.Node.Value), users); err != nil {
		return nil, err
	}
	return users.ToUserGroup(), nil
}

func (b BackendEtcd) FindUserByName(name string) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/users/name/%s", name), nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("user", name)
		}
	}

	users := new(model.Users)
	if err := json.Unmarshal([]byte(r.Node.Value), users); err != nil {
		return nil, err
	}
	return users.ToUserGroup(), nil
}

func (b BackendEtcd) Users() (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), "/users", nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("user", nil)
		}
		return nil, err
	}

	users := new(model.Users)
	if err := json.Unmarshal([]byte(r.Node.Value), users); err != nil {
		return nil, err
	}
	return users.ToUserGroup(), nil
}

func (b BackendEtcd) FindGroupByID(id int) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/groups/id/%d", id), nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("group", id)
		}
		return nil, err
	}

	groups := new(model.Groups)
	if err := json.Unmarshal([]byte(r.Node.Value), groups); err != nil {
		return nil, err
	}
	return groups.ToUserGroup(), nil
}

func (b BackendEtcd) FindGroupByName(name string) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/groups/name/%s", name), nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("group", name)
		}
		return nil, err
	}

	groups := new(model.Groups)
	if err := json.Unmarshal([]byte(r.Node.Value), groups); err != nil {
		return nil, err
	}
	return groups.ToUserGroup(), nil
}

func (b BackendEtcd) Groups() (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), "/groups", nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return nil, model.NewNotFoundError("group", nil)
		}
		return nil, err
	}

	groups := new(model.Groups)
	if err := json.Unmarshal([]byte(r.Node.Value), groups); err != nil {
		return nil, err
	}
	return groups.ToUserGroup(), nil
}

func (b BackendEtcd) HighestUserID() int {
	return 0
}

func (b BackendEtcd) LowestUserID() int {
	return 0
}

func (b BackendEtcd) HighestGroupID() int {
	return 0
}

func (b BackendEtcd) LowestGroupID() int {
	return 0
}

func (b BackendEtcd) Create(path string, v map[string]model.UserGroup) error {
	bjson, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err := b.api.Set(context.Background(), path, string(bjson), nil); err != nil {
		return err
	}
	return nil
}
