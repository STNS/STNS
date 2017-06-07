package stns

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/STNS/STNS/test"
)

func TestGetByName(t *testing.T) {
	users := Attributes{
		"test1": &Attribute{
			ID: 1,
			User: &User{LinkUsers: []string{"foo", "bar"},
				Password: "foo",
			},
			Group: &Group{Users: []string{"foo", "bar"}},
		},
		"test3": &Attribute{},
		"test4": &Attribute{
			ID: 4,
		},
	}

	_users := users.GetByName("test1")

	for n, u := range _users {
		if u.ID != 1 {
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

	t3 := users.GetByName("test3")

	if t3 != nil {
		t.Error("ummatch user id test3")
	}

	t4 := users.GetByName("test4")
	if t4 == nil {
		t.Error("ummatch user id test4")
	}
}

func TestGetByID(t *testing.T) {
	users := Attributes{
		"test1": &Attribute{
			ID: 1,
		},
		"test3": &Attribute{},
		"test4": &Attribute{
			ID: 4,
		},
	}
	_users := users.GetByID("1")
	for n, u := range _users {
		if u.ID != 1 {
			t.Error("ummatch user id")
		}

		if n != "test1" {
			t.Error("ummatch user name")
		}
	}

	notfound := users.GetByID("2")
	if notfound != nil {
		t.Error("ummatch user id")
	}

	t3 := users.GetByID("3")

	if t3 != nil {
		t.Error("ummatch user id test3")
	}

	t4 := users.GetByID("4")
	if t4 == nil {
		t.Error("ummatch user id test4")
	}
}

func TestAttributePrevID(t *testing.T) {
	users := Attributes{
		"test1": &Attribute{ID: 1},
		"test2": &Attribute{ID: 3},
		"test3": &Attribute{ID: 5},
	}

	test.Assert(t, 1 == users.PrevID(1), fmt.Sprintf("AttributePrevID expected: %d got: %d", 1, users.PrevID(1)))
	test.Assert(t, 1 == users.PrevID(2), fmt.Sprintf("AttributePrevID expected: %d got: %d", 1, users.PrevID(2)))
	test.Assert(t, 3 == users.PrevID(5), fmt.Sprintf("AttributePrevID expected: %d got: %d", 3, users.PrevID(5)))
	test.Assert(t, 5 == users.PrevID(6), fmt.Sprintf("AttributePrevID expected: %d got: %d", 5, users.PrevID(6)))
}

func TestAttributeNextID(t *testing.T) {
	users := Attributes{
		"test1": &Attribute{ID: 1},
		"test2": &Attribute{ID: 3},
		"test3": &Attribute{ID: 5},
	}

	test.Assert(t, 3 == users.NextID(1), fmt.Sprintf("AttributeNextID expected: %d got: %d", 3, users.NextID(1)))
	test.Assert(t, 3 == users.NextID(2), fmt.Sprintf("AttributeNextID expected: %d got: %d", 3, users.NextID(2)))
	test.Assert(t, 5 == users.NextID(4), fmt.Sprintf("AttributeNextID expected: %d got: %d", 5, users.NextID(4)))
	test.Assert(t, 6 == users.NextID(6), fmt.Sprintf("AttributeNextID expected: %d got: %d", 6, users.NextID(6)))
}
