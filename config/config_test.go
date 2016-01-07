package config

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

[users.pyama]
id = 1001
group_id = 2001
directory = "/home/pyama"
shell = "/bin/bash"
keys = ["ssh-rsa aaa"]
`, tomlQuotedReplacer.Replace(configDir))

	includedContent := `
[groups.pepabo]
id = 3001
users = ["pyama"]
`
	_, err = configFile.WriteString(configContent)
	assertNoError(t, err)

	_, err = includedFile.WriteString(includedContent)
	assertNoError(t, err)

	configFile.Close()
	includedFile.Close()
	defer os.Remove(configFile.Name())
	defer os.Remove(includedFile.Name())
	err = Load(configFile.Name())
	assertNoError(t, err)
	assert(t, All.Port == 9999, "not over write port")
	assert(t, All.Users["pyama"].Id == 1001, "unmatch id")
	assert(t, All.Users["pyama"].GroupId == 2001, "unmatch group id")
	assert(t, All.Users["pyama"].Directory == "/home/pyama", "unmatch directory")
	assert(t, All.Users["pyama"].Shell == "/bin/bash", "unmatch shell")
	assert(t, All.Users["pyama"].Keys[0] == "ssh-rsa aaa", "unmatch key")
	assert(t, All.Groups["pepabo"].Id == 3001, "unmatch group id")
	assert(t, All.Groups["pepabo"].Users[0] == "pyama", "unmatch group users")
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
