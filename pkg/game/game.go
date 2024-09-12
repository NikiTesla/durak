package durak

import (
	"durak/pkg/domain"
	"durak/pkg/repository"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	cardsAmount    = 36
	welcomeMessage = "Welcome to Durak online (offline) game. Please login and wait other players to join the game"
)

type Game struct {
	storage *repository.MemoryStorage
	players map[string]domain.Player

	bito []domain.Card

	logger *log.Entry
	port   string
}

func NewGame(port string, logger *log.Entry) *Game {
	return &Game{
		port:    port,
		logger:  logger,
		storage: repository.NewMemoryStorage(),
		players: make(map[string]domain.Player),
		bito:    make([]domain.Card, 0, cardsAmount),
	}
}

func (g *Game) Run() error {
	rtr := http.NewServeMux()

	rtr.HandleFunc("GET /hello", g.indexPage)
	rtr.HandleFunc("POST /login", g.loginHandler)
	rtr.HandleFunc("POST /register", g.registerHandler)

	rtr.Handle("/api/", http.StripPrefix("/api", g.jwtMiddleware(rtr)))
	rtr.HandleFunc("GET /api/healthz", g.healthz)
	rtr.HandleFunc("GET /api/start_game", g.startGame)

	return http.ListenAndServe(g.port, rtr)
}

func (g *Game) indexPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(welcomeMessage))
}

func (g *Game) healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK!"))
}

func (g *Game) createPlayer(username string) error {
	if err := g.storage.CreatePlayer(g.logger.Context, username); err != nil {
		return fmt.Errorf("failed to create user, err: %w", err)
	}
	return nil
}

func (g *Game) startGame(w http.ResponseWriter, r *http.Request) {
	players, err := g.storage.GetPlayers(r.Context())
	if err != nil {
		g.logger.WithError(err).Error("failed to get players")
		http.Error(w, "failed to get players", http.StatusInternalServerError)
		return
	}

	if len(players)*6 > cardsAmount {
		g.logger.Error("too many players registered for game")
		http.Error(w, "too many players, restart the game", http.StatusTooManyRequests)
		return
	}

	g.logger.Infof("game was started with players: %v", players)
	deck := domain.NewDeck()
	deck.Shuffle()

	// distribution of cards
	for range 6 {
		for _, player := range players {
			card, err := deck.GetCard()
			if err != nil {
				panic(err)
			}

			player.GetCard(card)
		}
	}

	fmt.Printf("Deck: %v\n", deck)
	for _, player := range players {
		fmt.Println(player)
	}
}
