package durak

import "errors"

var (
	ErrUsernameIsEmpty = errors.New("Username is empty")
	ErrPasswordIsEmpty = errors.New("Password is empty")
)
