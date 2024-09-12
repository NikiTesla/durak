package domain

import (
	"strings"
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

func (p *Player) String() string {
	return p.name
}

func (p *Player) TakeCard(card Card) {
	p.cards = append(p.cards, card)
}

func (p *Player) GetHand() string {
	if len(p.cards) == 0 {
		return ""
	}

	result := strings.Builder{}
	cardImageRowsAmount := strings.Split(p.cards[0].String(), "\n")
	for i := range cardImageRowsAmount {
		for _, card := range p.cards {
			result.WriteString(strings.Split(card.String(), "\n")[i] + "\t\t")
		}
		result.WriteString("\n")
	}

	return result.String()
}

func (p *Player) SetReady() {
	p.isReady = true
}

func (p *Player) IsReady() bool {
	return p.isReady
}
