package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	redis "gopkg.in/redis.v5"
)

func TestBackendRedis_FindUserByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				id: 1,
			},
			want: map[string]UserGroup{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "test",
					},
					Password: "$6$/C5VdIWEaQVD4Y9D$CQz5Qc99yKucuwvVWIrc2cgnLCOgTbq/QXvKGCXa3f3gYx3xc0/EOhyHAUehS92J9iy8IUqhpnGXpaKYVMoZK1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})

				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     10,
				}

				got, err := b.FindUserByID(tt.args.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("BackendRedis.FindUserByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BackendRedis.FindUserByID() = %v, want %v", got, tt.want)
				}

				d := Users{}
				result, _ := redis.Get(fmt.Sprintf(userIDKey, tt.args.id)).Result()
				if err := json.Unmarshal([]byte(result), &d); err != nil {
					t.Errorf("BackendRedis.FindUserByID() error = %v", err)
				}

				if !reflect.DeepEqual(d["test"], tt.want["test"]) {
					t.Errorf("BackendRedis.FindUserByID() = %v, want %v", d, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_FindUserByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "test",
			},
			want: map[string]UserGroup{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "test",
					},
					Password: "$6$/C5VdIWEaQVD4Y9D$CQz5Qc99yKucuwvVWIrc2cgnLCOgTbq/QXvKGCXa3f3gYx3xc0/EOhyHAUehS92J9iy8IUqhpnGXpaKYVMoZK1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {

				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     10,
				}

				got, err := b.FindUserByName(tt.args.name)
				if (err != nil) != tt.wantErr {
					t.Errorf("BackendRedis.FindUserByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BackendRedis.FindUserByName() = %v, want %v", got, tt.want)
				}

				d := Users{}
				result, _ := redis.Get(fmt.Sprintf(userNameKey, tt.args.name)).Result()
				if err := json.Unmarshal([]byte(result), &d); err != nil {
					t.Errorf("BackendRedis.FindUserByName() error = %v", err)
				}

				if !reflect.DeepEqual(d["test"], tt.want["test"]) {
					t.Errorf("BackendRedis.FindUserByName() = %v, want %v", d, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_Users(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			want: map[string]UserGroup{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "test",
					},
					Password: "$6$/C5VdIWEaQVD4Y9D$CQz5Qc99yKucuwvVWIrc2cgnLCOgTbq/QXvKGCXa3f3gYx3xc0/EOhyHAUehS92J9iy8IUqhpnGXpaKYVMoZK1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     10,
				}

				got, err := b.Users()
				if (err != nil) != tt.wantErr {
					t.Errorf("BackendRedis.Users() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BackendRedis.Users() = %v, want %v", got, tt.want)
				}

				d := Users{}
				result, _ := redis.Get(usersKey).Result()
				if err := json.Unmarshal([]byte(result), &d); err != nil {
					t.Errorf("BackendRedis.Users() error = %v", err)
				}

				if !reflect.DeepEqual(d["test"], tt.want["test"]) {
					t.Errorf("BackendRedis.Users() = %v, want %v", d, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_FindGroupByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				id: 1,
			},
			want: map[string]UserGroup{
				"test": &Group{
					Base: Base{
						ID:   1,
						Name: "test",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     10,
				}

				got, err := b.FindGroupByID(tt.args.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("BackendRedis.FindGroupByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BackendRedis.FindGroupByID() = %v, want %v", got, tt.want)
				}

				d := Groups{}
				result, _ := redis.Get(fmt.Sprintf(groupIDKey, tt.args.id)).Result()
				if err := json.Unmarshal([]byte(result), &d); err != nil {
					t.Errorf("BackendRedis.FindGroupByID() error = %v", err)
				}

				if !reflect.DeepEqual(d["test"], tt.want["test"]) {
					t.Errorf("BackendRedis.FindGroupByID() = %v, want %v", d, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_FindGroupByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "test",
			},
			want: map[string]UserGroup{
				"test": &Group{
					Base: Base{
						ID:   1,
						Name: "test",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     10,
				}

				got, err := b.FindGroupByName(tt.args.name)
				if (err != nil) != tt.wantErr {
					t.Errorf("BackendRedis.FindGroupByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BackendRedis.FindGroupByName() = %v, want %v", got, tt.want)
				}

				d := Groups{}
				result, _ := redis.Get(fmt.Sprintf(groupNameKey, tt.args.name)).Result()
				if err := json.Unmarshal([]byte(result), &d); err != nil {
					t.Errorf("BackendRedis.FindGroupByName() error = %v", err)
				}

				if !reflect.DeepEqual(d["test"], tt.want["test"]) {
					t.Errorf("BackendRedis.FindGroupByName() = %v, want %v", d, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_Groups(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			want: map[string]UserGroup{
				"test": &Group{
					Base: Base{
						ID:   1,
						Name: "test",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     10,
				}

				got, err := b.Groups()
				if (err != nil) != tt.wantErr {
					t.Errorf("BackendRedis.Groups() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BackendRedis.Groups() = %v, want %v", got, tt.want)
				}

				d := Groups{}
				result, _ := redis.Get(groupsKey).Result()
				if err := json.Unmarshal([]byte(result), &d); err != nil {
					t.Errorf("BackendRedis.Groups() error = %v", err)
				}

				if !reflect.DeepEqual(d["test"], tt.want["test"]) {
					t.Errorf("BackendRedis.Groups() = %v, want %v", d, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_HighestUserID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "ok",
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     10,
				}

				if got := b.HighestUserID(); got != tt.want {
					t.Errorf("BackendRedis.HighestUserID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_LowestUserID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "ok",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     1,
				}

				if got := b.LowestUserID(); got != tt.want {
					t.Errorf("BackendRedis.LowestUserID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_HighestGroupID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "ok",
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     10,
				}

				if got := b.HighestGroupID(); got != tt.want {
					t.Errorf("BackendRedis.HighestGroupID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_LowestGroupID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "ok",
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, origin := range []Backend{
				BackendDummy{},
				BackendNil{},
			} {
				redis := redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
				})
				b := &BackendRedis{
					backend: origin,
					Client:  redis,
					ttl:     1,
				}

				if got := b.LowestGroupID(); got != tt.want {
					t.Errorf("BackendRedis.LowestGroupID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestBackendRedis_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		redis := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		})

		b := &BackendRedis{
			backend: BackendDummy{},
			Client:  redis,
			ttl:     10,
		}

		err := b.Set(usersKey, "test", time.Duration(10)*time.Second).Err()
		if err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := b.CreateUser(nil); (err != nil) != tt.wantErr {
				t.Errorf("BackendRedis.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, e := b.Get(usersKey).Result()
			if e == nil {
				t.Error("BackendRedis.CreateUser() can't purge cache")
			}

		})
	}
}

func TestBackendRedis_UpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		args    UserGroup
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
			args: &User{
				Base: Base{
					ID:   1,
					Name: "test",
				},
				Password: "$6$/C5VdIWEaQVD4Y9D$CQz5Qc99yKucuwvVWIrc2cgnLCOgTbq/QXvKGCXa3f3gYx3xc0/EOhyHAUehS92J9iy8IUqhpnGXpaKYVMoZK1",
			},
		},
	}
	for _, tt := range tests {
		redis := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		})

		b := &BackendRedis{
			backend: BackendDummy{},
			Client:  redis,
			ttl:     10,
		}

		keys := []string{
			usersKey,
			fmt.Sprintf(userIDKey, 1),
			fmt.Sprintf(userNameKey, "test"),
		}
		for _, k := range keys {
			err := b.Set(k, "test", time.Duration(10)*time.Second).Err()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := b.UpdateUser(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("BackendRedis.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, k := range keys {
				_, e := b.Get(k).Result()
				if e == nil {
					t.Error("BackendRedis.UpdateUser() can't purge cache")
				}
			}
		})
	}
}

func TestBackendRedis_DeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		args    UserGroup
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
			args: &User{
				Base: Base{
					ID:   1,
					Name: "test",
				},
				Password: "$6$/C5VdIWEaQVD4Y9D$CQz5Qc99yKucuwvVWIrc2cgnLCOgTbq/QXvKGCXa3f3gYx3xc0/EOhyHAUehS92J9iy8IUqhpnGXpaKYVMoZK1",
			},
		},
	}
	for _, tt := range tests {
		redis := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		})

		b := &BackendRedis{
			backend: BackendDummy{},
			Client:  redis,
			ttl:     10,
		}

		keys := []string{
			usersKey,
			fmt.Sprintf(userIDKey, 1),
			fmt.Sprintf(userNameKey, "test"),
		}
		for _, k := range keys {
			err := b.Set(k, "test", time.Duration(10)*time.Second).Err()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := b.DeleteUser(tt.args.GetID()); (err != nil) != tt.wantErr {
				t.Errorf("BackendRedis.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, k := range keys {
				_, e := b.Get(k).Result()
				if e == nil {
					t.Error("BackendRedis.DeleteUser() can't purge cache")
				}
			}
		})
	}
}

func TestBackendRedis_CreateGroup(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		redis := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		})

		b := &BackendRedis{
			backend: BackendDummy{},
			Client:  redis,
			ttl:     10,
		}

		err := b.Set(groupsKey, "test", time.Duration(10)*time.Second).Err()
		if err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := b.CreateGroup(nil); (err != nil) != tt.wantErr {
				t.Errorf("BackendRedis.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, e := b.Get(groupsKey).Result()
			if e == nil {
				t.Error("BackendRedis.CreateGroup() can't purge cache")
			}

		})
	}
}

func TestBackendRedis_UpdateGroup(t *testing.T) {
	tests := []struct {
		name    string
		args    UserGroup
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
			args: &Group{
				Base: Base{
					ID:   1,
					Name: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		redis := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		})

		b := &BackendRedis{
			backend: BackendDummy{},
			Client:  redis,
			ttl:     10,
		}

		keys := []string{
			groupsKey,
			fmt.Sprintf(groupIDKey, 1),
			fmt.Sprintf(groupNameKey, "test"),
		}
		for _, k := range keys {
			err := b.Set(k, "test", time.Duration(10)*time.Second).Err()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := b.UpdateGroup(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("BackendRedis.UpdateGroup() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, k := range keys {
				_, e := b.Get(k).Result()
				if e == nil {
					t.Error("BackendRedis.UpdateGroup() can't purge cache")
				}
			}
		})
	}
}

func TestBackendRedis_DeleteGroup(t *testing.T) {
	tests := []struct {
		name    string
		args    UserGroup
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
			args: &Group{
				Base: Base{
					ID:   1,
					Name: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		redis := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		})

		b := &BackendRedis{
			backend: BackendDummy{},
			Client:  redis,
			ttl:     10,
		}

		keys := []string{
			groupsKey,
			fmt.Sprintf(groupIDKey, 1),
			fmt.Sprintf(groupNameKey, "test"),
		}
		for _, k := range keys {
			err := b.Set(k, "test", time.Duration(10)*time.Second).Err()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := b.DeleteGroup(tt.args.GetID()); (err != nil) != tt.wantErr {
				t.Errorf("BackendRedis.DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, k := range keys {
				_, e := b.Get(k).Result()
				if e == nil {
					t.Error("BackendRedis.DeleteGroup() can't purge cache")
				}
			}
		})
	}
}
