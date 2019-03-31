package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"

	"gopkg.in/redis.v5"
)

type BackendRedis struct {
	*redis.Client
	backend Backend
	ttl     int
	logger  *log.Logger
}

const defaultRedisTTL = 600

const (
	userIDKey         = "user_id_%d"
	userNameKey       = "user_name_%s"
	usersKey          = "users"
	groupIDKey        = "group_id_%d"
	groupNameKey      = "group_name_%s"
	groupsKey         = "groups"
	userLowestIDKey   = "lowest_user_id"
	userHighestIDKey  = "highest_user_id"
	groupLowestIDKey  = "lowest_group_id"
	groupHighestIDKey = "highest_group_id"
)

var timeToNotCache *time.Time

func NewBackendRedis(b Backend, logger *log.Logger, host string, port int, password string, ttl int, db int) (*BackendRedis, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})
	_, err := redis.Ping().Result()
	if err != nil {
		return nil, err
	}

	if ttl == 0 {
		ttl = defaultRedisTTL
	}
	return &BackendRedis{
		Client:  redis,
		backend: b,
		ttl:     ttl,
		logger:  logger,
	}, nil
}

func (b *BackendRedis) checkSession() error {
	if timeToNotCache == nil || timeToNotCache.Unix() < time.Now().Unix() {
		_, err := b.Ping().Result()
		if err != nil {
			tn := time.Now().Add(time.Duration(10) * time.Minute)
			timeToNotCache = &tn
			b.logger.Warn(fmt.Sprintf("redis has been disconnected. reconnect time is %s", timeToNotCache))
		} else {
			timeToNotCache = nil
		}
	}

	if timeToNotCache != nil {
		return fmt.Errorf("will not cache until %s", timeToNotCache)
	}
	return nil
}

func (b *BackendRedis) GetCache(key string) string {
	if err := b.checkSession(); err != nil {
		return ""
	}
	v, e := b.Get(key).Result()

	if e == nil {
		return v
	}
	return ""
}

func (b *BackendRedis) SetCache(key string, src interface{}) {
	var mjson []byte
	if err := b.checkSession(); err != nil {
		return
	}

	j, err := json.Marshal(src)
	if err == nil {
		mjson = j
	}

	if len(mjson) == 0 {
		b.logger.Warn(fmt.Sprintf("json unmarshal error. [value: %T]\n", src))
		return
	}

	err = b.Set(key, mjson, time.Duration(b.ttl)*time.Second).Err()
	if err != nil {
		b.logger.Warn(fmt.Printf("Cannot set redis value: %s err: %s", key, err))
	}
}

func (b *BackendRedis) DelCache(key string) {
	if err := b.checkSession(); err != nil {
		return
	}

	err := b.Del(key).Err()
	if err != nil {
		b.logger.Warn(fmt.Printf("Cannot del redis value: %s err: %s", key, err))
	}
}

func (b *BackendRedis) FindUserByID(id int) (map[string]UserGroup, error) {
	v := b.GetCache(fmt.Sprintf(userIDKey, id))
	if v != "" {
		d := map[string]UserGroup{}
		if json.Unmarshal([]byte(v), d) == nil {
			return d, nil
		}
	}
	u, err := b.backend.FindUserByID(id)
	if err != nil {
		return nil, err
	}
	b.SetCache(fmt.Sprintf(userIDKey, id), u)
	return u, nil
}

func (b *BackendRedis) FindUserByName(name string) (map[string]UserGroup, error) {
	v := b.GetCache(fmt.Sprintf(userNameKey, name))
	if v != "" {
		d := map[string]UserGroup{}
		if json.Unmarshal([]byte(v), d) == nil {
			return d, nil
		}
	}
	u, err := b.backend.FindUserByName(name)
	if err != nil {
		return nil, err
	}
	b.SetCache(fmt.Sprintf(userNameKey, name), u)
	return u, nil
}

func (b *BackendRedis) Users() (map[string]UserGroup, error) {
	v := b.GetCache(usersKey)
	if v != "" {
		d := map[string]UserGroup{}
		if json.Unmarshal([]byte(v), d) == nil {
			return d, nil
		}
	}
	u, err := b.backend.Users()
	if err != nil {
		return nil, err
	}
	b.SetCache(usersKey, u)
	return u, nil
}

func (b *BackendRedis) FindGroupByID(id int) (map[string]UserGroup, error) {
	v := b.GetCache(fmt.Sprintf(groupIDKey, id))
	if v != "" {
		d := map[string]UserGroup{}
		if json.Unmarshal([]byte(v), d) == nil {
			return d, nil
		}
	}
	u, err := b.backend.FindGroupByID(id)
	if err != nil {
		return nil, err
	}
	b.SetCache(fmt.Sprintf(groupIDKey, id), u)
	return u, nil
}

func (b *BackendRedis) FindGroupByName(name string) (map[string]UserGroup, error) {
	v := b.GetCache(fmt.Sprintf(groupNameKey, name))
	if v != "" {
		d := map[string]UserGroup{}
		if json.Unmarshal([]byte(v), d) == nil {
			return d, nil
		}
	}
	u, err := b.backend.FindGroupByName(name)
	if err != nil {
		return nil, err
	}
	b.SetCache(fmt.Sprintf(groupNameKey, name), u)
	return u, nil
}

func (b *BackendRedis) Groups() (map[string]UserGroup, error) {
	v := b.GetCache(groupsKey)
	if v != "" {
		if v != "" {
			d := map[string]UserGroup{}
			if json.Unmarshal([]byte(v), d) == nil {
				return d, nil
			}
		}
	}
	u, err := b.backend.Groups()
	if err != nil {
		return nil, err
	}
	b.SetCache(groupsKey, u)
	return u, nil
}

func (b *BackendRedis) HighestUserID() int {
	v := b.GetCache(userHighestIDKey)
	if v != "" {
		if v != "" {
			var d int
			if json.Unmarshal([]byte(v), d) == nil {
				return d
			}
		}
	}
	i := b.backend.HighestUserID()
	b.SetCache(userHighestIDKey, strconv.Itoa(i))
	return i
}

func (b *BackendRedis) LowestUserID() int {
	v := b.GetCache(userLowestIDKey)
	if v != "" {
		d := 0
		if json.Unmarshal([]byte(v), d) == nil {
			return d
		}
	}
	i := b.backend.LowestUserID()
	b.SetCache(userLowestIDKey, strconv.Itoa(i))
	return i
}

func (b *BackendRedis) HighestGroupID() int {
	v := b.GetCache(groupHighestIDKey)
	if v != "" {
		d := 0
		if json.Unmarshal([]byte(v), d) == nil {
			return d
		}
	}
	i := b.backend.HighestGroupID()
	b.SetCache(groupHighestIDKey, strconv.Itoa(i))
	return i
}

func (b *BackendRedis) LowestGroupID() int {
	v := b.GetCache(groupLowestIDKey)
	if v != "" {
		d := 0
		if json.Unmarshal([]byte(v), d) == nil {
			return d
		}
	}
	i := b.backend.LowestGroupID()
	b.SetCache(groupLowestIDKey, strconv.Itoa(i))
	return i
}
func (b *BackendRedis) CreateUser(u UserGroup) error {
	if err := b.CreateUser(u); err != nil {
		return err
	}
	b.DelCache(usersKey)
	return nil
}
func (b *BackendRedis) DeleteUser(id int) error {
	us, err := b.backend.FindUserByID(id)
	if err != nil {
		return err
	}

	user := new(User)
	for _, u := range us {
		user = u.(*User)
		break
	}
	if err := b.backend.DeleteUser(id); err != nil {
		return err
	}
	b.DelCache(usersKey)
	b.DelCache(fmt.Sprintf(userIDKey, id))
	b.DelCache(fmt.Sprintf(userNameKey, user.Name))
	return nil
}
func (b *BackendRedis) UpdateUser(us UserGroup) error {
	cu, err := b.backend.FindUserByID(us.GetID())
	if err != nil {
		return err
	}

	user := new(User)
	for _, u := range cu {
		user = u.(*User)
		break
	}

	if err := b.backend.UpdateUser(us); err != nil {
		return err
	}
	b.DelCache(usersKey)
	b.DelCache(fmt.Sprintf(userIDKey, user.GetID()))
	b.DelCache(fmt.Sprintf(userNameKey, user.GetName()))
	return nil
}

func (b *BackendRedis) CreateGroup(u UserGroup) error {
	if err := b.backend.CreateGroup(u); err != nil {
		return err
	}
	b.DelCache(groupsKey)
	return nil
}
func (b *BackendRedis) DeleteGroup(id int) error {
	cu, err := b.backend.FindGroupByID(id)
	if err != nil {
		return err
	}

	group := new(Group)
	for _, u := range cu {
		group = u.(*Group)
		break
	}

	if err := b.backend.DeleteGroup(id); err != nil {
		return err
	}

	b.DelCache(groupsKey)
	b.DelCache(fmt.Sprintf(groupIDKey, id))
	b.DelCache(fmt.Sprintf(groupNameKey, group.Name))
	return nil
}
func (b *BackendRedis) UpdateGroup(us UserGroup) error {
	cu, err := b.backend.FindGroupByID(us.GetID())
	if err != nil {
		return err
	}
	group := new(Group)
	for _, u := range cu {
		group = u.(*Group)
		break
	}
	if err := b.backend.UpdateGroup(us); err != nil {
		return err
	}
	b.DelCache(groupsKey)
	b.DelCache(fmt.Sprintf(groupIDKey, group.GetID()))
	b.DelCache(fmt.Sprintf(groupNameKey, group.GetName()))
	return nil
}
