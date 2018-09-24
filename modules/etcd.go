package modules

import (
	"crypto/aes"
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

func NewBackendEtcd(c *stns.Config) (model.GetterBackend, error) {
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

	block, err := aes.NewCipher(c.Etcd.SecretKey)
	if err != nil {
		return nil, err
	}

	return BackendEtcd{
		api:    etcd.NewKeysAPI(cli),
		config: c,
		block:  block,
	}, nil
}

func (b BackendEtcd) decrypt(v string) []byte {
	decryptedText := make([]byte, len(v))
	b.block.Decrypt(decryptedText, []byte(v))
	return decryptedText
}

func (b BackendEtcd) encrypt(v string) string {
	cipherText := make([]byte, len(v))
	b.block.Encrypt(cipherText, []byte(v))
	return string(cipherText)
}

func (b BackendEtcd) FindUserByID(id int) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/users/id/%d", id), nil)
	if err != nil {
		return nil, err
	}

	users := new(model.Users)
	if err := json.Unmarshal(b.decrypt(r.Node.Value), users); err != nil {
		return nil, err
	}
	return users.ToUserGroup(), nil
}

func (b BackendEtcd) FindUserByName(name string) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/users/name/%s", name), nil)
	if err != nil {
		return nil, err
	}

	users := new(model.Users)
	if err := json.Unmarshal(b.decrypt(r.Node.Value), users); err != nil {
		return nil, err
	}
	return users.ToUserGroup(), nil
}

func (b BackendEtcd) Users() (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), "/users", nil)
	if err != nil {
		return nil, err
	}

	users := new(model.Users)
	if err := json.Unmarshal(b.decrypt(r.Node.Value), users); err != nil {
		return nil, err
	}
	return users.ToUserGroup(), nil
}

func (b BackendEtcd) FindGroupByID(id int) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/groups/id/%d", id), nil)
	if err != nil {
		return nil, err
	}

	groups := new(model.Groups)
	if err := json.Unmarshal(b.decrypt(r.Node.Value), groups); err != nil {
		return nil, err
	}
	return groups.ToUserGroup(), nil
}

func (b BackendEtcd) FindGroupByName(name string) (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), fmt.Sprintf("/groups/name/%s", name), nil)
	if err != nil {
		return nil, err
	}

	groups := new(model.Groups)
	if err := json.Unmarshal(b.decrypt(r.Node.Value), groups); err != nil {
		return nil, err
	}
	return groups.ToUserGroup(), nil
}

func (b BackendEtcd) Groups() (map[string]model.UserGroup, error) {
	r, err := b.api.Get(context.Background(), "/groups", nil)
	if err != nil {
		return nil, err
	}

	groups := new(model.Groups)
	if err := json.Unmarshal(b.decrypt(r.Node.Value), groups); err != nil {
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
