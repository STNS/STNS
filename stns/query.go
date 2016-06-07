package stns

import "reflect"

type Query struct {
	config   *Config
	resource string
	column   string
	value    string
}

func (q *Query) getConfigByType() Attributes {
	if q.resource == "user" {
		return q.config.Users
	} else if q.resource == "group" {
		return q.config.Groups
	} else if q.resource == "sudo" {
		return q.config.Sudoers
	}
	return nil
}

func (q *Query) GetMinId() int {
	if q.resource == "user" {
		return MinUserId
	} else if q.resource == "group" {
		return MinGroupId
	}
	return 0
}

func (q *Query) getAttribute() Attributes {
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

func (q *Query) Get() Attributes {
	attr := q.getAttribute()
	if attr != nil && !reflect.ValueOf(attr).IsNil() {
		q.mergeLinkValue(attr)
	}
	return attr
}

func (q *Query) mergeLinkValue(attr Attributes) {
	for k, v := range attr {
		mergeValue := []string{}
		linker := q.getLinker(v)
		if linker != nil && !reflect.ValueOf(linker).IsNil() &&
			linker.LinkTargetValue() != nil && !reflect.ValueOf(linker.LinkTargetValue()).IsNil() {
			for _, linkValue := range linker.LinkTargetValue() {
				linkValues := map[string][]string{k: linker.LinkValue()}

				q.recursiveSetLinkValue(linkValue, linkValues)
				for _, val := range linkValues {
					mergeValue = append(mergeValue, val...)
				}
				linker.SetLinkValue(RemoveDuplicates(mergeValue))
			}
		}
	}
}

func (q *Query) getLinker(attr *Attribute) Linker {
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

	c := q.getConfigByType()

	if c != nil && !reflect.ValueOf(c).IsNil() {
		linker := q.getLinker(c.GetByName(name)[name])

		if linker != nil && !reflect.ValueOf(linker).IsNil() && len(linker.LinkValue()) > 0 {
			result[name] = linker.LinkValue()
			if linker.LinkTargetValue() != nil || !reflect.ValueOf(linker.LinkTargetValue()).IsNil() {
				for _, next_name := range linker.LinkTargetValue() {
					q.recursiveSetLinkValue(next_name, result)
				}
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
