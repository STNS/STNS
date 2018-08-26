package stns

import "testing"

func TestNewConfig(t *testing.T) {
	c, err := NewConfig("test.toml")
	if err != nil {
		t.Fatalf(err.Error())
	}

	u := (*c.Users)["test"]
	if u.ID != 10001 {
		t.Errorf("config cannot parse User")
	}

	// include config user
	i := (*c.Users)["bar"]
	if i.ID != 99 {
		t.Errorf("include config cannot parse User")
	}

	g := (*c.Groups)["test"]
	if g.ID != 10001 {
		t.Errorf("config cannot parse Group")
	}

	if c.BasicAuth.User != "foo" && c.BasicAuth.Password != "bar" {
		t.Errorf("config cannot parse basic auth")
	}
}

func TestNewConfigError(t *testing.T) {
	_, err := NewConfig("../config/test-absent.toml")
	if err == nil {
		t.Errorf("Config cannot handle errors")
	}
}
