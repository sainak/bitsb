package errors

import (
	"net/http"
)

// Error is a custom error wrapper with more information
type Error struct {
	StatusCode int
	Message    string
}

func New(statusCode int, message string) error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func GetErrStatusCode(err error) int {
	switch e := err.(type) {
	case *Error:
		return e.StatusCode
	default:
		return http.StatusInternalServerError
	}
}

var (
	// BadRequest errors
	ErrBadRequest = &Error{
		StatusCode: http.StatusBadRequest,
		Message:    "bad request",
	}
	ErrBadCursor = &Error{
		StatusCode: http.StatusBadRequest,
		Message:    "bad cursor",
	}
	ErrBadInputParam = &Error{
		StatusCode: http.StatusBadRequest,
		Message:    "bad input param, check your input params",
	}
	ErrEmptyRequest = &Error{
		StatusCode: http.StatusBadRequest,
		Message:    "empty request body",
	}
	ErrInvalidLocation = &Error{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid location",
	}

	// Unauthorized errors
	ErrUnauthorized = &Error{
		StatusCode: http.StatusUnauthorized,
		Message:    "unauthorized",
	}
	ErrInvalidToken = &Error{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid token",
	}
	ErrExpiredToken = &Error{
		StatusCode: http.StatusUnauthorized,
		Message:    "expired token",
	}
	ErrInvalidCredentials = &Error{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid credentials",
	}

	// Forbidden errors
	ErrForbidden = &Error{
		StatusCode: http.StatusForbidden,
		Message:    "forbidden",
	}

	// Not Found errors
	ErrNotFound = &Error{
		StatusCode: http.StatusNotFound,
		Message:    "entity not found",
	}

	// Conflict errors
	ErrConflict = &Error{
		StatusCode: http.StatusConflict,
		Message:    "database conflict occurred",
	}
	ErrEntityAlreadyExist = &Error{
		StatusCode: http.StatusConflict,
		Message:    "entity already exist",
	}

	// Internal Server errors
	ErrInternalServerError = &Error{
		StatusCode: http.StatusInternalServerError,
		Message:    "an internal server error occurred, we are checking...",
	}
)
