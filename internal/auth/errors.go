package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	ErrInvalidRefresh     = errors.New("auth: invalid refresh token")
	ErrValidation         = errors.New("auth: validation failed")
)
