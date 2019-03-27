package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	httpIntegration = flag.Bool("integration-http", false, "run http integration tests")
)

// testHost is host url for test
const testHost = "http://localhost:1104"

// testEndpoint is endpoint for API
const testEndpoint = testHost + "/v1"

// testMime is request content-type
const testMime = "application/json"

func TestMain(m *testing.M) {
	flag.Parse()
	result := m.Run()
	os.Exit(result)
}

func TestHTTPGetUserList(t *testing.T) {
	if !*httpIntegration {
		t.Skip()
	}

	res, err := http.Get(testEndpoint + "/users")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("GetUserList API returned wrong status code: got %v want %v",
			res.StatusCode, http.StatusOK)
	}

	var users interface{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		t.Fatal(err)
	}

	expectedCount := 203
	gotCount := len(users.([]interface{}))
	if gotCount != expectedCount {
		t.Errorf("GetUsers API returned wrong count: got %v expected %v",
			gotCount, expectedCount)
	}

	expectedHighestID := "10002"
	gotHighestID := res.Header.Get("User-Highest-Id")
	if gotHighestID != expectedHighestID {
		t.Errorf("GetUsers API returned wrong highest id: got %v expected %v",
			gotHighestID, expectedHighestID)
	}

	expectedLowestID := "99"
	gotLowestID := res.Header.Get("User-Lowest-Id")
	if gotLowestID != expectedLowestID {
		t.Errorf("GetUsers API returned wrong lowest id: got %v expected %v",
			gotLowestID, expectedLowestID)
	}
}

func TestHTTPGetUserByName(t *testing.T) {
	if !*httpIntegration {
		t.Skip()
	}

	res, err := http.Get(testEndpoint + "/users?name=foo")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("GetUserList API returned wrong status code: got %v want %v",
			res.StatusCode, http.StatusOK)
	}

	var users interface{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		t.Fatal(err)
	}

	fst := users.([]interface{})[0].(map[string]interface{})

	expectedName := "foo"
	gotName, ok := fst["name"].(string)
	if ok && gotName != expectedName {
		t.Errorf("GetUsers API returned wrong user: got %s expected %s",
			gotName, expectedName)
	}

	expectedCount := 1
	gotCount := len(users.([]interface{}))
	if gotCount != expectedCount {
		t.Errorf("GetUsers API returned wrong count: got %v expected %v",
			gotCount, expectedCount)
	}
}

func TestHTTPGetUserByID(t *testing.T) {
	if !*httpIntegration {
		t.Skip()
	}

	res, err := http.Get(testEndpoint + "/users?id=10001")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("GetUserList API returned wrong status code: got %v want %v",
			res.StatusCode, http.StatusOK)
	}

	var users interface{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		t.Fatal(err)
	}

	fst := users.([]interface{})[0].(map[string]interface{})

	expectedName := "test"
	gotName, ok := fst["name"].(string)
	if ok && gotName != expectedName {
		t.Errorf("GetUsers  API returned wrong user: got %s expected %s",
			gotName, expectedName)
	}

	expectedCount := 1
	gotCount := len(users.([]interface{}))
	if gotCount != expectedCount {
		t.Errorf("GetUsers API returned wrong count: got %v expected %v",
			gotCount, expectedCount)
	}
}

func TestHTTPGetGroupList(t *testing.T) {
	if !*httpIntegration {
		t.Skip()
	}

	res, err := http.Get(testEndpoint + "/groups")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("GetGroupList API returned wrong status code: got %v want %v",
			res.StatusCode, http.StatusOK)
	}

	var groups interface{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		t.Fatal(err)
	}

	expectedCount := 3
	gotCount := len(groups.([]interface{}))
	if gotCount != expectedCount {
		t.Errorf("GetGroups API returned wrong count: got %v expected %v",
			gotCount, expectedCount)
	}

	expectedHighestID := "10002"
	gotHighestID := res.Header.Get("Group-Highest-Id")
	if gotHighestID != expectedHighestID {
		t.Errorf("GetGroups API returned wrong highest id: got %v expected %v",
			gotHighestID, expectedHighestID)
	}

	expectedLowestID := "100"
	gotLowestID := res.Header.Get("Group-Lowest-Id")
	if gotLowestID != expectedLowestID {
		t.Errorf("GetGroups API returned wrong lowest id: got %v expected %v",
			gotLowestID, expectedLowestID)
	}
}

func TestHTTPGetGroupByName(t *testing.T) {
	if !*httpIntegration {
		t.Skip()
	}

	res, err := http.Get(testEndpoint + "/groups?name=bar")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("GetGroupList API returned wrong status code: got %v want %v",
			res.StatusCode, http.StatusOK)
	}

	var groups interface{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		t.Fatal(err)
	}

	fst := groups.([]interface{})[0].(map[string]interface{})

	expectedName := "bar"
	gotName, ok := fst["name"].(string)
	if ok && gotName != expectedName {
		t.Errorf("GetGroups API returned wrong group: got %s expected %s",
			gotName, expectedName)
	}

	expectedCount := 1
	gotCount := len(groups.([]interface{}))
	if gotCount != expectedCount {
		t.Errorf("GetGroups API returned wrong count: got %v expected %v",
			gotCount, expectedCount)
	}
}

func TestHTTPGetGroupByID(t *testing.T) {
	if !*httpIntegration {
		t.Skip()
	}

	res, err := http.Get(testEndpoint + "/groups?id=10001")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("GetGroupList API returned wrong status code: got %v want %v",
			res.StatusCode, http.StatusOK)
	}

	var groups interface{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		t.Fatal(err)
	}

	fst := groups.([]interface{})[0].(map[string]interface{})

	expectedName := "test"
	gotName, ok := fst["name"].(string)
	if ok && gotName != expectedName {
		t.Errorf("GetGroups  API returned wrong group: got %s expected %s",
			gotName, expectedName)
	}

	expectedCount := 1
	gotCount := len(groups.([]interface{}))
	if gotCount != expectedCount {
		t.Errorf("GetGroups API returned wrong count: got %v expected %v",
			gotCount, expectedCount)
	}
}

func TestHTTPRoot(t *testing.T) {
	if !*httpIntegration {
		t.Skip()
	}

	res, err := http.Get(testHost + "/")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("STNS Status API returned wrong status code: got %v want %v", res.StatusCode, http.StatusOK)
	}

	expected := "Hello! STNS!!1"
	got := string(body)
	if got != expected {
		t.Errorf("STNS Status API returned wrong body: got %v expected %v", got, expected)
	}
}
