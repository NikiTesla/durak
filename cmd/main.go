package main

import (
	"durak/pkg/api"
	durak "durak/pkg/game"

	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.NewEntry(log.StandardLogger())
	logger.Logger.SetLevel(log.DebugLevel)

	game := durak.NewGame(logger)

	app := api.NewApp(":7070", game, logger)

	if err := app.Start(); err != nil {
		log.WithError(err).Fatal("game failed")
	}
}
