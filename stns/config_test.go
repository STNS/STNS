package stns

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configDir, err := ioutil.TempDir("", "stns-config-test")
	assertNoError(t, err)

	configFile, err := ioutil.TempFile("", "stns-config-test")
	assertNoError(t, err)

	includedFile, err := os.Create(filepath.Join(configDir, "sub.conf"))

	configContent := fmt.Sprintf(`
port = 9999
include = "%s/*.conf"

[users.example]
id = 1001
group_id = 2001
directory = "/home/example"
shell = "/bin/bash"
keys = ["ssh-rsa aaa"]
link_users = ["example2", "example3"]
`, tomlQuotedReplacer.Replace(configDir))

	includedContent := `
[groups.pepabo]
id = 3001
users = ["example"]
[sudoers.example]
password = "p@ssword"
`
	_, err = configFile.WriteString(configContent)
	assertNoError(t, err)

	_, err = includedFile.WriteString(includedContent)
	assertNoError(t, err)

	configFile.Close()
	includedFile.Close()
	defer os.Remove(configFile.Name())
	defer os.Remove(includedFile.Name())
	name := configFile.Name()

	config, err := LoadConfig(name)
	assertNoError(t, err)
	assert(t, config.Port == 9999, "not over write port")
	assert(t, config.Users["example"].Id == 1001, "unmatch id")
	assert(t, config.Users["example"].GroupId == 2001, "unmatch group id")
	assert(t, config.Users["example"].Directory == "/home/example", "unmatch directory")
	assert(t, config.Users["example"].Shell == "/bin/bash", "unmatch shell")
	assert(t, config.Users["example"].Keys[0] == "ssh-rsa aaa", "unmatch key")
	assert(t, config.Users["example"].LinkUsers[0] == "example2", "unmach link_users")
	assert(t, config.Users["example"].LinkUsers[1] == "example3", "unmach link_users")
	assert(t, config.Groups["pepabo"].Id == 3001, "unmatch group id")
	assert(t, config.Groups["pepabo"].Users[0] == "example", "unmatch group users")
	assert(t, config.Sudoers["example"].Password == "p@ssword", "unmatch password")
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
