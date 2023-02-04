package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/sainak/bitsb/root/delivery/http/handler"
)

func RegisterRoutes(router *chi.Mux) {
	router.Get("/ping", handler.Ping)
}
