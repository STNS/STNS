package modules

import (
	"reflect"
	"testing"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
)

func TestBackendEtcd_FindUserByID(t *testing.T) {
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
		want    map[string]model.UserGroup
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				config: &stns.Config{
					Etcd: &stns.Etcd{
						Endpoints: []string{"http://127.0.0.1:2379"},
					},
				},
			},
			args: args{
				id: 1,
			},
			want: map[string]model.UserGroup{
				"test": &model.User{
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
				config: &stns.Config{
					Etcd: &stns.Etcd{
						Endpoints: []string{"http://127.0.0.1:2379"},
					},
				},
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
			if err := b.Create("/users/id/1", tt.want); err != nil {
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
