package api

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sainak/bitsb/apperrors"
)

// ErrorResponse represent the api error struct
type ErrorResponse struct {
	Message string `json:"message"`
}

func RespondForError(w http.ResponseWriter, r *http.Request, err error) {
	e := apperrors.ParseError(err)
	render.Status(r, e.StatusCode)
	render.JSON(w, r, ErrorResponse{e.Message})
}
