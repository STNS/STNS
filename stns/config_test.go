package stns

import "testing"

func TestNewConfig(t *testing.T) {
	c, err := NewConfig("test.toml")
	if err != nil {
		t.Fatalf(err.Error())
	}

	u := c.Users["test"]
	if u.Name != "test" {
		t.Errorf("Config cannot parse User")
	}

	g := c.Groups["test"]
	if g.Name != "test" {
		t.Errorf("Config cannot parse Group")
	}
}

func TestNewConfigError(t *testing.T) {
	_, err := NewConfig("../config/test-absent.toml")
	if err == nil {
		t.Errorf("Config cannot handle errors")
	}
}
