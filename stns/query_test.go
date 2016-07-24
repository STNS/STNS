package stns

import (
	"sort"
	"testing"

	"github.com/STNS/STNS/test"
)

func TestGet(t *testing.T) {
	config, _ := LoadConfig("./fixtures/query_01.conf")
	query := Query{&config, nil, "user", "name", "example1"}
	resource := query.Get()
	assertUserGet(t, resource)

	query = Query{&config, nil, "user", "id", "1001"}
	resource = query.Get()
	assertUserGet(t, resource)

	query = Query{&config, nil, "group", "name", "example_group1"}
	resource = query.Get()
	assertGroupGet(t, resource)

	query = Query{&config, nil, "group", "id", "3001"}
	resource = query.Get()
	assertGroupGet(t, resource)

	query = Query{&config, nil, "sudo", "name", "example1"}
	resource = query.Get()
	assertSudoGet(t, resource)

}

func TestNull(t *testing.T) {
	config, _ := LoadConfig("./fixtures/query_02.conf")
	query := Query{&config, nil, "user", "name", "example1"}
	resource := query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null 1")

	query = Query{&config, nil, "user", "name", "example2"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null 2")

	query = Query{&config, nil, "user", "name", "example3"}
	resource = query.Get()
	test.Assert(t, len(resource) == 0, "unmatch resource count null 3")

	query = Query{&config, nil, "user", "list", ""}
	resource = query.Get()
	test.Assert(t, len(resource) == 3, "unmatch resource count list 1")

	query = Query{&config, nil, "group", "name", "example1"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null 4")

	query = Query{&config, nil, "group", "name", "example2"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null 5")

	query = Query{&config, nil, "group", "name", "example3"}
	resource = query.Get()
	test.Assert(t, len(resource) == 0, "unmatch resource count null 6")

	query = Query{&config, nil, "group", "list", ""}
	resource = query.Get()
	test.Assert(t, len(resource) == 3, "unmatch resource count list 2")

	query = Query{&config, nil, "sudo", "name", "example1"}
	resource = query.Get()
	test.Assert(t, len(resource) == 0, "unmatch resource count null 7")

	nullconfig, _ := LoadConfig("./fixtures/query_03.conf")
	nq := Query{&nullconfig, nil, "user", "name", "example1"}
	resource = nq.Get()
	test.Assert(t, len(resource) == 0, "unmatch resource count null 8")

}

func assertUserGet(t *testing.T, resource Attributes) {
	test.Assert(t, len(resource) == 1, "unmatch resource count")
	test.Assert(t, resource["example1"].Id == 1001, "unmatch id")
	test.Assert(t, resource["example1"].Directory == "/home/example1", "unmatch directory")
	test.Assert(t, resource["example1"].Shell == "/bin/bash", "unmatch shell")

	test.Assert(t, len(resource["example1"].Keys) == 3, "unmatch key length")
	sort.Strings(resource["example1"].Keys)
	test.Assert(t, resource["example1"].Keys[0] == "ssh-rsa aaa", "unmatch key1")
	test.Assert(t, resource["example1"].Keys[1] == "ssh-rsa bbb", "unmatch key2")
	test.Assert(t, resource["example1"].Keys[2] == "ssh-rsa ccc", "unmatch key3")
}

func assertGroupGet(t *testing.T, resource Attributes) {
	test.Assert(t, len(resource) == 1, "unmatch resource count")
	test.Assert(t, resource["example_group1"].Id == 3001, "unmatch id")
	test.Assert(t, len(resource["example_group1"].Users) == 3, "unmatch group user count")
	sort.Strings(resource["example_group1"].Users)
	test.Assert(t, resource["example_group1"].Users[0] == "example", "unmatch group user1")
	test.Assert(t, resource["example_group1"].Users[1] == "example1", "unmatch group user2")
	test.Assert(t, resource["example_group1"].Users[2] == "example2", "unmatch group user3")
}

func assertSudoGet(t *testing.T, resource Attributes) {
	test.Assert(t, len(resource) == 1, "unmatch resource count")
	test.Assert(t, resource["example1"].Password == "p@ssword1", "unmatch password")
}

func TestGetList(t *testing.T) {
	config, _ := LoadConfig("./fixtures/query_01.conf")
	query := Query{&config, nil, "user", "list", ""}
	resource := query.Get()
	test.Assert(t, len(resource) == 3, "unmatch resource count")

	test.Assert(t, resource["example1"].Id == 1001, "unmatch id1")
	test.Assert(t, len(resource["example1"].Keys) == 3, "unmatch key length1")
	sort.Strings(resource["example1"].Keys)
	test.Assert(t, resource["example1"].Keys[0] == "ssh-rsa aaa", "unmatch key1")
	test.Assert(t, resource["example1"].Keys[1] == "ssh-rsa bbb", "unmatch key2")
	test.Assert(t, resource["example1"].Keys[2] == "ssh-rsa ccc", "unmatch key3")

	test.Assert(t, resource["example2"].Id == 1002, "unmatch id2")
	test.Assert(t, len(resource["example2"].Keys) == 3, "unmatch key length2")

	test.Assert(t, resource["example3"].Id == 1003, "unmatch id3")
	test.Assert(t, len(resource["example3"].Keys) == 3, "unmatch key length3")

	query = Query{&config, nil, "group", "list", ""}
	resource = query.Get()
	test.Assert(t, len(resource) == 2, "unmatch group resource count")
	test.Assert(t, resource["example_group1"].Id == 3001, "unmatch group id")
	test.Assert(t, len(resource["example_group1"].Users) == 3, "unmatch user count")
	sort.Strings(resource["example_group1"].Users)
	test.Assert(t, resource["example_group1"].Users[0] == "example", "unmatch group user1")
	test.Assert(t, resource["example_group1"].Users[1] == "example1", "unmatch group user2")
	test.Assert(t, resource["example_group1"].Users[2] == "example2", "unmatch group user3")

	test.Assert(t, resource["example_group2"].Id == 3002, "unmatch group  id")
	test.Assert(t, len(resource["example_group2"].Users) == 1, "unmatch group  user count")
	test.Assert(t, resource["example_group2"].Users[0] == "example2", "unmatch group user1")

	query = Query{&config, nil, "sudo", "list", ""}
	resource = query.Get()
	test.Assert(t, len(resource) == 2, "unmatch sudo resource count")
}
