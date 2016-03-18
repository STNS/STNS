package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	"github.com/STNS/STNS/config"
	"github.com/ant0ine/go-json-rest/rest/test"
)

func TestHandler(t *testing.T) {
	createHandlerTestConfig()

	recorded := test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{
  "example": {
    "id": 1000,
    "group_id": 2000,
    "directory": "/home/example",
    "shell": "/bin/bash",
    "gecos": "",
    "keys": [
      "ssh-rsa aaa"
    ],
    "link_users": null,
    "users": null
  }
}`)
	recorded = test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/id/1000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{
  "example": {
    "id": 1000,
    "group_id": 2000,
    "directory": "/home/example",
    "shell": "/bin/bash",
    "gecos": "",
    "keys": [
      "ssh-rsa aaa"
    ],
    "link_users": null,
    "users": null
  }
}`)
	recorded = test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/id/1001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/name/example_group", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{
  "example_group": {
    "id": 3000,
    "group_id": 0,
    "directory": "",
    "shell": "",
    "gecos": "",
    "keys": null,
    "link_users": null,
    "users": [
      "example"
    ]
  }
}`)
	recorded = test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/name/example_group3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/id/3000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{
  "example_group": {
    "id": 3000,
    "group_id": 0,
    "directory": "",
    "shell": "",
    "gecos": "",
    "keys": null,
    "link_users": null,
    "users": [
      "example"
    ]
  }
}`)

	recorded = test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/id/3001", nil))
	recorded.CodeIs(404)
}

func TestBasicAuth(t *testing.T) {
	createAuthTestConfig()
	// simple request fails
	recorded := test.RunRequest(t, getHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil))
	recorded.CodeIs(401)

	// auth with wrong cred and right method fails
	wrongCredReq := test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil)
	encoded := base64.StdEncoding.EncodeToString([]byte("admin:AdmIn"))
	wrongCredReq.Header.Set("Authorization", "Basic "+encoded)
	recorded = test.RunRequest(t, getHandler(), wrongCredReq)
	recorded.CodeIs(401)

	rightCredReq := test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil)
	encoded = base64.StdEncoding.EncodeToString([]byte("admin:Admin"))
	rightCredReq.Header.Set("Authorization", "Basic "+encoded)
	recorded = test.RunRequest(t, getHandler(), rightCredReq)
	recorded.CodeIs(200)

}

func createHandlerTestConfig() {
	configFile, _ := ioutil.TempFile("", "stns-config-handler")
	configContent := `port = 9999
[users.example]
id = 1000
group_id = 2000
directory = "/home/example"
shell = "/bin/bash"
keys = ["ssh-rsa aaa"]

[groups.example_group]
id = 3000
users = ["example"]
`
	_, _ = configFile.WriteString(configContent)
	configFile.Close()
	name := configFile.Name()
	defer os.Remove(name)
	config.Load(&name)
}

func createAuthTestConfig() {
	configFile, _ := ioutil.TempFile("", "stns-config-auth")
	configContent := `port = 9999
user = "admin"
password = "Admin"
[users.example]
id = 1000
`
	_, _ = configFile.WriteString(configContent)
	configFile.Close()
	name := configFile.Name()
	defer os.Remove(name)
	config.Load(&name)
}
