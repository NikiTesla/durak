package durak

import (
	"context"
	"durak/pkg/domain"
	"durak/pkg/repository"
	"errors"
	"fmt"
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
}

func NewGame(logger *log.Entry) *Game {
	return &Game{
		logger:           logger,
		storage:          repository.NewMemoryStorage(),
		players:          make(map[string]*domain.Player),
		bito:             make([]domain.Card, 0, cardsAmount),
		isGameInProgress: atomic.Bool{},
	}
}

func (g *Game) Start(ctx context.Context) error {
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

func (g *Game) AddCardToTable(card *domain.Card) error {
	var valid bool
	for _, baseCard := range domain.Cards {
		if card.Suit == baseCard.Suit && card.Rank == baseCard.Rank {
			valid = true
		}
	}

	if !valid {
		return errors.New("card is unknwon")
	}

	g.tableCards[card] = nil

	return nil
}

func (g *Game) clearTable() error {
	for card1, card2 := range g.tableCards {
		g.bito = append(g.bito, *card1, *card2)
	}
	return nil
}
