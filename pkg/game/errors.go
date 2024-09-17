package durak

import "errors"

var (
	ErrContextIsDone    = errors.New("context is done")
	ErrGameIsInProgress = errors.New("game is already in progress")
	ErrTooManyPlayers   = errors.New("too many players, restart the game")
)
