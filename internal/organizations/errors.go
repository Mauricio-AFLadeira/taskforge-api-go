package organizations

import "errors"

var (
	ErrValidation = errors.New("organizations: validation failed")
	ErrNotFound   = errors.New("organizations: not found")
	ErrForbidden  = errors.New("organizations: forbidden")
	ErrConflict   = errors.New("organizations: conflict")
)
