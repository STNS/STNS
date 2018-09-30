package main

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
)

var m sync.Mutex

func etcdTestConfig() *stns.Config {
	return &stns.Config{
		Modules: map[string]interface{}{
			"etcd": map[string]interface{}{
				"endpoints": []interface{}{"http://127.0.0.1:2379"},
				"sync":      true,
			},
		},
	}
}

func TestBackendEtcd_FindUserByID(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		params  model.UserGroup
		want    map[string]model.UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				config: etcdTestConfig(),
			},
			args: args{
				id: 1,
			},
			params: &model.User{
				Base: model.Base{
					ID:   1,
					Name: "user1",
				},
			},
			want: map[string]model.UserGroup{
				"user1": &model.User{
					Base: model.Base{
						ID:   1,
						Name: "user1",
					},
				},
			},
		},
		{
			name: "notfound",
			fields: fields{
				config: etcdTestConfig(),
			},
			args: args{
				id: 99999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			if err := b.Create("/users/name/user1", tt.params); err != nil {
				t.Fatal(err)
			}
			got, err := b.FindUserByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.FindUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendEtcd.FindUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendEtcd_FindUserByName(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		params  model.UserGroup
		want    map[string]model.UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				config: etcdTestConfig(),
			},
			args: args{
				name: "user1",
			},
			params: &model.User{
				Base: model.Base{
					ID:   1,
					Name: "user1",
				},
			},
			want: map[string]model.UserGroup{
				"user1": &model.User{
					Base: model.Base{
						ID:   1,
						Name: "user1",
					},
				},
			},
		},
		{
			name: "notfound",
			fields: fields{
				config: etcdTestConfig(),
			},
			args: args{
				name: "notfound",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			if err := b.Create("/users/name/user1", tt.params); err != nil {
				t.Fatal(err)
			}
			got, err := b.FindUserByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.FindUserByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendEtcd.FindUserByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendEtcd_Users(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	tests := []struct {
		name    string
		fields  fields
		params  model.Users
		want    map[string]model.UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				config: etcdTestConfig(),
			},
			params: model.Users{
				"user1": &model.User{
					Base: model.Base{
						ID:   1,
						Name: "user1",
					},
				},
				"user10": &model.User{
					Base: model.Base{
						ID:   10,
						Name: "user10",
					},
				},
			},
			want: map[string]model.UserGroup{
				"user1": &model.User{
					Base: model.Base{
						ID:   1,
						Name: "user1",
					},
				},
				"user10": &model.User{
					Base: model.Base{
						ID:   10,
						Name: "user10",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range tt.params {
				if err := b.Create(fmt.Sprintf("/users/name/%s", v.Name), v); err != nil {
					t.Fatal(err)
				}
			}
			got, err := b.Users()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.Users() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendEtcd.Users() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendEtcd_FindGroupByID(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		params  model.UserGroup
		want    map[string]model.UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				config: etcdTestConfig(),
			},
			args: args{
				id: 1,
			},
			params: &model.Group{
				Base: model.Base{
					ID:   1,
					Name: "group1",
				},
			},
			want: map[string]model.UserGroup{
				"group1": &model.Group{
					Base: model.Base{
						ID:   1,
						Name: "group1",
					},
				},
			},
		},
		{
			name: "notfound",
			fields: fields{
				config: etcdTestConfig(),
			},
			args: args{
				id: 99999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			if err := b.Create("/groups/name/group1", tt.params); err != nil {
				t.Fatal(err)
			}
			got, err := b.FindGroupByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.FindGroupByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendEtcd.FindGroupByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendEtcd_FindGroupByName(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		params  model.UserGroup
		want    map[string]model.UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				config: etcdTestConfig(),
			},
			args: args{
				name: "group1",
			},
			params: &model.Group{
				Base: model.Base{
					ID:   1,
					Name: "group1",
				},
			},
			want: map[string]model.UserGroup{
				"group1": &model.Group{
					Base: model.Base{
						ID:   1,
						Name: "group1",
					},
				},
			},
		},
		{
			name: "notfound",
			fields: fields{
				config: etcdTestConfig(),
			},
			args: args{
				name: "notfound",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			if err := b.Create("/groups/name/group1", tt.params); err != nil {
				t.Fatal(err)
			}
			got, err := b.FindGroupByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.FindGroupByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendEtcd.FindGroupByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendEtcd_Groups(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	tests := []struct {
		name    string
		fields  fields
		params  model.Groups
		want    map[string]model.UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				config: etcdTestConfig(),
			},
			params: model.Groups{
				"group1": &model.Group{
					Base: model.Base{
						ID:   1,
						Name: "group1",
					},
				},
				"group10": &model.Group{
					Base: model.Base{
						ID:   10,
						Name: "group10",
					},
				},
			},
			want: map[string]model.UserGroup{
				"group1": &model.Group{
					Base: model.Base{
						ID:   1,
						Name: "group1",
					},
				},
				"group10": &model.Group{
					Base: model.Base{
						ID:   10,
						Name: "group10",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range tt.params {
				if err := b.Create(fmt.Sprintf("/groups/name/%s", v.Name), v); err != nil {
					t.Fatal(err)
				}
			}

			got, err := b.Groups()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.Groups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendEtcd.Groups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendEtcd_highlowUserID(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	tests := []struct {
		name    string
		fields  fields
		params  model.Users
		want    int
		wantErr bool
		high    bool
	}{
		{
			name: "high",
			fields: fields{
				config: etcdTestConfig(),
			},
			params: model.Users{
				"user1": &model.User{
					Base: model.Base{
						ID:   1,
						Name: "user1",
					},
				},
				"user10": &model.User{
					Base: model.Base{
						ID:   10,
						Name: "user10",
					},
				},
			},
			want: 10,
			high: true,
		},
		{
			name: "low",
			fields: fields{
				config: etcdTestConfig(),
			},
			params: model.Users{
				"user1": &model.User{
					Base: model.Base{
						ID:   1,
						Name: "user1",
					},
				},
				"user10": &model.User{
					Base: model.Base{
						ID:   10,
						Name: "user10",
					},
				},
			},
			want: 1,
			high: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			b := _b.(BackendEtcd)

			for _, v := range tt.params {
				if err := b.Create(fmt.Sprintf("/users/name/%s", v.Name), v); err != nil {
					t.Fatal(err)
				}
			}

			got := b.highlowUserID(tt.high)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.highlowUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendEtcd.highlowUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendEtcd_highlowGroupID(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	tests := []struct {
		name    string
		fields  fields
		params  model.Groups
		want    int
		wantErr bool
		high    bool
	}{
		{
			name: "high",
			fields: fields{
				config: etcdTestConfig(),
			},
			params: model.Groups{
				"group1": &model.Group{
					Base: model.Base{
						ID:   1,
						Name: "group1",
					},
				},
				"group10": &model.Group{
					Base: model.Base{
						ID:   10,
						Name: "group10",
					},
				},
			},
			want: 10,
			high: true,
		},
		{
			name: "low",
			fields: fields{
				config: etcdTestConfig(),
			},
			params: model.Groups{
				"group1": &model.Group{
					Base: model.Base{
						ID:   1,
						Name: "group1",
					},
				},
				"group10": &model.Group{
					Base: model.Base{
						ID:   10,
						Name: "group10",
					},
				},
			},
			want: 1,
			high: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			b := _b.(BackendEtcd)

			for _, v := range tt.params {
				if err := b.Create(fmt.Sprintf("/groups/name/%s", v.Name), v); err != nil {
					t.Fatal(err)
				}
			}

			got := b.highlowGroupID(tt.high)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.highlowGroupID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendEtcd.highlowGroupID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendEtcd_syncConfig(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	type fields struct {
		config *stns.Config
	}
	tests := []struct {
		name      string
		fields    fields
		wantUser  map[string]model.UserGroup
		wantGroup map[string]model.UserGroup
		wantErr   bool
	}{
		{
			name: "ok",
			fields: fields{
				config: &stns.Config{
					Modules: map[string]interface{}{
						"etcd": map[string]interface{}{
							"endpoints": []interface{}{"http://127.0.0.1:2379"},
							"sync":      true,
						},
					},
					Users: &model.Users{
						"user2": &model.User{
							Base: model.Base{
								ID:   2,
								Name: "user2",
							},
						},
						"user20": &model.User{
							Base: model.Base{
								ID:   20,
								Name: "user20",
							},
						},
					},
					Groups: &model.Groups{
						"group1": &model.Group{
							Base: model.Base{
								ID:   1,
								Name: "group1",
							},
						},
					},
				},
			},
			wantUser: map[string]model.UserGroup{
				"user2": &model.User{
					Base: model.Base{
						ID:   2,
						Name: "user2",
					},
				},
				"user20": &model.User{
					Base: model.Base{
						ID:   20,
						Name: "user20",
					},
				},
			},
			wantGroup: map[string]model.UserGroup{
				"group1": &model.Group{
					Base: model.Base{
						ID:   1,
						Name: "group1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendEtcd(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			back := b.(BackendEtcd)
			if err := back.syncConfig(); (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.syncConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := back.Users()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.Users() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantUser) {
				t.Errorf("BackendEtcd.Users() = %v, want %v", got, tt.wantUser)
			}

			got, err = back.Groups()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendEtcd.Groups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantGroup) {
				t.Errorf("BackendEtcd.Groups() = %v, want %v", got, tt.wantGroup)
			}
		})
	}
}
