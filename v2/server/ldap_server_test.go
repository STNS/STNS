package server

import (
	"net"
	"testing"

	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/stns"
	"github.com/labstack/gommon/log"
	"github.com/nmcclain/ldap"
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

type userBackend struct {
	model.Backend
	users map[string]model.UserGroup
}

func (b userBackend) FindUserByName(name string) (map[string]model.UserGroup, error) {
	return b.users, nil
}

func TestLDAPHandlerBind(t *testing.T) {
	hash, err := sha512_crypt.New().Generate([]byte("secret"), nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		storedHash  string
		rawPassword string
		want        ldap.LDAPResultCode
	}{
		{
			name:        "valid password",
			storedHash:  hash,
			rawPassword: "secret",
			want:        ldap.LDAPResultSuccess,
		},
		{
			name:        "wrong password",
			storedHash:  hash,
			rawPassword: "wrong",
			want:        ldap.LDAPResultInvalidCredentials,
		},
		{
			name:        "empty password",
			storedHash:  hash,
			rawPassword: "",
			want:        ldap.LDAPResultInvalidCredentials,
		},
		{
			name:        "user without password",
			storedHash:  "",
			rawPassword: "anything",
			want:        ldap.LDAPResultInvalidCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ldapHandler{
				backend: userBackend{
					users: map[string]model.UserGroup{
						"user": &model.User{
							Base:     model.Base{ID: 1, Name: "user"},
							Password: tt.storedHash,
						},
					},
				},
				logger: log.New("test"),
				config: &stns.Config{LDAP: &stns.LDAP{BaseDN: "dc=stns,dc=local"}},
			}

			client, server := net.Pipe()
			defer client.Close()
			defer server.Close()

			got, err := h.Bind("cn=user,dc=stns,dc=local", tt.rawPassword, client)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("unmatch result code want:%d got:%d", tt.want, got)
			}
		})
	}
}
