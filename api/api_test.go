package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/STNS/STNS/config"
)

func TestGet(t *testing.T) {
	loadConfig()
	query := Query{"user", "name", "pyama1"}
	resource := query.Get()

	assert(t, len(resource) == 1, "unmach resource count")
	assert(t, resource["pyama1"].Id == 1001, "unmach id")
	assert(t, len(resource["pyama1"].Keys) == 3, "unmatch key length")
	sort.Strings(resource["pyama1"].Keys)
	assert(t, resource["pyama1"].Keys[0] == "ssh-rsa aaa", "unmach key1")
	assert(t, resource["pyama1"].Keys[1] == "ssh-rsa bbb", "unmach key2")
	assert(t, resource["pyama1"].Keys[2] == "ssh-rsa ccc", "unmach key3")
}
func TestGetList(t *testing.T) {
	loadConfig()
	query := Query{"user", "list", ""}
	resource := query.Get()
	assert(t, len(resource) == 5, "unmach resource count")

	assert(t, resource["pyama1"].Id == 1001, "unmach id1")
	assert(t, len(resource["pyama1"].Keys) == 3, "unmatch key length1")
	sort.Strings(resource["pyama1"].Keys)
	assert(t, resource["pyama1"].Keys[0] == "ssh-rsa aaa", "unmach key1")
	assert(t, resource["pyama1"].Keys[1] == "ssh-rsa bbb", "unmach key2")
	assert(t, resource["pyama1"].Keys[2] == "ssh-rsa ccc", "unmach key3")

	fmt.Println(resource["pyama2"].Keys)
	assert(t, resource["pyama2"].Id == 1002, "unmach id2")
	assert(t, len(resource["pyama2"].Keys) == 3, "unmatch key length2")

	assert(t, resource["pyama3"].Id == 1003, "unmach id3")
	assert(t, len(resource["pyama3"].Keys) == 3, "unmatch key length3")

	assert(t, resource["pyama4"].Id == 1004, "unmach id4")
	assert(t, len(resource["pyama4"].Keys) == 2, "unmatch key length4")
	sort.Strings(resource["pyama4"].Keys)
	assert(t, resource["pyama4"].Keys[0] == "ssh-rsa ddd", "unmach key1")
	assert(t, resource["pyama4"].Keys[1] == "ssh-rsa eee", "unmach key2")
}

func loadConfig() {
	configFile, _ := ioutil.TempFile("", "stns-config-test")
	configContent := `port = 9999
[users.pyama1]
id = 1001
group_id = 2001
directory = "/home/pyama1"
shell = "/bin/bash"
keys = ["ssh-rsa aaa"]
link_users = ["pyama2"]

[users.pyama2]
id = 1002
group_id = 2001
directory = "/home/pyama2"
shell = "/bin/bash"
keys = ["ssh-rsa bbb"]
link_users = ["pyama3"]

[users.pyama3]
id = 1003
group_id = 2001
directory = "/home/pyama3"
shell = "/bin/bash"
keys = ["ssh-rsa ccc"]
link_users = ["pyama1"]

[users.pyama4]
id = 1004
group_id = 2001
directory = "/home/pyama4"
shell = "/bin/bash"
keys = ["ssh-rsa ddd"]
link_users = ["pyama5"]

[users.pyama5]
id = 1005
group_id = 2001
directory = "/home/pyama5"
shell = "/bin/bash"
keys = ["ssh-rsa eee"]
link_users = ["pyama4"]

[groups.pepabo]
id = 3001
users = ["pyama"]
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
