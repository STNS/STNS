package api

import (
	"reflect"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/STNS/config"
	"github.com/ant0ine/go-json-rest/rest"
)

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

type Query struct {
	resource string
	column   string
	value    string
}

func (q *Query) getConfigByType() attribute.UserGroups {
	if q.resource == "user" {
		return config.All.Users
	} else if q.resource == "group" {
		return config.All.Groups
	}
	return nil
}

func (q *Query) getAttribute() attribute.UserGroups {
	resource := q.getConfigByType()
	if resource != nil && !reflect.ValueOf(resource).IsNil() {

		if q.column == "id" {
			return resource.GetById(q.value)
		} else if q.column == "name" {
			return resource.GetByName(q.value)
		} else if q.column == "list" {
			return resource
		}
	}
	return nil
}

func (q *Query) Get() attribute.UserGroups {
	attr := q.getAttribute()
	if attr != nil && !reflect.ValueOf(attr).IsNil() {
		q.mergeLinkValue(attr)
	}
	return attr
}

func (q *Query) mergeLinkValue(attr attribute.UserGroups) {

	for k, v := range attr {
		linker := q.getLinker(v)
		mergeValue := []string{}
		if linker != nil && !reflect.ValueOf(linker).IsNil() &&
			linker.LinkTargetColumnValue() != nil && !reflect.ValueOf(linker.LinkTargetColumnValue()).IsNil() {
			for _, linkValue := range linker.LinkTargetColumnValue() {
				linkValues := map[string][]string{k: linker.LinkValues()}

				q.recursiveSetLinkValue(linkValue, linkValues)
				for _, val := range linkValues {
					mergeValue = append(mergeValue, val...)
				}
				linker.SetLinkValue(RemoveDuplicates(mergeValue))
			}
		}
	}
}

func (q *Query) getLinker(attr *attribute.All) attribute.Linker {
	if attr != nil && !reflect.ValueOf(attr).IsNil() {
		if q.resource == "user" {
			return attr.User
		} else if q.resource == "group" {
			return attr.Group
		}
	}
	return nil
}

func (q *Query) recursiveSetLinkValue(name string, result map[string][]string) {
	if result[name] != nil {
		return
	}

	config := q.getConfigByType()
	linker := q.getLinker(config.GetByName(name)[name])

	if linker != nil && !reflect.ValueOf(linker).IsNil() && len(linker.LinkValues()) > 0 {
		result[name] = linker.LinkValues()
		if linker.LinkTargetColumnValue() != nil || !reflect.ValueOf(linker.LinkTargetColumnValue()).IsNil() {
			for _, next_name := range linker.LinkTargetColumnValue() {
				q.recursiveSetLinkValue(next_name, result)
			}
		}
	}
}

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

func (q *Query) Response(w rest.ResponseWriter, r *rest.Request) {
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
