package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/redis.v5"
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
			redis := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			})
			b := &BackendRedis{
				backend: BackendDummy{},
				Client:  redis,
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
			redis := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			})
			b := &BackendRedis{
				backend: BackendDummy{},
				Client:  redis,
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
			redis := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			})
			b := &BackendRedis{
				backend: BackendDummy{},
				Client:  redis,
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
			redis := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			})
			b := &BackendRedis{
				backend: BackendDummy{},
				Client:  redis,
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
			redis := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			})
			b := &BackendRedis{
				backend: BackendDummy{},
				Client:  redis,
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
			redis := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			})
			b := &BackendRedis{
				backend: BackendDummy{},
				Client:  redis,
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
		})
	}
}
