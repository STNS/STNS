package main

import (
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/pyama86/palk/internal"
)

func main() {
	palk.LoadConfig()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/:resource_name/:column/:value", GetAttr),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":1104", api.MakeHandler()))
}

func GetAttr(w rest.ResponseWriter, r *rest.Request) {
	var attr *palk.Attr
	var resource map[string]*palk.Attr

	value := r.PathParam("value")
	column := r.PathParam("column")
	resource_name := r.PathParam("resource_name")

	if resource_name == "user" {
		resource = palk.AllConfig.Users
	} else if resource_name == "group" {
		resource = palk.AllConfig.Groups
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

func _GetByName(name string, resource map[string]*palk.Attr) *palk.Attr {
	attr := resource[name]
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		return nil
	}
	attr.Name = name
	return attr
}

func _GetById(_id string, resource map[string]*palk.Attr) *palk.Attr {
	id, _ := strconv.Atoi(_id)
	for k, u := range resource {
		if u.Id == id {
			u.Name = k
			return u
		}
	}
	return nil
}
