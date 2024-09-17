package api

import (
	durak "durak/pkg/game"

	log "github.com/sirupsen/logrus"
)

type App struct {
	port string

	game *durak.Game

	logger *log.Entry
}

func NewApp(port string, game *durak.Game, logger *log.Entry) *App {
	return &App{
		port:   port,
		game:   game,
		logger: logger,
	}
}
