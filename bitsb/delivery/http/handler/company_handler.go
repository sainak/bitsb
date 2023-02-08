package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/utils/repo"
	"github.com/sainak/bitsb/utils/response"
)

type CompanyHandler struct {
	service domain.CompanyServiceProvider
}

func NewCompanyHandler(service domain.CompanyServiceProvider) *CompanyHandler {
	return &CompanyHandler{
		service: service,
	}
}

func (h CompanyHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	cursor := r.URL.Query().Get("cursor")
	limit := getLimit(r)

	filters := make(repo.Filters)
	if query != "" {
		filters["name:ilike"] = query
	}

	companies, nextCursor, err := h.service.ListAll(r.Context(), cursor, limit, filters)

	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Cursor", nextCursor)
	render.JSON(w, r, companies)
}

func (h CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	company, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, company)
}

func (h CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &domain.CompanyForm{}
	if err := render.Bind(r, data); err != nil {
		response.RespondForError(w, r, err)
		return
	}

	company := &domain.Company{
		Name:       data.Name,
		LocationID: data.LocationID,
	}

	if err := h.service.Create(r.Context(), company); err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, company)
}

func (h CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}

	data := &domain.CompanyForm{}
	if err := render.Bind(r, data); err != nil {
		response.RespondForError(w, r, err)
		return
	}

	company := &domain.Company{
		ID:         id,
		Name:       data.Name,
		LocationID: data.LocationID,
	}

	if err := h.service.Update(r.Context(), company); err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, company)
}

func (h CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
