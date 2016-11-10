package stns

import (
	"net/http"
	"strconv"

	"github.com/STNS/STNS/settings"
	"github.com/ant0ine/go-json-rest/rest"
)

type responser interface {
	Response()
}

// ----------------------------------------
// v1
// ----------------------------------------
type v1_ResponseFormat struct {
	Items Attributes `json:"items"`
	query *Query
	w     rest.ResponseWriter
	r     *rest.Request
}

func (self *v1_ResponseFormat) Response() {
	if self.Items == nil {
		rest.NotFound(self.w, self.r)
	} else {
		self.w.WriteJson(self.Items)
	}
}

// ----------------------------------------
// v2
// ----------------------------------------
type v2_MetaData struct {
	ApiVersion float64 `json:"api_version"`
	Result     string  `json:"result"`
	MinId      int     `json:"min_id"`
}

type v2_ResponseFormat struct {
	MetaData *v2_MetaData `json:"metadata"`
	Items    Attributes   `json:"items"`
	query    *Query
	w        rest.ResponseWriter
	r        *rest.Request
}

func (self *v2_ResponseFormat) Response() {
	if self.Items == nil {
		self.w.WriteHeader(http.StatusNotFound)
	}

	response := v2_ResponseFormat{
		MetaData: &v2_MetaData{
			ApiVersion: settings.API_VERSION,
			Result:     settings.SUCCESS,
			MinId:      self.query.GetMinId(),
		},
		Items: self.Items,
	}
	self.w.WriteJson(response)
	return
}

// ----------------------------------------
// v3
// ----------------------------------------
type v3_ResponseFormat struct {
	Items Attributes `json:"items"`
	query *Query
	w     rest.ResponseWriter
	r     *rest.Request
}

type v3User struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Password  string   `json:"password"`
	GroupId   int      `json:"group_id"`
	Directory string   `json:"directory"`
	Shell     string   `json:"shell"`
	Gecos     string   `json:"gecos"`
	Keys      []string `json:"keys"`
}

type v3Group struct {
	Id    int      `json:"id"`
	Name  string   `json:"name"`
	Users []string `json:"users"`
}

type v3Sudo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type v3Users struct {
	items []*v3User
}
type v3Groups struct {
	items []*v3Group
}

type v3Sudoers struct {
	items []*v3Sudo
}

type v3Resource interface {
	buildResource(string, *Attribute) interface{}
}

func NewV3Resource(q *Query) v3Resource {
	switch q.resource {
	case "user":
		return v3Users{}
	case "group":
		return v3Groups{}
	case "sudo":
		return v3Sudoers{}
	}
	return nil
}

func (self v3Users) buildResource(n string, u *Attribute) interface{} {
	if u.User != nil {
		return &v3User{
			Name:      n,
			Id:        u.Id,
			Password:  u.Password,
			GroupId:   u.GroupId,
			Directory: u.Directory,
			Shell:     u.Shell,
			Gecos:     u.Gecos,
			Keys:      u.Keys,
		}
	}
	return nil
}

func (self v3Groups) buildResource(n string, g *Attribute) interface{} {
	if g.Group != nil {
		return &v3Group{
			Name:  n,
			Id:    g.Id,
			Users: g.Users,
		}
	}
	return nil
}

func (self v3Sudoers) buildResource(n string, u *Attribute) interface{} {
	if u.User != nil {
		return &v3Sudo{
			Name:     n,
			Password: u.Password,
		}
	}
	return nil
}

func (self *v3_ResponseFormat) Response() {
	if len(self.Items) == 0 {
		rest.NotFound(self.w, self.r)
		return
	}

	self.w.Header().Set("X-STNS-MIN-ID", strconv.Itoa(self.query.GetMinId()))

	resource := NewV3Resource(self.query)
	resources := []interface{}{}

	for n, u := range self.Items {
		resources = append(resources, resource.buildResource(n, u))
		if self.query.column != "list" {
			break
		}
	}

	if len(resources) > 0 {
		if self.query.column == "list" {
			self.w.WriteJson(resources)
		} else {
			self.w.WriteJson(resources[0])
		}
	} else {
		rest.NotFound(self.w, self.r)
	}
}

func NewResponder(q *Query, w rest.ResponseWriter, r *rest.Request) responser {
	res := q.Get()
	switch r.URL.Path[1:3] {
	case "v2":
		return &v2_ResponseFormat{
			Items: res,
			query: q,
			w:     w,
			r:     r,
		}
	case "v3":
		return &v3_ResponseFormat{
			Items: res,
			query: q,
			w:     w,
			r:     r,
		}
	default:
		return &v1_ResponseFormat{
			Items: res,
			query: q,
			w:     w,
			r:     r,
		}
	}
}
