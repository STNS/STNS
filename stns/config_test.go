package stns

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("./fixtures/base.conf")
	assertNoError(t, err)
	assert(t, config.Port == 9999, "not over write port")
	assert(t, config.TLSCa == "ca.pem", "unmatch tls ca")
	assert(t, config.TLSCert == "tls.crt", "unmatch tls crt")
	assert(t, config.TLSKey == "tls.key", "unmatch tls key")

	assert(t, config.Users["example"].ID == 1001, "unmatch id")
	assert(t, config.Users["example"].GroupID == 2002, "unmatch group id")
	assert(t, config.Users["example"].Directory == "/home/example", "unmatch directory")
	assert(t, config.Users["example"].Shell == "/bin/bash", "unmatch shell")
	assert(t, config.Users["example"].Keys[0] == "ssh-rsa aaa", "unmatch key")
	assert(t, config.Users["example"].LinkUsers[0] == "example2", "unmach link_users")
	assert(t, config.Users["example"].LinkUsers[1] == "example3", "unmach link_users")
	assert(t, config.Users["example"].SetupCommands[0] == "commands", "unmatch commands")
	assert(t, config.Groups["pepabo"].ID == 3001, "unmatch group id")
	assert(t, config.Groups["pepabo"].Users[0] == "example", "unmatch group users")
	assert(t, config.Sudoers["example"].Password == "p@ssword", "unmatch password")
	assert(t, config.Users["example"].LinkUsers[1] == "example3", "unmach link_users")
}

func TestDuplicateID(t *testing.T) {
	_, err := LoadConfig("./fixtures/duplicate_id.conf")
	assert(t, err.Error() == "Duplicate id is not allowed user_id:1001", "TestDuplicateID")
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
