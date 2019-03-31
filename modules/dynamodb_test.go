package main

import (
	"reflect"
	"sync"
	"testing"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
)

var dm sync.Mutex

func dynamodbTestConfig() *stns.Config {
	return &stns.Config{
		Modules: map[string]interface{}{
			"dynamodb": map[string]interface{}{
				"sync": true,
			},
		},
	}
}

func TestBackendDynamodb_FindUserByID(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
				config: dynamodbTestConfig(),
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
				config: dynamodbTestConfig(),
			},
			args: args{
				id: 99999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			if tt.params != nil {
				if err := b.CreateUser(tt.params); err != nil {
					t.Fatal(err)
				}
			}
			got, err := b.FindUserByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.FindUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendDynamodb.FindUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendDynamodb_FindUserByName(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
				config: dynamodbTestConfig(),
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
				config: dynamodbTestConfig(),
			},
			args: args{
				name: "notfound",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			if tt.params != nil {
				if err := b.CreateUser(tt.params); err != nil {
					t.Fatal(err)
				}
			}
			got, err := b.FindUserByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.FindUserByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendDynamodb.FindUserByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendDynamodb_Users(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
				config: dynamodbTestConfig(),
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
			b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range tt.params {
				if err := b.CreateUser(v); err != nil {
					t.Fatal(err)
				}
			}
			got, err := b.Users()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.Users() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendDynamodb.Users() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendDynamodb_FindGroupByID(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
				config: dynamodbTestConfig(),
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
				config: dynamodbTestConfig(),
			},
			args: args{
				id: 99999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			if tt.params != nil {
				if err := b.CreateGroup(tt.params); err != nil {
					t.Fatal(err)
				}
			}
			got, err := b.FindGroupByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.FindGroupByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendDynamodb.FindGroupByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendDynamodb_FindGroupByName(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
				config: dynamodbTestConfig(),
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
				config: dynamodbTestConfig(),
			},
			args: args{
				name: "notfound",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			if tt.params != nil {
				if err := b.CreateGroup(tt.params); err != nil {
					t.Fatal(err)
				}
			}
			got, err := b.FindGroupByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.FindGroupByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendDynamodb.FindGroupByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendDynamodb_Groups(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
				config: dynamodbTestConfig(),
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
			b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range tt.params {
				if err := b.CreateGroup(v); err != nil {
					t.Fatal(err)
				}
			}

			got, err := b.Groups()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.Groups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendDynamodb.Groups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendDynamodb_highlowUserID(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
				config: dynamodbTestConfig(),
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
				config: dynamodbTestConfig(),
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
			_b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			b := _b.(BackendDynamoDB)

			for _, v := range tt.params {
				if err := b.CreateUser(v); err != nil {
					t.Fatal(err)
				}
			}

			got := 0
			if tt.high {
				got = b.HighestUserID()
			} else {
				got = b.LowestUserID()
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.highlowUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendDynamodb.highlowUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendDynamodb_highlowGroupID(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
				config: dynamodbTestConfig(),
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
				config: dynamodbTestConfig(),
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
			_b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			b := _b.(BackendDynamoDB)

			for _, v := range tt.params {
				if err := b.CreateGroup(v); err != nil {
					t.Fatal(err)
				}
			}

			got := 0
			if tt.high {
				got = b.HighestGroupID()
			} else {
				got = b.LowestGroupID()
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.highlowGroupID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendDynamodb.highlowGroupID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendDynamodb_syncConfig(t *testing.T) {
	dm.Lock()
	defer dm.Unlock()
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
						"dynamodb": map[string]interface{}{
							"sync": true,
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
			b, err := NewBackendDynamodb(tt.fields.config)
			if err != nil {
				t.Fatal(err)
			}
			back := b.(BackendDynamoDB)

			// check modify result
			if err := b.CreateUser(
				&model.User{
					Base: model.Base{
						ID:   2,
						Name: "before modify",
					},
				},
			); err != nil {
				t.Fatal(err)
			}

			if err := syncConfig(back, tt.fields.config); (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.syncConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := back.Users()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.Users() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantUser) {
				t.Errorf("BackendDynamodb.Users() = %v, want %v", got, tt.wantUser)
			}

			got, err = back.Groups()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendDynamodb.Groups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantGroup) {
				t.Errorf("BackendDynamodb.Groups() = %v, want %v", got, tt.wantGroup)
			}
		})
	}
}
