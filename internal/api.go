package stns

import (
	"reflect"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

func GetAttr(w rest.ResponseWriter, r *rest.Request) {
	var attr *Attr
	var resource map[string]*Attr

	value := r.PathParam("value")
	column := r.PathParam("column")
	resource_name := r.PathParam("resource_name")

	if resource_name == "user" {
		resource = AllConfig.Users
	} else if resource_name == "group" {
		resource = AllConfig.Groups
	}

	if column == "id" {
		attr = _GetById(value, resource)
	} else if column == "name" {
		attr = _GetByName(value, resource)
	}

	if attr == nil || reflect.ValueOf(attr).IsNil() {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(attr)
}

func _GetByName(name string, resource map[string]*Attr) *Attr {
	attr := resource[name]
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		return nil
	}
	attr.Name = name
	return attr
}

func _GetById(_id string, resource map[string]*Attr) *Attr {
	id, _ := strconv.Atoi(_id)
	for k, u := range resource {
		if u.Id == id {
			u.Name = k
			return u
		}
	}
	return nil
}
