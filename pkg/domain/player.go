package domain

import (
	"fmt"
)

type Player struct {
	name    string
	cards   []Card
	isReady bool
}

func NewPlayer(name string) *Player {
	return &Player{
		name:    name,
		cards:   make([]Card, 0, 6),
		isReady: false,
	}
}

func (p *Player) GetCard(card Card) {
	p.cards = append(p.cards, card)
}

func (p *Player) String() string {
	return fmt.Sprintf("Player %s has cards: %v", p.name, p.cards)
}

func (p *Player) SetReady() {
	p.isReady = true
}

func (p *Player) IsReady() bool {
	return p.isReady
}
