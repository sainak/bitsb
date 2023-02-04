package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"message":      "pong",
		"current_time": time.Now(),
	})
}
