package handler

import (
	"net/http"
	"strconv"

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
