package response

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse represent the response error struct
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func RespondForError(w http.ResponseWriter, r *http.Request, err error) {
	er := ErrorResponse{
		Status: "unexpected_error",
	}
	er.Message = err.Error()
	// todo: set status
	render.Status(r, 400)
	render.JSON(w, r, er)
}
