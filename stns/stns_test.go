package stns

import (
	"crypto/tls"
	"encoding/base64"
	"testing"

	stns_test "github.com/STNS/STNS/test"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
)

func TestHandlerV1User(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example":{"id":1000,"password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null,"setup_commands":null}}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/id/1000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example":{"id":1000,"password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null,"setup_commands":null}}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/id/1001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example":{"id":1000,"password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null,"setup_commands":null}}`)
}

func TestHandlerV1Group(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/name/example_group", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example_group":{"id":3000,"users":["example"],"link_groups":null}}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/name/example_group3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/id/3000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example_group":{"id":3000,"users":["example"],"link_groups":null}}`)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/id/3001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example_group":{"id":3000,"users":["example"],"link_groups":null}}`)
}

func TestHandlerV1Sudo(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/sudo/name/example_sudo", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example_sudo":{"id":0,"password":"p@ssword","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"link_users":null,"setup_commands":null}}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/sudo/name/example_notfound", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/sudo/id/1001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/sudo/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example_sudo":{"id":0,"password":"p@ssword","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"link_users":null,"setup_commands":null}}`)
}

func TestHandlerV2User(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/name/example", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2.1,"result":"success"},"items":{"example":{"id":1000,"password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null,"setup_commands":null}}}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/name/example3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/id/1000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2.1,"result":"success"},"items":{"example":{"id":1000,"password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null,"setup_commands":null}}}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/id/1001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2.1,"result":"success"},"items":{"example":{"id":1000,"password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null,"setup_commands":null}}}`)
}

func TestHandlerV2Group(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/name/example_group", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2.1,"result":"success"},"items":{"example_group":{"id":3000,"users":["example"],"link_groups":null}}}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/name/example_group3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/id/3000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2.1,"result":"success"},"items":{"example_group":{"id":3000,"users":["example"],"link_groups":null}}}`)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/id/3001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2.1,"result":"success"},"items":{"example_group":{"id":3000,"users":["example"],"link_groups":null}}}`)
}

func TestHandlerV2Sudo(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/sudo/name/example_sudo", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2.1,"result":"success"},"items":{"example_sudo":{"id":0,"password":"p@ssword","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"link_users":null,"setup_commands":null}}}`)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/sudo/name/example_notfound", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/sudo/id/1001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/sudo/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2.1,"result":"success"},"items":{"example_sudo":{"id":0,"password":"p@ssword","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"link_users":null,"setup_commands":null}}}`)
}

func TestHandlerV3User(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/user/name/example", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"id":1000,"prev_id":0,"next_id":0,"name":"example","password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"setup_commands":null}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/user/name/example3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/user/id/1000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"id":1000,"prev_id":0,"next_id":0,"name":"example","password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"setup_commands":null}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/user/id/1001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/user/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`[{"id":1000,"prev_id":0,"next_id":0,"name":"example","password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"setup_commands":null}]`)
}

func TestHandlerv3Group(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/group/name/example_group", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"id":3000,"prev_id":0,"next_id":0,"name":"example_group","users":["example"]}`)
	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/group/name/example_group3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/group/id/3000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"id":3000,"prev_id":0,"next_id":0,"name":"example_group","users":["example"]}`)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/group/id/3001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/group/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`[{"id":3000,"prev_id":0,"next_id":0,"name":"example_group","users":["example"]}]`)
}

func TestHandlerv3Sudo(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/sudo/name/example_sudo", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"name":"example_sudo","password":"p@ssword"}`)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/sudo/name/example_notfound", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/sudo/id/1001", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/sudo/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`[{"name":"example_sudo","password":"p@ssword"}]`)
}

func TestBasicAuth(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_02.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	// simple request fails
	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil))
	recorded.CodeIs(401)

	// auth with wrong cred and right method fails
	wrongCredReq := test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil)
	encoded := base64.StdEncoding.EncodeToString([]byte("admin:AdmIn"))
	wrongCredReq.Header.Set("Authorization", "Basic "+encoded)
	recorded = test.RunRequest(t, s.newAPIHandler(), wrongCredReq)
	recorded.CodeIs(401)

	rightCredReq := test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil)
	encoded = base64.StdEncoding.EncodeToString([]byte("admin:Admin"))
	rightCredReq.Header.Set("Authorization", "Basic "+encoded)
	recorded = test.RunRequest(t, s.newAPIHandler(), rightCredReq)
	recorded.CodeIs(200)
}
func TestNewHTTPServerTLSAuth(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_03.conf")
	s := NewServer(config, "", "", false)
	h := s.newHTTPServer()
	stns_test.Assert(t, h.TLSConfig.ClientAuth == tls.RequireAndVerifyClientCert, "unmatch auth tls config")
}

func TestNewHTTPServerTLSNonAuth(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_02.conf")
	s := NewServer(config, "", "", false)
	h := s.newHTTPServer()
	stns_test.Assert(t, h.TLSConfig == nil, "unmatch non auth tls config")
}

func TestTLSKeysNotExists(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_02.conf")
	s := NewServer(config, "", "", false)
	stns_test.Assert(t, s.tlsKeysExists() == false, "unmatch tls keys not exists")
}

func TestTLSKeysExists(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_03.conf")
	s := NewServer(config, "", "", false)
	stns_test.Assert(t, s.tlsKeysExists() == true, "unmatch tls keys not exists")
}

func TestHandlerPeriodInUsername(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_04.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/user/name/example%2e1", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"id":1000,"prev_id":0,"next_id":0,"name":"example.1","password":"p@ssword","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"setup_commands":null}`)
}

func TestHandlerPrevNextID(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_05.conf")
	s := NewServer(config, "", "", false)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/user/list", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`[{"id":1001,"prev_id":0,"next_id":1002,"name":"example1","password":"","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"setup_commands":null},{"id":1002,"prev_id":1001,"next_id":1003,"name":"example2","password":"","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"setup_commands":null},{"id":1003,"prev_id":1002,"next_id":0,"name":"example3","password":"","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"setup_commands":null}]`)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/user/id/1001", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"id":1001,"prev_id":0,"next_id":1002,"name":"example1","password":"","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"setup_commands":null}`)

	recorded = test.RunRequest(t, s.newAPIHandler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v3/group/id/1001", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"id":1001,"prev_id":0,"next_id":1002,"name":"example1","users":["example1"]}`)
}
