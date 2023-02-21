package apperrors

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/lib/pq"
)

func ParseError(err error) *Error {
	switch e := err.(type) {
	case *Error:
		return e

	case *strconv.NumError:
		return &Error{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf(`%q is not a valid number.`, e.Num),
		}

	case *pq.Error:
		//https://github.com/lib/pq/blob/922c00e176fb3960d912dc2c7f67ea2cf18d27b0/error.go#L78
		switch e.Code {
		case "23502":
			// not-null constraint violation
			return &Error{
				StatusCode: http.StatusConflict,
				Message:    fmt.Sprint("some required data was left out:", e.Message),
			}
		case "23505":
			// unique constraint violation
			return &Error{
				StatusCode: http.StatusConflict,
				Message:    fmt.Sprint("this record already exists:", e.Message),
			}
		}

	default:
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		case io.EOF:
			return ErrEmptyRequest
		}
	}

	return ErrInternalServerError
}
