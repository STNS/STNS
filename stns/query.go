package stns

import "reflect"

type Query struct {
	config   *Config
	backend  Backend
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

func (q *Query) Get() Attributes {
	attr := q.getConfigByType()
	if attr != nil {
		if q.column == "id" {
			return attr.GetById(q.value)
		} else if q.column == "name" {
			return q.mergeBackendPassword(
				attr.GetByName(q.value),
				q.value,
			)
		} else if q.column == "list" {
			return attr
		}
	}
	return nil
}

func (q *Query) mergeBackendPassword(attr Attributes, name string) Attributes {
	if q.resource == "user" && q.backend != nil && !reflect.ValueOf(q.backend).IsNil() {
		bu := q.backend.UserFindByName(name)
		if bu != nil {
			// this is pointer so over write value
			attr[name].Password = bu.Password
		}
	}
	return attr
}

func (q *Query) GetMinId() int {
	if q.resource == "user" {
		return MinUserId
	} else if q.resource == "group" {
		return MinGroupId
	}
	return 0
}
