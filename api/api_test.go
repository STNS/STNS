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

	query = Query{"group", "name", "example_group"}
	resource = query.Get()
	assertGroupGet(t, resource)

	query = Query{"group", "id", "3001"}
	resource = query.Get()
	assertGroupGet(t, resource)
}

func assertUserGet(t *testing.T, resource attribute.UserGroups) {
	assert(t, len(resource) == 1, "unmach resource count")
	assert(t, resource["example1"].Id == 1001, "unmach id")
	assert(t, len(resource["example1"].Keys) == 3, "unmatch key length")
	sort.Strings(resource["example1"].Keys)
	assert(t, resource["example1"].Keys[0] == "ssh-rsa aaa", "unmach key1")
	assert(t, resource["example1"].Keys[1] == "ssh-rsa bbb", "unmach key2")
	assert(t, resource["example1"].Keys[2] == "ssh-rsa ccc", "unmach key3")
}

func assertGroupGet(t *testing.T, resource attribute.UserGroups) {
	assert(t, len(resource) == 1, "unmach resource count")
	assert(t, resource["example_group"].Id == 3001, "unmach id")
	assert(t, len(resource["example_group"].Users) == 2, "unmach user count")
	assert(t, resource["example_group"].Users[0] == "example", "unmach user1")
	assert(t, resource["example_group"].Users[1] == "example1", "unmach user2")
}

func TestGetList(t *testing.T) {
	loadConfig()
	query := Query{"user", "list", ""}
	resource := query.Get()
	assert(t, len(resource) == 5, "unmach resource count")

	assert(t, resource["example1"].Id == 1001, "unmach id1")
	assert(t, len(resource["example1"].Keys) == 3, "unmatch key length1")
	sort.Strings(resource["example1"].Keys)
	assert(t, resource["example1"].Keys[0] == "ssh-rsa aaa", "unmach key1")
	assert(t, resource["example1"].Keys[1] == "ssh-rsa bbb", "unmach key2")
	assert(t, resource["example1"].Keys[2] == "ssh-rsa ccc", "unmach key3")

	assert(t, resource["example2"].Id == 1002, "unmach id2")
	assert(t, len(resource["example2"].Keys) == 3, "unmatch key length2")

	assert(t, resource["example3"].Id == 1003, "unmach id3")
	assert(t, len(resource["example3"].Keys) == 3, "unmatch key length3")

	assert(t, resource["example4"].Id == 1004, "unmach id4")
	assert(t, len(resource["example4"].Keys) == 2, "unmatch key length4")
	sort.Strings(resource["example4"].Keys)
	assert(t, resource["example4"].Keys[0] == "ssh-rsa ddd", "unmach key1")
	assert(t, resource["example4"].Keys[1] == "ssh-rsa eee", "unmach key2")
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

[users.example4]
id = 1004
group_id = 2001
directory = "/home/example4"
shell = "/bin/bash"
keys = ["ssh-rsa ddd"]
link_users = ["example5"]

[users.example5]
id = 1005
group_id = 2001
directory = "/home/example5"
shell = "/bin/bash"
keys = ["ssh-rsa eee"]
link_users = ["example4"]

[groups.example_group]
id = 3001
users = ["example", "example1"]
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
