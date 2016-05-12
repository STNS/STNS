package stns

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("./fixtures/base.conf")

	assertNoError(t, err)
	assert(t, config.Port == 9999, "not over write port")
	assert(t, config.Salt == true, "not ovver write salt")
	assert(t, config.Stretching == 1000, "not ovver write stretching_count")
	assert(t, config.Users["example"].Id == 1001, "unmatch id")
	assert(t, config.Users["example"].GroupId == 2002, "unmatch group id")
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
