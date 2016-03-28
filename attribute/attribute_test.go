package attribute

import (
	"reflect"
	"testing"
)

func TestGetByName(t *testing.T) {
	users := AllAttribute{
		"test1": &All{
			Id: 1,
			User: &User{LinkUsers: []string{"foo", "bar"},
				Password: "foo",
			},
			Group: &Group{Users: []string{"foo", "bar"}},
		},
	}
	_users := users.GetByName("test1")

	for n, u := range _users {
		if u.Id != 1 {
			t.Error("ummatch user id")
		}

		if !reflect.DeepEqual(u.LinkUsers, []string{"foo", "bar"}) {
			t.Error("ummatch link user")
		}

		if !reflect.DeepEqual(u.Users, []string{"foo", "bar"}) {
			t.Error("ummatch link user")
		}

		if u.Password != "foo" {
			t.Error("ummatch password")
		}

		if n != "test1" {
			t.Error("ummatch user name")
		}
	}
	notfound := users.GetByName("test2")
	if notfound != nil {
		t.Error("ummatch user id")
	}
}
func TestGetById(t *testing.T) {
	users := AllAttribute{
		"test1": &All{
			Id: 1,
		},
	}
	_users := users.GetById("1")
	for n, u := range _users {
		if u.Id != 1 {
			t.Error("ummatch user id")
		}

		if n != "test1" {
			t.Error("ummatch user name")
		}
	}

	notfound := users.GetByName("test2")
	if notfound != nil {
		t.Error("ummatch user id")
	}
}
