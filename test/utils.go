package test

import (
	"strings"
	"testing"
)

// AssertNoError util method
func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

// Assert util method
func Assert(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}

// TomlQuotedReplacer util method
var TomlQuotedReplacer = strings.NewReplacer(
	"\t", "\\t",
	"\n", "\\n",
	"\r", "\\r",
	"\"", "\\\"",
	"\\", "\\\\",
)
