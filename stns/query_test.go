package stns

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/STNS/STNS/test"
)

func TestGet(t *testing.T) {
	config := fullConfig()
	query := Query{config, "user", "name", "example1"}
	resource := query.Get()
	assertUserGet(t, resource)

	query = Query{config, "user", "id", "1001"}
	resource = query.Get()
	assertUserGet(t, resource)

	query = Query{config, "group", "name", "example_group1"}
	resource = query.Get()
	assertGroupGet(t, resource)

	query = Query{config, "group", "id", "3001"}
	resource = query.Get()
	assertGroupGet(t, resource)

	query = Query{config, "sudo", "name", "example1"}
	resource = query.Get()
	assertSudoGet(t, resource)

}

func TestNull(t *testing.T) {
	config := nullErrorConfig()
	query := Query{config, "user", "name", "example1"}
	resource := query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null1")

	query = Query{config, "user", "name", "example2"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null2")

	query = Query{config, "user", "name", "example3"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null2")

	query = Query{config, "user", "list", ""}
	resource = query.Get()
	test.Assert(t, len(resource) == 3, "unmatch resource count list")

	query = Query{config, "group", "name", "example1"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null3")

	query = Query{config, "group", "name", "example2"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null4")

	query = Query{config, "group", "name", "example3"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null4")

	query = Query{config, "group", "list", ""}
	resource = query.Get()
	test.Assert(t, len(resource) == 3, "unmatch resource count list")

	query = Query{config, "sudo", "name", "example1"}
	resource = query.Get()
	test.Assert(t, len(resource) == 1, "unmatch resource count null5")
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
	config := fullConfig()
	query := Query{config, "user", "list", ""}
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

	query = Query{config, "group", "list", ""}
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

	query = Query{config, "sudo", "list", ""}
	resource = query.Get()
	test.Assert(t, len(resource) == 2, "unmatch sudo resource count")
}

func fullConfig() *Config {
	configFile, _ := ioutil.TempFile("", "stns-config-test")
	configContent := `port = 9999
[users.example1]
id = 1001
group_id = 2001
directory = "/home/example1"
shell = "/bin/bash"
keys = ["ssh-rsa aaa"]
link_users = ["example2"]

[users.example2]
id = 1002
group_id = 2001
directory = "/home/example2"
shell = "/bin/bash"
keys = ["ssh-rsa bbb"]
link_users = ["example3"]

[users.example3]
id = 1003
group_id = 2001
directory = "/home/example3"
shell = "/bin/bash"
keys = ["ssh-rsa ccc"]
link_users = ["example1"]

[groups.example_group1]
id = 3001
users = ["example", "example1"]
link_groups = ["example_group2"]

[groups.example_group2]
id = 3002
users = ["example2"]

[sudoers.example1]
password = "p@ssword1"

[sudoers.example2]
password = "p@ssword2"
`
	_, _ = configFile.WriteString(configContent)
	configFile.Close()
	defer os.Remove(configFile.Name())
	name := configFile.Name()
	config, _ := LoadConfig(name)
	return config
}

func nullErrorConfig() *Config {
	configFile, _ := ioutil.TempFile("", "stns-config-test")
	configContent := `port = 9999
[users.example1]
link_users = ["none"]

[users.example2]
id = 1001

[users.example3]

[groups.example1]
link_groups = ["none"]

[groups.example2]
id = 1001

[groups.example3]

[sudoers.example1]
`
	_, _ = configFile.WriteString(configContent)
	configFile.Close()
	defer os.Remove(configFile.Name())
	name := configFile.Name()
	config, _ := LoadConfig(name)
	return config
}
