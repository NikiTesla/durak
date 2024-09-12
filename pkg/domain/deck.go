package domain

import (
	"fmt"
	"math/rand"
	"sync"
)

type Deck struct {
	mu    sync.RWMutex
	cards []Card
}

func NewDeck() *Deck {
	return &Deck{
		mu:    sync.RWMutex{},
		cards: Cards,
	}
}

func (d *Deck) Shuffle() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i := range d.cards {
		j := rand.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) String() string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return fmt.Sprintf("%v", d.cards)
}

func (d *Deck) GetCard() (Card, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.cards) == 0 {
		return Card{}, ErrDeckIsEmpty
	}

	card := d.cards[0]
	d.cards = d.cards[1:]

	return card, nil
}
