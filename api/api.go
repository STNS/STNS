package api

import (
	"reflect"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/STNS/config"
)

func Get(w rest.ResponseWriter, r *rest.Request) {
	var attr attribute.UserGroups
	var resource attribute.UserGroups

	value := r.PathParam("value")
	column := r.PathParam("column")
	resource_name := r.PathParam("resource_name")

	if resource_name == "user" {
		resource = config.All.Users
	} else if resource_name == "group" {
		resource = config.All.Groups
	}

	if column == "id" {
		attr = resource.GetById(value)
	} else if column == "name" {
		attr = resource.GetByName(value)
	}

	if attr == nil || reflect.ValueOf(attr).IsNil() {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(attr)
}
func GetList(w rest.ResponseWriter, r *rest.Request) {
	var resource attribute.UserGroups
	resource_name := r.PathParam("resource_name")

	if resource_name == "user" {
		resource = config.All.Users
	} else if resource_name == "group" {
		resource = config.All.Groups
	}

	if resource == nil || reflect.ValueOf(resource).IsNil() {
		rest.NotFound(w, r)
		return
	}

	w.WriteJson(resource)
}
