package errs

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrValidation = errors.New("validation error")
)
