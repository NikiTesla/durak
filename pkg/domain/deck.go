package domain

import (
	"fmt"
	"math/rand"
)

type Deck struct {
	cards []Card
}

func NewDeck() Deck {
	return Deck{
		cards: Cards,
	}
}

func (d *Deck) Shuffle() {
	for i := range d.cards {
		j := rand.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) String() string {
	return fmt.Sprintf("%v", d.cards)
}

func (d *Deck) GetCard() (Card, error) {
	if len(d.cards) == 0 {
		return Card{}, ErrDeckIsEmpty
	}

	card := d.cards[0]
	d.cards = d.cards[1:]

	return card, nil
}
