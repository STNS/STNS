package model

import "fmt"

type Base struct {
	ID   int    `toml:"id" json:"id" yaml:"id"`
	Name string `toml:"name" json:"name" yaml:"name"`
}

func (b *Base) GetID() int {
	return b.ID
}

func (b *Base) GetName() string {
	return b.Name

}

func (b *Base) setName(n string) {
	b.Name = n
}

type NotFoundError struct {
	resource string
	id       int
	name     string
}

func NewNotFoundError(r string, v interface{}) NotFoundError {
	e := NotFoundError{resource: r}
	switch v.(type) {
	case int:
		e.id = v.(int)
	case string:
		e.name = v.(string)
	}

	return e
}

func (e NotFoundError) Error() string {
	if e.id != 0 {
		return fmt.Sprintf("%s id %d is not found", e.resource, e.id)
	}
	if e.name != "" {
		return fmt.Sprintf("%s name %s is not found", e.resource, e.name)
	}
	return fmt.Sprintf("%s name not found", e.resource)
}

func errorHandler(res map[string]UserGroup, err error, value interface{}, rtype string) (map[string]UserGroup, error) {
	if err != nil {
		return nil, err
	}

	if res == nil || len(res) == 0 {
		return nil, NewNotFoundError(rtype, value)
	}
	return res, nil
}
