package durak

import (
	"context"
	"durak/pkg/domain"
	"durak/pkg/repository"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	cardsAmount    = 36
	welcomeMessage = "Welcome to Durak online (offline) game. Please login and wait other players to join the game"
)

type Game struct {
	mu sync.RWMutex

	storage *repository.MemoryStorage
	players map[string]*domain.Player
	bito    []domain.Card

	isGameInProgress atomic.Bool

	logger *log.Entry
	port   string
}

func NewGame(port string, logger *log.Entry) *Game {
	return &Game{
		port:             port,
		logger:           logger,
		storage:          repository.NewMemoryStorage(),
		players:          make(map[string]*domain.Player),
		bito:             make([]domain.Card, 0, cardsAmount),
		isGameInProgress: atomic.Bool{},
	}
}

func (g *Game) run(ctx context.Context) error {
	if swapped := g.isGameInProgress.CompareAndSwap(false, true); !swapped {
		return ErrGameIsInProgress
	}

	players, err := g.storage.GetPlayers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get players")
	}

	if len(players)*6 > cardsAmount {
		return ErrTooManyPlayers
	}

	g.logger.Infof("game was started with players: %v", players)
	deck := domain.NewDeck()
	deck.Shuffle()

	// waiting for readiness of players
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

waitingForReadiness:
	for {
		select {
		case <-ticker.C:
			if g.arePlayersReady() {
				break waitingForReadiness
			}
		case <-ctx.Done():
			return errors.Join(err, ErrContextIsDone)
		}
	}

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

	g.logger.Info("all players are ready")
	fmt.Printf("Deck: %v\n", deck)
	for _, player := range players {
		fmt.Println(player)
	}

	return nil
}

func (g *Game) createPlayer(username string) error {
	if err := g.storage.CreatePlayer(g.logger.Context, username); err != nil {
		return fmt.Errorf("failed to create user, err: %w", err)
	}
	return nil
}

func (g *Game) playerIsReady(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		g.logger.Fatalf("usernameKey holds not string value but %T", username)
	}

	player, ok := g.players[username]
	if !ok {
		g.logger.Fatalf("player with username %s was not found", username)
	}
	player.SetReady()
}

func (g *Game) arePlayersReady() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, player := range g.players {
		if !player.IsReady() {
			return false
		}
	}

	return true
}
