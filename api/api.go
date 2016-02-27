package api

import (
	"reflect"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/STNS/config"
	"github.com/ant0ine/go-json-rest/rest"
)

type Query struct {
	resource string
	column   string
	value    string
}

func (q Query) Get() attribute.UserGroups {
	var attr attribute.UserGroups
	var resource attribute.UserGroups

	if q.resource == "user" {
		resource = config.All.Users
	} else if q.resource == "group" {
		resource = config.All.Groups
	}
	if q.column == "id" {
		attr = resource.GetById(q.value)
	} else if q.column == "name" {
		attr = resource.GetByName(q.value)
	} else if q.column == "list" {
		attr = resource
	}

	// merge user keys
	q.mergeKey(attr)
	return attr
}
func (q Query) mergeKey(attr attribute.UserGroups) {
	// merge user keys
	for k, u := range attr {
		merge_keys := []string{}
		if u.LinkUsers != nil || !reflect.ValueOf(u.LinkUsers).IsNil() {
			for _, link_users_name := range u.LinkUsers {
				link_keys := map[string][]string{k: u.Keys}

				q.recursiveSetLinkKey(link_users_name, link_keys)
				for _, user_keys := range link_keys {
					merge_keys = append(merge_keys, user_keys...)
				}

				u.Keys = RemoveDuplicates(merge_keys)
			}
		}
	}
}

// ref:http://qiita.com/yukitomo/items/2e6be0f26905d8e3dd22
func member(n string, xs []string) bool {
	for _, x := range xs {
		if n == x {
			return true
		}
	}
	return false
}

func RemoveDuplicates(xs []string) []string {
	ys := make([]string, 0, len(xs))
	for _, x := range xs {
		if !member(x, ys) {
			ys = append(ys, x)
		}
	}
	return ys
}

func (q Query) recursiveSetLinkKey(name string, result map[string][]string) {
	if result[name] != nil {
		return
	}

	user := config.All.Users.GetByName(name)
	if user != nil && len(user[name].Keys) > 0 {
		result[name] = user[name].Keys
		if user[name].LinkUsers != nil || !reflect.ValueOf(user[name].LinkUsers).IsNil() {
			for _, nest_user_name := range user[name].LinkUsers {
				q.recursiveSetLinkKey(nest_user_name, result)
			}
		}
	}
}

func Get(w rest.ResponseWriter, r *rest.Request) {
	value := r.PathParam("value")
	column := r.PathParam("column")
	resource_name := r.PathParam("resource_name")
	query := Query{resource_name, column, value}
	query.Response(w, r)
}
func GetList(w rest.ResponseWriter, r *rest.Request) {
	resource_name := r.PathParam("resource_name")
	query := Query{resource_name, "list", ""}
	query.Response(w, r)
}

func (q Query) Response(w rest.ResponseWriter, r *rest.Request) {
	attr := q.Get()
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		rest.NotFound(w, r)
		return
	}

	w.WriteJson(attr)
}

func HealthChech(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("success")
}
