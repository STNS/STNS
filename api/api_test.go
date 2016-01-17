package api

import (
	"testing"

	"github.com/pyama86/STNS/attribute"
)

func TestGetByName(t *testing.T) {
	users := attribute.UserGroups{
		"test1": &attribute.All{
			Id: 1,
		},
	}
	_users := users.GetByName("test1")

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
func TestGetById(t *testing.T) {
	users := attribute.UserGroups{
		"test1": &attribute.All{
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
