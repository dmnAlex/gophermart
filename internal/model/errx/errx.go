package errx

import "errors"

var (
	ErrBadRequest      = errors.New("bad request")
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrAlreadyAccepted = errors.New("already accepted")
	ErrConflict        = errors.New("conflict")
	ErrUnprocessable   = errors.New("unprocessable")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInternalError   = errors.New("internal server error")
)
