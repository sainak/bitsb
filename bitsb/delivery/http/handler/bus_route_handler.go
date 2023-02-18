package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/api"
	"github.com/sainak/bitsb/domain/errors"
	"github.com/sainak/bitsb/domain/middleware"
	"github.com/sainak/bitsb/pkg/handler"
)

type TicketPriceResponse struct {
	TicketPrice float64 `json:"ticket_price"`
}

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
	limit := handler.GetLimit(r)

	l := r.URL.Query().Get("locations")
	// convert csv string of locations to int64 array
	var locations []int64
	if l != "" {
		lA := strings.Split(l, ",")
		locations = make([]int64, len(lA))
		for i, v := range lA {
			locations[i], _ = strconv.ParseInt(v, 10, 64)
		}
	}

	logrus.Debug(locations)

	busRoutes, nextCursor, err := h.service.ListAll(r.Context(), cursor, limit, locations)
	if err != nil {
		logrus.Error(err)
		api.RespondForError(w, r, err)
		return
	}

	w.Header().Set("X-Cursor", nextCursor)
	render.JSON(w, r, busRoutes)
}

func (h *BusRouteHandler) BusesForUser(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limit := handler.GetLimit(r)

	user := r.Context().Value(middleware.UserCtxKey).(domain.User)
	homeLocation := user.HomeLocationID.ValueOrZero()
	workLocation := user.WorkLocationID.ValueOrZero()
	if homeLocation == 0 || workLocation == 0 {
		api.RespondForError(w, r, errors.New(http.StatusBadRequest, "user has no home or work location"))
		return
	}

	busRoutes, nextCursor, err := h.service.ListAll(r.Context(), cursor, limit, []int64{homeLocation, workLocation})
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	w.Header().Set("X-Cursor", nextCursor)
	render.JSON(w, r, busRoutes)
}

func (h *BusRouteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}
	busRoute, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, busRoute)
}

func (h *BusRouteHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &domain.BusRouteForm{}
	err := render.Bind(r, data)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	busRoute := &domain.BusRoute{
		Name:        data.Name,
		Number:      data.Number,
		StartTime:   data.StartTime,
		EndTime:     data.EndTime,
		Interval:    data.Interval,
		LocationIDS: data.LocationIDS,
		MaxPrice:    data.MaxPrice,
		MinPrice:    data.MinPrice,
	}

	if err = h.service.Create(r.Context(), busRoute); err != nil {
		api.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, busRoute)
}

func (h *BusRouteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}
	data := &domain.BusRouteForm{}
	err = render.Bind(r, data)
	if err != nil {
		api.RespondForError(w, r, err)
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

	if err = h.service.Update(r.Context(), busRoute); err != nil {
		api.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, busRoute)
}

func (h *BusRouteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		api.RespondForError(w, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.NoContent(w, r)
}

func (h *BusRouteHandler) TicketPrice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		logrus.Error(err, "id")
		api.RespondForError(w, r, err)
		return
	}
	start, err := strconv.ParseInt(r.URL.Query().Get("start"), 10, 64)
	if err != nil {
		logrus.Error(err, "start")
		api.RespondForError(w, r, err)
		return
	}
	end, err := strconv.ParseInt(r.URL.Query().Get("end"), 10, 64)
	if err != nil {
		logrus.Error(err, "end")

		api.RespondForError(w, r, err)
		return
	}

	ticketPrice, err := h.service.CalculateTicketPrice(r.Context(), id, start, end)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	render.JSON(w, r, render.M{
		"ticket_price": ticketPrice,
	})
}
