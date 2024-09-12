package domain

import "errors"

var (
	ErrNotFound    = errors.New("Not found")
	ErrDeckIsEmpty = errors.New("Deck is empty")
)
