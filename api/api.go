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

func (q Query) Get(value string) attribute.UserGroups {
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
	return attr
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
	resource := query.Get()
	if resource == nil || reflect.ValueOf(resource).IsNil() {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(resource)
}

func HealthChech(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("success")
}
