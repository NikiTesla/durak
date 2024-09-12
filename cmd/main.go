package main

import (
	durak "durak/pkg/game"

	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.NewEntry(log.StandardLogger())
	logger.Logger.SetLevel(log.DebugLevel)

	game := durak.NewGame(":7070", logger)

	if err := game.Run(); err != nil {
		log.WithError(err).Fatal("game failed")
	}
}
