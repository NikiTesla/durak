package repository

import (
	"context"
	"durak/pkg/domain"
	"sync"
)

type MemoryStorage struct {
	players map[string]*domain.Player
	mx      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		players: make(map[string]*domain.Player),
	}
}

func (m *MemoryStorage) CreatePlayer(_ context.Context, username string) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.players[username] = domain.NewPlayer(username)

	return nil
}

func (m *MemoryStorage) DeletePlayer(_ context.Context, username string) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.players, username)

	return nil
}

func (m *MemoryStorage) GetPlayers(_ context.Context) (map[string]*domain.Player, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.players, nil
}

func (m *MemoryStorage) GetPlayer(username string) (*domain.Player, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	player, ok := m.players[username]
	if !ok {
		return player, domain.ErrNotFound
	}
	return player, nil
}
