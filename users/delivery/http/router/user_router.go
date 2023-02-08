package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/users/delivery/http/handler"
)

func RegisterRoutes(
	router *chi.Mux,
	service domain.UserServiceProvider,
	jwtMiddleware func(next http.Handler) http.Handler,
) {
	h := handler.New(service)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
		r.Post("/register", h.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Get("/user", h.GetCurrentUser)
	})
}
