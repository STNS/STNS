package model

import (
	"fmt"
)

type UserGroup interface {
	GetID() int
	GetName() string
	setName(string)
	linkValues() []string
	setLinkValues([]string)
	value() []string
}

type Backend interface {
	GetterBackend
	SetterBackend
}

type GetterBackends []GetterBackend
type GetterBackend interface {
	FindUserByID(int) (map[string]UserGroup, error)
	FindUserByName(string) (map[string]UserGroup, error)
	FindGroupByID(int) (map[string]UserGroup, error)
	FindGroupByName(string) (map[string]UserGroup, error)
	Users() (map[string]UserGroup, error)
	Groups() (map[string]UserGroup, error)
	HighestUserID() int
	LowestUserID() int
	HighestGroupID() int
	LowestGroupID() int
}

type SetterBackend interface {
	Create(string, UserGroup) error
	Delete(string) error
	Update(string, UserGroup) error
}

func mergeUserGroup(m1, m2 map[string]UserGroup) map[string]UserGroup {
	ans := map[string]UserGroup{}

	for k, v := range m1 {
		ans[k] = v
	}
	for k, v := range m2 {
		ans[k] = v
	}
	return (ans)
}

func SyncConfig(resourceName string, b Backend, configResources, backendResources map[string]UserGroup) error {
	if configResources != nil {
		for _, cu := range configResources {
			found := false

			if backendResources != nil {
				for _, eu := range backendResources {
					if cu.GetID() == eu.GetID() {
						found = true
						break
					}
				}
			}

			if !found {
				if err := b.Create(fmt.Sprintf("/%s/name/%s", resourceName, cu.GetName()), cu); err != nil {
					return err
				}
			}
		}
	}

	if backendResources != nil {
		for _, eu := range backendResources {
			found := false
			if configResources != nil {
				for _, cu := range configResources {
					if cu.GetID() == eu.GetID() {
						found = true
						break
					}
				}
			}
			if !found {
				if err := b.Delete(fmt.Sprintf("/%s/name/%s", resourceName, eu.GetName())); err != nil {
					return err
				}
			}
		}

	}
	return nil
}
