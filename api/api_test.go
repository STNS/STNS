package api

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/STNS/config"
)

func TestGet(t *testing.T) {
	loadConfig()
	query := Query{"user", "name", "example1"}
	resource := query.Get()
	assertUserGet(t, resource)

	query = Query{"user", "id", "1001"}
	resource = query.Get()
	assertUserGet(t, resource)

	query = Query{"group", "name", "example_group1"}
	resource = query.Get()
	assertGroupGet(t, resource)

	query = Query{"group", "id", "3001"}
	resource = query.Get()
	assertGroupGet(t, resource)

	query = Query{"sudo", "name", "example1"}
	resource = query.Get()
	assertSudoGet(t, resource)

}

func TestNull(t *testing.T) {
	nullErrorConfig()
	query := Query{"user", "name", "example1"}
	resource := query.Get()
	assert(t, len(resource) == 1, "unmatch resource count null1")

	query = Query{"user", "name", "example2"}
	resource = query.Get()
	assert(t, len(resource) == 1, "unmatch resource count null2")

	query = Query{"user", "name", "example3"}
	resource = query.Get()
	assert(t, len(resource) == 1, "unmatch resource count null2")

	query = Query{"user", "list", ""}
	resource = query.Get()
	assert(t, len(resource) == 3, "unmatch resource count list")

	query = Query{"group", "name", "example1"}
	resource = query.Get()
	assert(t, len(resource) == 1, "unmatch resource count null3")

	query = Query{"group", "name", "example2"}
	resource = query.Get()
	assert(t, len(resource) == 1, "unmatch resource count null4")

	query = Query{"group", "name", "example3"}
	resource = query.Get()
	assert(t, len(resource) == 1, "unmatch resource count null4")

	query = Query{"group", "list", ""}
	resource = query.Get()
	assert(t, len(resource) == 3, "unmatch resource count list")

	query = Query{"sudo", "name", "example1"}
	resource = query.Get()
	assert(t, len(resource) == 1, "unmatch resource count null5")
}

func assertUserGet(t *testing.T, resource attribute.UserGroups) {
	assert(t, len(resource) == 1, "unmatch resource count")
	assert(t, resource["example1"].Id == 1001, "unmatch id")
	assert(t, resource["example1"].Directory == "/home/example1", "unmatch directory")
	assert(t, resource["example1"].Shell == "/bin/bash", "unmatch shell")

	assert(t, len(resource["example1"].Keys) == 3, "unmatch key length")
	sort.Strings(resource["example1"].Keys)
	assert(t, resource["example1"].Keys[0] == "ssh-rsa aaa", "unmatch key1")
	assert(t, resource["example1"].Keys[1] == "ssh-rsa bbb", "unmatch key2")
	assert(t, resource["example1"].Keys[2] == "ssh-rsa ccc", "unmatch key3")
}

func assertGroupGet(t *testing.T, resource attribute.UserGroups) {
	assert(t, len(resource) == 1, "unmatch resource count")
	assert(t, resource["example_group1"].Id == 3001, "unmatch id")
	assert(t, len(resource["example_group1"].Users) == 3, "unmatch group user count")
	sort.Strings(resource["example_group1"].Users)
	assert(t, resource["example_group1"].Users[0] == "example", "unmatch group user1")
	assert(t, resource["example_group1"].Users[1] == "example1", "unmatch group user2")
	assert(t, resource["example_group1"].Users[2] == "example2", "unmatch group user3")
}

func assertSudoGet(t *testing.T, resource attribute.UserGroups) {
	assert(t, len(resource) == 1, "unmatch resource count")
	assert(t, resource["example1"].Password == "p@ssword1", "unmatch password")
}

func TestGetList(t *testing.T) {
	loadConfig()
	query := Query{"user", "list", ""}
	resource := query.Get()
	assert(t, len(resource) == 3, "unmatch resource count")

	assert(t, resource["example1"].Id == 1001, "unmatch id1")
	assert(t, len(resource["example1"].Keys) == 3, "unmatch key length1")
	sort.Strings(resource["example1"].Keys)
	assert(t, resource["example1"].Keys[0] == "ssh-rsa aaa", "unmatch key1")
	assert(t, resource["example1"].Keys[1] == "ssh-rsa bbb", "unmatch key2")
	assert(t, resource["example1"].Keys[2] == "ssh-rsa ccc", "unmatch key3")

	assert(t, resource["example2"].Id == 1002, "unmatch id2")
	assert(t, len(resource["example2"].Keys) == 3, "unmatch key length2")

	assert(t, resource["example3"].Id == 1003, "unmatch id3")
	assert(t, len(resource["example3"].Keys) == 3, "unmatch key length3")

	query = Query{"group", "list", ""}
	resource = query.Get()
	assert(t, len(resource) == 2, "unmatch group resource count")
	assert(t, resource["example_group1"].Id == 3001, "unmatch group id")
	assert(t, len(resource["example_group1"].Users) == 3, "unmatch user count")
	sort.Strings(resource["example_group1"].Users)
	assert(t, resource["example_group1"].Users[0] == "example", "unmatch group user1")
	assert(t, resource["example_group1"].Users[1] == "example1", "unmatch group user2")
	assert(t, resource["example_group1"].Users[2] == "example2", "unmatch group user3")

	assert(t, resource["example_group2"].Id == 3002, "unmatch group  id")
	assert(t, len(resource["example_group2"].Users) == 1, "unmatch group  user count")
	assert(t, resource["example_group2"].Users[0] == "example2", "unmatch group user1")

	query = Query{"sudo", "list", ""}
	resource = query.Get()
	assert(t, len(resource) == 2, "unmatch sudo resource count")
}

func loadConfig() {
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
	config.Load(&name)

}

func nullErrorConfig() {
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
	config.Load(&name)
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assert(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}

var tomlQuotedReplacer = strings.NewReplacer(
	"\t", "\\t",
	"\n", "\\n",
	"\r", "\\r",
	"\"", "\\\"",
	"\\", "\\\\",
)
