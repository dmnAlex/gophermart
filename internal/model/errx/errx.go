package errx

import "errors"

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrBadRequest       = errors.New("bad request")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrConflict         = errors.New("conflict")
	ErrUnprocessable    = errors.New("unprocessable")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInternalError    = errors.New("internal server error")
)
