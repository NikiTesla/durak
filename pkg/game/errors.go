package durak

import "errors"

var (
	ErrUsernameIsEmpty            = errors.New("username is empty")
	ErrPasswordIsEmpty            = errors.New("password is empty")
	ErrAuthorisationHeaderMissing = errors.New("authorization header missing")
	ErrInvalidToken               = errors.New("invalid token")

	ErrContextIsDone    = errors.New("context is done")
	ErrGameIsInProgress = errors.New("game is already in progress")
	ErrTooManyPlayers   = errors.New("too many players, restart the game")
)
