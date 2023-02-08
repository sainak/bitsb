package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrBadCursor          = errors.New("bad cursor")
	ErrNotFound           = errors.New("not found")
)
