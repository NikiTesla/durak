package durak

import (
	"durak/pkg/domain"
	"encoding/json"
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
	playerRtr.HandleFunc("GET /get_hand", g.getPlayersCards)
	playerRtr.HandleFunc("POST /add_card", g.addCardToTable)
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

func (g *Game) addCardToTable(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var card domain.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		g.logger.WithError(err).Error("decoding card data from request")
		http.Error(w, "failed to read your card data", http.StatusInternalServerError)
		return
	}

	var valid bool
	for _, baseCard := range domain.Cards {
		if card.Suit == baseCard.Suit && card.Rank == baseCard.Rank {
			valid = true
		}
	}

	if !valid {
		http.Error(w, "you've sent unknown card", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"card": card.String()})
}
