package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/utils/repo"
	"github.com/sainak/bitsb/utils/response"
)

type LocationHandler struct {
	service domain.LocationServiceProvider
}

func NewLocationHandler(service domain.LocationServiceProvider) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

func (l *LocationHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	cursor := r.URL.Query().Get("cursor")
	limit := getLimit(r)

	logrus.Print("query: ", query)

	filters := make(repo.Filters)
	if query != "" {
		filters["name:ilike"] = query
	}

	locations, nextCursor, err := l.service.ListAll(r.Context(), cursor, limit, filters)

	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Cursor", nextCursor)
	render.JSON(w, r, locations)
}

func (l *LocationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	location, err := l.service.GetByID(r.Context(), id)
	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, location)
}

func (l *LocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &domain.LocationForm{}
	if err := render.Bind(r, data); err != nil {
		response.RespondForError(w, r, err)
		return
	}

	location := &domain.Location{
		Name: data.Name,
	}

	if err := l.service.Create(r.Context(), location); err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, location)
}

func (l *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	data := &domain.LocationForm{}
	if err = render.Bind(r, data); err != nil {
		response.RespondForError(w, r, err)
		return
	}

	location := &domain.Location{
		ID:   id,
		Name: data.Name,
	}

	if err = l.service.Update(r.Context(), location); err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, location)
}

func (l *LocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}

	if err = l.service.Delete(r.Context(), id); err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
