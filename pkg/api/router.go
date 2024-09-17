package api

import (
	"durak/pkg/domain"
	durak "durak/pkg/game"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const welcomeMessage = "Welcome to Durak online (offline) game. Please login and wait other players to join the game"

func (a *App) Start() error {
	rtr := a.InitRouter()

	a.logger.Infof("game service was started and is listening addr %s", a.port)
	return http.ListenAndServe(a.port, rtr)
}

func (a *App) InitRouter() http.Handler {
	rtr := http.NewServeMux()

	rtr.HandleFunc("GET /hello", a.indexPage)
	rtr.HandleFunc("POST /login", a.loginHandler)
	rtr.HandleFunc("POST /register", a.registerHandler)

	playerRtr := http.NewServeMux()
	playerRtr.HandleFunc("POST /ready", a.playerIsReady)
	playerRtr.HandleFunc("GET /get_hand", a.getPlayersCards)
	playerRtr.HandleFunc("POST /add_card", a.addCardToTable)
	rtr.Handle("/player/", http.StripPrefix("/player", a.playerMiddleware(playerRtr)))

	adminRtr := http.NewServeMux()
	adminRtr.HandleFunc("/start_game", a.startGame)
	rtr.Handle("/admin/", http.StripPrefix("/admin", a.adminMiddleware(adminRtr)))

	return rtr
}

func (a *App) indexPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(welcomeMessage))
}

func (a *App) startGame(w http.ResponseWriter, r *http.Request) {
	if err := a.game.Start(r.Context()); err != nil {
		switch {
		case errors.Is(err, durak.ErrGameIsInProgress):
			http.Error(w, "Game is already in progress", http.StatusAlreadyReported)
		case errors.Is(err, durak.ErrContextIsDone):
			w.Write([]byte("You've ended the game"))
		case errors.Is(err, durak.ErrTooManyPlayers):
			http.Error(w, "There are too many players. Please, restart the game", http.StatusTooManyRequests)
		default:
			a.logger.WithError(err).Error("game failed")
			http.Error(w, "error occured during the game", http.StatusInternalServerError)
		}
	}
}

func (a *App) addCardToTable(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var card domain.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		a.logger.WithError(err).Error("decoding card data from request")
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

func (a *App) playerIsReady(w http.ResponseWriter, r *http.Request) {
	username, err := a.getUsernameFromRequest(r)
	if err != nil {
		http.Error(w, "failed to get username from request", http.StatusUnauthorized)
		return
	}

	if err := a.game.PlayerIsReady(username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getPlayersCards(w http.ResponseWriter, r *http.Request) {
	username, err := a.getUsernameFromRequest(r)
	if err != nil {
		http.Error(w, "failed to get username from request", http.StatusUnauthorized)
		return
	}

	cardsData, err := a.game.GetPlayersCards(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(cardsData)
}

func (a *App) getUsernameFromRequest(r *http.Request) (string, error) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		return "", fmt.Errorf("usernameKey holds not string value but %T", username)
	}

	return username, nil
}
