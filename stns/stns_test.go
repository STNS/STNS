package stns

import (
	"encoding/base64"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
)

func TestHandlerV1User(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := Create(config, "", "", false, nil)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example":{"id":1000,"password":"p@ssword","hash_type":"sha256","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null}}`)
	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/id/1000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example":{"id":1000,"password":"p@ssword","hash_type":"sha256","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null}}`)
	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/id/1001", nil))
	recorded.CodeIs(404)
}

func TestHandlerV1Group(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := Create(config, "", "", false, nil)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/name/example_group", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example_group":{"id":3000,"users":["example"],"link_groups":null}}`)
	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/name/example_group3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/id/3000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example_group":{"id":3000,"users":["example"],"link_groups":null}}`)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/group/id/3001", nil))
	recorded.CodeIs(404)
}

func TestHandlerV1Sudo(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := Create(config, "", "", false, nil)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/sudo/name/example_sudo", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"example_sudo":{"id":0,"password":"p@ssword","hash_type":"sha512","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"link_users":null}}`)
	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/sudo/name/example_notfound", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/sudo/id/1001", nil))
	recorded.CodeIs(404)
}

func TestHandlerV2User(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := Create(config, "", "", false, nil)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/name/example", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2,"salt_enable":false,"stretching_number":0,"hash_type":"sha256","result":"success","min_id":1000},"items":{"example":{"id":1000,"password":"p@ssword","hash_type":"sha256","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null}}}`)
	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/name/example3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/id/1000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2,"salt_enable":false,"stretching_number":0,"hash_type":"sha256","result":"success","min_id":1000},"items":{"example":{"id":1000,"password":"p@ssword","hash_type":"sha256","group_id":2000,"directory":"/home/example","shell":"/bin/bash","gecos":"","keys":["ssh-rsa aaa"],"link_users":null}}}`)
	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/user/id/1001", nil))
	recorded.CodeIs(404)
}
func TestHandlerV2Group(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := Create(config, "", "", false, nil)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/name/example_group", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2,"salt_enable":false,"stretching_number":0,"hash_type":"sha256","result":"success","min_id":3000},"items":{"example_group":{"id":3000,"users":["example"],"link_groups":null}}}`)
	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/name/example_group3", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/id/3000", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2,"salt_enable":false,"stretching_number":0,"hash_type":"sha256","result":"success","min_id":3000},"items":{"example_group":{"id":3000,"users":["example"],"link_groups":null}}}`)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/group/id/3001", nil))
	recorded.CodeIs(404)
}
func TestHandlerV2Sudo(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_01.conf")
	s := Create(config, "", "", false, nil)
	s.SetMiddleWare(rest.DefaultCommonStack)

	recorded := test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/sudo/name/example_sudo", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs(`{"metadata":{"api_version":2,"salt_enable":false,"stretching_number":0,"hash_type":"sha256","result":"success","min_id":0},"items":{"example_sudo":{"id":0,"password":"p@ssword","hash_type":"sha512","group_id":0,"directory":"","shell":"","gecos":"","keys":null,"link_users":null}}}`)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/sudo/name/example_notfound", nil))
	recorded.CodeIs(404)

	recorded = test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/v2/sudo/id/1001", nil))
	recorded.CodeIs(404)

}
func TestBasicAuth(t *testing.T) {
	config, _ := LoadConfig("./fixtures/stns_02.conf")
	s := Create(config, "", "", false, nil)
	s.SetMiddleWare(rest.DefaultCommonStack)

	// simple request fails
	recorded := test.RunRequest(t, s.Handler(), test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil))
	recorded.CodeIs(401)

	// auth with wrong cred and right method fails
	wrongCredReq := test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil)
	encoded := base64.StdEncoding.EncodeToString([]byte("admin:AdmIn"))
	wrongCredReq.Header.Set("Authorization", "Basic "+encoded)
	recorded = test.RunRequest(t, s.Handler(), wrongCredReq)
	recorded.CodeIs(401)

	rightCredReq := test.MakeSimpleRequest("GET", "http://localhost:9999/user/name/example", nil)
	encoded = base64.StdEncoding.EncodeToString([]byte("admin:Admin"))
	rightCredReq.Header.Set("Authorization", "Basic "+encoded)
	recorded = test.RunRequest(t, s.Handler(), rightCredReq)
	recorded.CodeIs(200)

}
