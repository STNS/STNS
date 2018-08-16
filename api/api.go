package api

import "github.com/STNS/STNS/model"

func toSlice(ug map[string]model.UserGroup) []model.UserGroup {
	if ug == nil {
		return nil
	}
	r := []model.UserGroup{}
	for _, v := range ug {
		r = append(r, v)
	}
	return r
}
