package stns

import "github.com/ant0ine/go-json-rest/rest"

type Handler struct {
	config *Config
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
	res := NewResponder(q, w, r)
	res.Response()
}
