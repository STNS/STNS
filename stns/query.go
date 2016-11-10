package stns

// Query query object
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

// Get get resouce from endpoint(id,name,list)
func (q *Query) Get() Attributes {
	resource := q.getConfigByType()
	if resource != nil {
		if q.column == "id" {
			return resource.GetByID(q.value)
		} else if q.column == "name" {
			return resource.GetByName(q.value)
		} else if q.column == "list" {
			return resource
		}
	}
	return nil
}

// GetMinID get the minimum id of the specified resource
func (q *Query) GetMinID() int {
	if q.resource == "user" {
		return minUserID
	} else if q.resource == "group" {
		return minGroupID
	}
	return 0
}
