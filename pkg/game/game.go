package durak

import (
	"context"
	"durak/pkg/domain"
	"durak/pkg/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	cardsAmount    = 36
	welcomeMessage = "Welcome to Durak online (offline) game. Please login and wait other players to join the game"
)

type Game struct {
	storage *repository.MemoryStorage

	players      map[string]*domain.Player
	bito         []domain.Card
	tableCards   map[*domain.Card]*domain.Card
	playersQueue []*domain.Player

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

func (g *Game) updatePlayers(ctx context.Context) error {
	players, err := g.storage.GetPlayers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get players, err: %w", err)
	}
	g.players = players

	return nil
}

func (g *Game) run(ctx context.Context) error {
	if swapped := g.isGameInProgress.CompareAndSwap(false, true); !swapped {
		return ErrGameIsInProgress
	}

	if err := g.updatePlayers(ctx); err != nil {
		return fmt.Errorf("updating players, err: %w", err)
	}

	// getting players' queue
	if len(g.players)*6 > cardsAmount {
		return ErrTooManyPlayers
	}

	g.playersQueue = make([]*domain.Player, 0, len(g.players))
	for _, player := range g.players {
		g.playersQueue = append(g.playersQueue, player)
	}

	g.logger.Infof("game was started with players: %v", g.players)
	deck := domain.NewDeck()
	deck.Shuffle()

	// waiting for readiness of players
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	g.logger.Info("waiting for players to be ready")
waitingForReadiness:
	for {
		select {
		case <-ticker.C:
			if g.arePlayersReady() {
				break waitingForReadiness
			}
		case <-ctx.Done():
			return errors.Join(ctx.Err(), ErrContextIsDone)
		}
	}

	// distribution of cards
	for range 6 {
		for _, player := range g.players {
			card, err := deck.GetCard()
			if err != nil {
				panic(err)
			}

			player.TakeCard(card)
		}
	}

	g.logger.Info("all players are ready")
	fmt.Printf("Trump is:\n%s\n", deck.GetTrumpCard())

	for range ctx.Done() {
		return errors.Join(ctx.Err(), ErrContextIsDone)
	}

	return nil
}

func (g *Game) arePlayersReady() bool {
	for _, player := range g.players {
		if !player.IsReady() {
			return false
		}
	}

	return true
}

func (g *Game) createPlayer(username string) error {
	if err := g.storage.CreatePlayer(g.logger.Context, username); err != nil {
		return fmt.Errorf("failed to create user, err: %w", err)
	}
	return nil
}

func (g *Game) playerIsReady(w http.ResponseWriter, r *http.Request) {
	player, err := g.getPlayer(r)
	if err != nil {
		g.logger.WithError(err).Error("getting player")
		http.Error(w, "cannot get player's info", http.StatusInternalServerError)
		return
	}

	player.SetReady()
}

func (g *Game) getPlayersCards(w http.ResponseWriter, r *http.Request) {
	player, err := g.getPlayer(r)
	if err != nil {
		g.logger.WithError(err).Error("getting player")
		http.Error(w, "cannot get player's info", http.StatusInternalServerError)
		return
	}

	g.logger.Infof("player's hand is:\n %s", player.GetHand())
	json.NewEncoder(w).Encode(map[string]string{"hand": player.GetHand()})
}

func (g *Game) getPlayer(r *http.Request) (*domain.Player, error) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		return nil, fmt.Errorf("usernameKey holds not string value but %T", username)
	}

	player, ok := g.players[username]
	if !ok {
		return nil, fmt.Errorf("player with username %s was not found", username)
	}

	return player, nil
}
