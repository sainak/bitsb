package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	locationHandler "github.com/sainak/bitsb/bitsb/delivery/http/handler"
	"github.com/sainak/bitsb/domain"
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
			r.Post("/", h.Create)
		})
		r.Route("/location", func(r chi.Router) {
			r.Get("/{id}", h.GetByID)
			r.Patch("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
}

func RegisterCompanyRoutes(
	router *chi.Mux,
	service domain.CompanyServiceProvider,
	jwtMiddleware func(next http.Handler) http.Handler,
) {
	h := locationHandler.NewCompanyHandler(service)

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Route("/companies", func(r chi.Router) {
			r.Get("/", h.ListAll)
			r.Post("/", h.Create)
		})
		r.Route("/company", func(r chi.Router) {
			r.Get("/{id}", h.GetByID)
			r.Patch("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
}
