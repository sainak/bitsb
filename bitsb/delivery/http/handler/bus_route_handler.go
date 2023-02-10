package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/utils/response"
)

type BusRouteHandler struct {
	service domain.BusRouteServiceProvider
}

func NewBusRouteHandler(service domain.BusRouteServiceProvider) *BusRouteHandler {
	return &BusRouteHandler{
		service: service,
	}
}

func (h *BusRouteHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limit := getLimit(r)

	l := r.URL.Query().Get("locations")
	var locations []int64
	if l != "" {
		lA := strings.Split(l, ",")
		locations = make([]int64, len(lA))
		for i, v := range lA {
			locations[i], _ = strconv.ParseInt(v, 10, 64)
		}
	}

	busRoutes, nextCursor, err := h.service.ListAll(r.Context(), cursor, limit, locations)
	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Cursor", nextCursor)
	render.JSON(w, r, busRoutes)
}

func (h *BusRouteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	busRoute, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, busRoute)
}

func (h *BusRouteHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &domain.BusRouteForm{}
	err := render.Bind(r, data)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}

	busRoute := &domain.BusRoute{
		Name:        data.Name,
		Number:      data.Number,
		StartTime:   data.StartTime,
		EndTime:     data.EndTime,
		Interval:    data.Interval,
		LocationIDS: data.LocationIDS,
	}

	err = h.service.Create(r.Context(), busRoute)
	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, busRoute)
}

func (h *BusRouteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	data := &domain.BusRouteForm{}
	err = render.Bind(r, data)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}

	busRoute := &domain.BusRoute{
		ID:          id,
		Name:        data.Name,
		Number:      data.Number,
		StartTime:   data.StartTime,
		EndTime:     data.EndTime,
		Interval:    data.Interval,
		LocationIDS: data.LocationIDS,
	}

	err = h.service.Update(r.Context(), busRoute)
	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, busRoute)
}

func (h *BusRouteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		response.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.NoContent(w, r)
}
