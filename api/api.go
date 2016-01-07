package api

import (
	"reflect"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/STNS/config"
)

func GetAttribute(w rest.ResponseWriter, r *rest.Request) {
	var attr *attribute.All
	var resource map[string]*attribute.All

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

func GetByName(name string, resource map[string]*attribute.All) *attribute.All {
	attr := resource[name]
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		return nil
	}
	attr.Name = name
	return attr
}

func GetById(_id string, resource map[string]*attribute.All) *attribute.All {
	id, _ := strconv.Atoi(_id)
	for k, u := range resource {
		if u.Id == id {
			u.Name = k
			return u
		}
	}
	return nil
}
