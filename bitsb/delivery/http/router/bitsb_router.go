package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	locationHandler "github.com/sainak/bitsb/bitsb/delivery/http/handler"
	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/utils/middleware"
)

func RegisterLocationRoutes(
	router *chi.Mux,
	service domain.LocationServiceProvider,
	jwtMiddleware func(next http.Handler) http.Handler,
) {
	h := locationHandler.NewLocationHandler(service)

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", h.ListAll)
			r.With(middleware.AccessAbove(domain.Admin)).Post("/", h.Create)
		})
		r.Route("/location", func(r chi.Router) {
			r.Get("/{id}", h.GetByID)
			r.With(middleware.AccessAbove(domain.Admin)).Patch("/{id}", h.Update)
			r.With(middleware.AccessAbove(domain.Admin)).Delete("/{id}", h.Delete)
		})
	})
}

func RegisterBusRouteRoutes(
	router *chi.Mux,
	service domain.BusRouteServiceProvider,
	jwtMiddleware func(next http.Handler) http.Handler,
) {
	h := locationHandler.NewBusRouteHandler(service)
	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.With(middleware.AccessAbove(domain.Admin)).Post("/bus-routes", h.Create)
		r.Route("/bus-route", func(r chi.Router) {
			r.With(middleware.AccessAbove(domain.Admin)).Patch("/{id}", h.Update)
			r.With(middleware.AccessAbove(domain.Admin)).Delete("/{id}", h.Delete)
		})
	})
}
