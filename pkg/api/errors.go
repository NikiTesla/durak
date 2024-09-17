package api

import "errors"

var (
	ErrUsernameIsEmpty            = errors.New("username is empty")
	ErrPasswordIsEmpty            = errors.New("password is empty")
	ErrAuthorisationHeaderMissing = errors.New("authorization header missing")
	ErrInvalidToken               = errors.New("invalid token")
)
