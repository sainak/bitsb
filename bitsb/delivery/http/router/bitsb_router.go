package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sainak/bitsb/bitsb"
	locationHandler "github.com/sainak/bitsb/bitsb/delivery/http/handler"
	"github.com/sainak/bitsb/users"
	"github.com/sainak/bitsb/users/delivery/http/middleware"
)

func RegisterLocationRoutes(
	router *chi.Mux,
	service bitsb.LocationServiceProvider,
	jwtMiddleware func(next http.Handler) http.Handler,
) {
	h := locationHandler.NewLocationHandler(service)

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", h.ListAll)
			r.With(middleware.AccessAbove(users.Admin)).Post("/", h.Create)
		})
		r.Route("/location", func(r chi.Router) {
			r.Get("/{id}", h.GetByID)
			r.With(middleware.AccessAbove(users.Admin)).Patch("/{id}", h.Update)
			r.With(middleware.AccessAbove(users.Admin)).Delete("/{id}", h.Delete)
		})
	})
}

func RegisterBusRouteRoutes(
	router *chi.Mux,
	service bitsb.BusRouteServiceProvider,
	jwtMiddleware func(next http.Handler) http.Handler,
) {
	h := locationHandler.NewBusRouteHandler(service)
	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Route("/bus-routes", func(r chi.Router) {
			r.Get("/", h.ListAll)
			r.Get("/for-user", h.BusesForUser)
			r.With(middleware.AccessAbove(users.Admin)).Post("/", h.Create)
		})
		r.Route("/bus-route", func(r chi.Router) {
			r.Get("/{id}", h.GetByID)
			r.Get("/{id}/ticket-price", h.TicketPrice)
			r.With(middleware.AccessAbove(users.Admin)).Patch("/{id}", h.Update)
			r.With(middleware.AccessAbove(users.Admin)).Delete("/{id}", h.Delete)
		})
	})
}
