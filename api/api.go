package api

import (
	"reflect"
	"strconv"

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
		attr = GetById(value, resource)
	} else if column == "name" {
		attr = GetByName(value, resource)
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

func GetByName(name string, resource attribute.UserGroups) attribute.UserGroups {
	attr := resource[name]
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		return nil
	}
	return attribute.UserGroups{
		name: attr,
	}
}

func GetById(_id string, resource attribute.UserGroups) attribute.UserGroups {
	id, _ := strconv.Atoi(_id)
	for k, u := range resource {
		if u.Id == id {
			return attribute.UserGroups{
				k: u,
			}
		}
	}
	return nil
}
