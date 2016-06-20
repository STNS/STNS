package stns

import (
	"net/http"

	"github.com/STNS/STNS/settings"
	"github.com/ant0ine/go-json-rest/rest"
)

type Handler struct {
	config *Config
}

type MetaData struct {
	ApiVersion float64 `json:"api_version"`
	Salt       bool    `json:"salt_enable"`
	Stretching int     `json:"stretching_number"`
	HashType   string  `json:"hash_type"`
	Result     string  `json:"result""`
	MinId      int     `json:"min_id""`
}

type ResponseFormat struct {
	MetaData *MetaData   `json:"metadata"`
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
	if r.URL.Path[1:3] == "v2" {
		if attr == nil {
			w.WriteHeader(http.StatusNotFound)
		}
		response := ResponseFormat{
			MetaData: &MetaData{
				ApiVersion: settings.API_VERSION,
				Salt:       h.config.Salt,
				Stretching: h.config.Stretching,
				Result:     settings.SUCCESS,
				HashType:   h.config.HashType,
				MinId:      q.GetMinId(),
			},
			Items: &attr,
		}
		w.WriteJson(response)
		return
	} else {
		if attr == nil {
			rest.NotFound(w, r)
		} else {
			w.WriteJson(attr)
		}
		return
	}
}
