package stns

import "github.com/ant0ine/go-json-rest/rest"

// Handler handler object
type Handler struct {
	config *Config
}

func (h *Handler) getQuery(r *rest.Request) *Query {
	value := r.PathParam("value")
	column := r.PathParam("column")
	resourceName := r.PathParam("resource_name")
	return &Query{h.config, resourceName, column, value}
}

func (h *Handler) getListQuery(r *rest.Request) *Query {
	resourceName := r.PathParam("resource_name")
	return &Query{h.config, resourceName, "list", ""}
}

// Get get a resource response
func (h *Handler) Get(w rest.ResponseWriter, r *rest.Request) {
	query := h.getQuery(r)
	h.Response(query, w, r)
}

// GetList get resource list response
func (h *Handler) GetList(w rest.ResponseWriter, r *rest.Request) {
	query := h.getListQuery(r)
	h.Response(query, w, r)
}

// Response proxy to the reponsder
func (h *Handler) Response(q *Query, w rest.ResponseWriter, r *rest.Request) {
	res := newResponder(q, w, r)
	res.Response()
}
