package domain

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrDeckIsEmpty = errors.New("deck is empty")
)
