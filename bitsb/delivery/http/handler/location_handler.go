package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/sainak/bitsb/api"
	"github.com/sainak/bitsb/bitsb"
	"github.com/sainak/bitsb/pkg/handler"
	"github.com/sainak/bitsb/pkg/repo"
)

type LocationHandler struct {
	service bitsb.LocationServiceProvider
}

func NewLocationHandler(service bitsb.LocationServiceProvider) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

func (l *LocationHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	cursor := r.URL.Query().Get("cursor")
	limit := handler.GetLimit(r)

	logrus.Print("query: ", query)

	filters := make(repo.Filters)
	if query != "" {
		filters["name:ilike"] = query
	}

	locations, nextCursor, err := l.service.ListAll(r.Context(), cursor, limit, filters)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	w.Header().Set("X-Cursor", nextCursor)
	render.JSON(w, r, locations)
}

func (l *LocationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	location, err := l.service.GetByID(r.Context(), id)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	render.JSON(w, r, location)
}

func (l *LocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &bitsb.LocationForm{}
	if err := render.Bind(r, data); err != nil {
		api.RespondForError(w, r, err)
		return
	}

	location := &bitsb.Location{
		Name: data.Name,
	}

	if err := l.service.Create(r.Context(), location); err != nil {
		api.RespondForError(w, r, err)
		return
	}

	render.JSON(w, r, location)
}

func (l *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	data := &bitsb.LocationForm{}
	if err = render.Bind(r, data); err != nil {
		api.RespondForError(w, r, err)
		return
	}

	location := &bitsb.Location{
		ID:   id,
		Name: data.Name,
	}

	if err = l.service.Update(r.Context(), location); err != nil {
		api.RespondForError(w, r, err)
		return
	}

	render.JSON(w, r, location)
}

func (l *LocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	if err = l.service.Delete(r.Context(), id); err != nil {
		api.RespondForError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
