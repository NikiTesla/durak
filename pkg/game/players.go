package durak

import (
	"context"
	"durak/pkg/domain"
	"encoding/json"
	"fmt"
)

func (g *Game) updatePlayers(ctx context.Context) error {
	players, err := g.storage.GetPlayers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get players, err: %w", err)
	}
	g.players = players

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

func (g *Game) CreatePlayer(ctx context.Context, username string) error {
	if err := g.storage.CreatePlayer(ctx, username); err != nil {
		return fmt.Errorf("failed to create user, err: %w", err)
	}
	return nil
}

func (g *Game) PlayerIsReady(username string) error {
	player, err := g.getPlayer(username)
	if err != nil {
		return fmt.Errorf("getting player by username, err: %w", err)
	}

	player.SetReady()

	return nil
}

func (g *Game) GetPlayersCards(username string) ([]byte, error) {
	player, err := g.getPlayer(username)
	if err != nil {
		return nil, fmt.Errorf("getting player by username, err: %w", err)
	}

	g.logger.Infof("player's hand is:\n %s", player.GetHand())

	data, err := json.Marshal(map[string]string{"hand": player.GetHand()})
	if err != nil {
		return nil, fmt.Errorf("marshalling data, err: %w", err)
	}

	return data, nil
}

func (g *Game) getPlayer(username string) (*domain.Player, error) {
	player, ok := g.players[username]
	if !ok {
		return nil, fmt.Errorf("player with username %s was not found", username)
	}

	return player, nil
}
