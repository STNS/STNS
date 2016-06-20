package stns

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

func (q *Query) Get() Attributes {
	resource := q.getConfigByType()
	if resource != nil {
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

func (q *Query) GetMinId() int {
	if q.resource == "user" {
		return MinUserId
	} else if q.resource == "group" {
		return MinGroupId
	}
	return 0
}
