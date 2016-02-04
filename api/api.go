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
		attr = resource.GetById(value)
	} else if q.column == "name" {
		attr = resource.GetByName(value)
	} else if q.column == "list" {
		attr = resource
	}
	return attr
}

func Get(w rest.ResponseWriter, r *rest.Request) {
	value := r.PathParam("value")
	column := r.PathParam("column")
	resource_name := r.PathParam("resource_name")
	query := Query{resource_name, column}

	attr := query.Get(value)
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(attr)
}
func GetList(w rest.ResponseWriter, r *rest.Request) {
	resource_name := r.PathParam("resource_name")

	query := Query{resource_name, "list"}
	resource := query.Get("")

	if resource == nil || reflect.ValueOf(resource).IsNil() {
		rest.NotFound(w, r)
		return
	}

	w.WriteJson(resource)
}

func HealthChech(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("success")
}
