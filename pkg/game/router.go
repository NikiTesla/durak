package durak

import (
	"errors"
	"net/http"
)

func (g *Game) Run() error {
	rtr := g.InitRouter()

	g.logger.Infof("game service was started and is listening addr %s", g.port)
	return http.ListenAndServe(g.port, rtr)
}

func (g *Game) InitRouter() http.Handler {
	rtr := http.NewServeMux()

	rtr.HandleFunc("GET /hello", g.indexPage)
	rtr.HandleFunc("POST /login", g.loginHandler)
	rtr.HandleFunc("POST /register", g.registerHandler)

	playerRtr := http.NewServeMux()
	playerRtr.HandleFunc("POST /ready", g.playerIsReady)
	rtr.Handle("/player/", http.StripPrefix("/player", g.playerMiddleware(playerRtr)))

	adminRtr := http.NewServeMux()
	adminRtr.HandleFunc("/start_game", g.startGame)
	rtr.Handle("/admin/", http.StripPrefix("/admin", g.adminMiddleware(adminRtr)))

	return rtr
}

func (g *Game) indexPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(welcomeMessage))
}

func (g *Game) startGame(w http.ResponseWriter, r *http.Request) {
	if err := g.run(r.Context()); err != nil {
		switch {
		case errors.Is(err, ErrGameIsInProgress):
			http.Error(w, "Game is already in progress", http.StatusAlreadyReported)
		case errors.Is(err, ErrContextIsDone):
			w.Write([]byte("You've ended the game"))
		case errors.Is(err, ErrTooManyPlayers):
			http.Error(w, "There are too many players. Please, restart the game", http.StatusTooManyRequests)
		default:
			g.logger.WithError(err).Error("game failed")
			http.Error(w, "error occured during the game", http.StatusInternalServerError)
		}
	}
}
