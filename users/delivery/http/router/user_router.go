package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/users/delivery/http/handler"
	"github.com/sainak/bitsb/utils/jwt"
)

func RegisterRoutes(router *chi.Mux, service domain.UserServiceProvider, j *jwt.JWT) {
	h := handler.New(service)
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/register", h.Register)
	r.Route("/user", func(r chi.Router) {
		r.Use(jwt.Authenticator(j))
		r.Get("/", h.GetCurrentUser)
	})
	router.Mount("/auth", r)
}
