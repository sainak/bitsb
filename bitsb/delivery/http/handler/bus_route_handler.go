package handler

import (
	"net/http"

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
