package api

import (
	"testing"

	"github.com/pyama86/STNS/attribute"
)

func TestGetByName(t *testing.T) {
	users := map[string]*attribute.All{
		"test1": &attribute.All{
			Id: 1,
		},
	}
	_users := GetByName("test1", users)

	for n, u := range _users {
		if u.Id != 1 {
			t.Error("ummatch user id")
		}

		if n != "test1" {
			t.Error("ummatch user name")
		}
	}
	notfound := GetByName("test2", users)
	if notfound != nil {
		t.Error("ummatch user id")
	}
}
func TestGetById(t *testing.T) {
	users := map[string]*attribute.All{
		"test1": &attribute.All{
			Id: 1,
		},
	}
	_users := GetById("1", users)
	for n, u := range _users {
		if u.Id != 1 {
			t.Error("ummatch user id")
		}

		if n != "test1" {
			t.Error("ummatch user name")
		}
	}

	notfound := GetByName("test2", users)
	if notfound != nil {
		t.Error("ummatch user id")
	}
}
