package model

import (
	"reflect"
	"testing"
)

func Test_tomlFileFindByID(t *testing.T) {
	type args struct {
		id   int
		list map[string]UserGroup
	}
	tests := []struct {
		name string
		args args
		want map[string]UserGroup
	}{
		{
			name: "ok",
			args: args{
				id: 1,
				list: map[string]UserGroup{
					"test": &User{
						Base: Base{
							ID:   1,
							Name: "user1",
						},
					},
				},
			},
			want: map[string]UserGroup{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "user1",
					},
				},
			},
		},
		{
			name: "not found",
			args: args{
				id: 2,
				list: map[string]UserGroup{
					"test": &User{
						Base: Base{
							ID:   1,
							Name: "user1",
						},
					},
				},
			},
			want: map[string]UserGroup{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tomlFileFindByID(tt.args.id, tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tomlFileFindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tomlFileFindByName(t *testing.T) {
	type args struct {
		name string
		list map[string]UserGroup
	}
	tests := []struct {
		name string
		args args
		want map[string]UserGroup
	}{
		{
			name: "ok",
			args: args{
				name: "user1",
				list: map[string]UserGroup{
					"test": &User{
						Base: Base{
							ID:   1,
							Name: "user1",
						},
					},
				},
			},
			want: map[string]UserGroup{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "user1",
					},
				},
			},
		},
		{
			name: "not found",
			args: args{
				name: "hoge",
				list: map[string]UserGroup{
					"test": &User{
						Base: Base{
							ID:   1,
							Name: "user1",
						},
					},
				},
			},
			want: map[string]UserGroup{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tomlFileFindByName(tt.args.name, tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tomlFileFindByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeLinkAttribute(t *testing.T) {
	type args struct {
		las    linkAttributers
		result map[string][]string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "user ok",
			args: args{
				las: linkAttributers{
					"user1": &User{
						Base: Base{
							Name: "user1",
						},
						Keys:      []string{"user1key"},
						LinkUsers: []string{"user2", "user3"},
					},
					"user2": &User{
						Base: Base{
							Name: "user2",
						},
						Keys:      []string{"user2key"},
						LinkUsers: []string{"user3"},
					},
					"user3": &User{
						Base: Base{
							Name: "user3",
						},
						Keys:      []string{"user3key"},
						LinkUsers: []string{"user1", "user2"},
					},
					"user4": &User{
						Base: Base{
							Name: "user4",
						},
						Keys:      []string{"user4key"},
						LinkUsers: []string{"user1", "user2", "user3"},
					},
				},
				result: map[string][]string{},
			},
			want: map[string][]string{
				"user1": []string{
					"user1key",
					"user2key",
					"user3key",
				},
				"user2": []string{
					"user1key",
					"user2key",
					"user3key",
				},
				"user3": []string{
					"user1key",
					"user2key",
					"user3key",
				},
				"user4": []string{
					"user1key",
					"user2key",
					"user3key",
					"user4key",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeLinkAttribute(tt.args.las, nil, nil, nil)
			for k, v := range tt.want {
				if !reflect.DeepEqual(tt.args.las[k].value(), v) {
					t.Errorf("mergeLinkAttribute() = user %v keys %v, want %v", k, tt.args.las[k].value(), v)
				}
			}
		})
	}
}

func Test_tomlHighLowID(t *testing.T) {
	type args struct {
		list      map[string]UserGroup
		highorlow int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "highest",
			args: args{
				highorlow: 0,
				list: map[string]UserGroup{
					"low": &User{
						Base: Base{
							ID:   1,
							Name: "user1",
						},
					},
					"high": &User{
						Base: Base{
							ID:   2,
							Name: "user2",
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "lowest",
			args: args{
				highorlow: 1,
				list: map[string]UserGroup{
					"low": &User{
						Base: Base{
							ID:   1,
							Name: "user1",
						},
					},
					"high": &User{
						Base: Base{
							ID:   2,
							Name: "user2",
						},
					},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tomlHighLowID(tt.args.highorlow, tt.args.list); got != tt.want {
				t.Errorf("tomlHighLowID() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestBackendTomlFile_Users(t *testing.T) {
	tests := []struct {
		name    string
		users   Users
		want    map[string]UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			users: Users{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "user1",
					},
				},
			},
			want: map[string]UserGroup{
				"test": &User{
					Base: Base{
						ID:   1,
						Name: "user1",
					},
				},
			},
		},
		{
			name:    "not found",
			users:   Users{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BackendTomlFile{
				users: &tt.users,
			}
			got, err := b.Users()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendTomlFile.Users() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendTomlFile.Users() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackendTomlFile_Groups(t *testing.T) {
	tests := []struct {
		name    string
		groups  Groups
		want    map[string]UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			groups: Groups{
				"test": &Group{
					Base: Base{
						ID:   1,
						Name: "group1",
					},
				},
			},
			want: map[string]UserGroup{
				"test": &Group{
					Base: Base{
						ID:   1,
						Name: "group1",
					},
				},
			},
		},
		{
			name:    "not found",
			groups:  Groups{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BackendTomlFile{
				groups: &tt.groups,
			}
			got, err := b.Groups()
			if (err != nil) != tt.wantErr {
				t.Errorf("BackendTomlFile.Groups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackendTomlFile.Groups() = %v, want %v", got, tt.want)
			}
		})
	}
}
