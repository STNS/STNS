package stns

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	c, err := NewConfig("test.toml")
	if err != nil {
		t.Fatal(err)
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

	if c.TokenAuth.Tokens[0] != "a" && c.TokenAuth.Tokens[1] != "b" {
		t.Errorf("config cannot parse token auth")
	}

	if c.TLS.Cert != "example_cert" && c.TLS.Key != "example_key" {
		t.Errorf("config cannot parse tls")
	}

	if c.AllowIPs[0] != "10.1.1.1/24" {
		t.Errorf("config cannot parse ip filter")
	}

	yc, err := NewConfig("test.yaml")
	if err != nil {
		t.Fatal(err)
	}
	yu := (*yc.Users)["test"]
	if yu.ID != 20000 {
		t.Errorf("config cannot parse User")
	}
}

func TestNewConfigError(t *testing.T) {
	_, err := NewConfig("../config/test-absent.toml")
	if err == nil {
		t.Errorf("Config cannot handle errors")
	}
}
