package test

import (
	"strings"
	"testing"
)

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func Assert(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}

var TomlQuotedReplacer = strings.NewReplacer(
	"\t", "\\t",
	"\n", "\\n",
	"\r", "\\r",
	"\"", "\\\"",
	"\\", "\\\\",
)
