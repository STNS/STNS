package stns

import (
	"reflect"
	"regexp"

	"github.com/STNS/STNS/settings"
	"github.com/ant0ine/go-json-rest/rest"
)

type Handler struct {
	config *Config
}

type Metadata struct {
	ApiVersion float32 `json:"api_version"`
	Salt       bool    `json:"salt_enable"`
	Stretching int     `json:"stretching_count"`
}

type ResponseFormat struct {
	Metadata *Metadata   `json:"metadata"`
	Items    *Attributes `json:"items"`
}

func (h *Handler) getQuery(r *rest.Request) *Query {
	value := r.PathParam("value")
	column := r.PathParam("column")
	resource_name := r.PathParam("resource_name")
	return &Query{h.config, resource_name, column, value}
}

func (h *Handler) getListQuery(r *rest.Request) *Query {
	resource_name := r.PathParam("resource_name")
	return &Query{h.config, resource_name, "list", ""}
}

func (h *Handler) Get(w rest.ResponseWriter, r *rest.Request) {
	query := h.getQuery(r)
	h.Response(query, w, r)
}

func (h *Handler) GetList(w rest.ResponseWriter, r *rest.Request) {
	query := h.getListQuery(r)
	h.Response(query, w, r)
}

func (h *Handler) Response(q *Query, w rest.ResponseWriter, r *rest.Request) {
	attr := q.Get()
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		rest.NotFound(w, r)
		return
	}

	v2 := regexp.MustCompile(`^/v2`)
	if v2.MatchString(r.URL.Path) {
		response := ResponseFormat{
			&Metadata{
				settings.ApiVersion,
				h.config.Salt,
				h.config.Stretching,
			},
			&attr,
		}
		w.WriteJson(response)
	} else {
		w.WriteJson(attr)
	}
}
