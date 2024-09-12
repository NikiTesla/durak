package domain

import (
	"fmt"
)

type Card struct {
	rank Rank
	suit Suit
}

type Rank int

type Suit string

const (
	Hearts   Suit = "Hearts"
	Diamonds Suit = "Diamonds"
	Clubs    Suit = "Clubs"
	Spades   Suit = "Spades"
)

const (
	Six Rank = iota + 6
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

var Cards = []Card{
	{Six, Hearts}, {Seven, Hearts}, {Eight, Hearts}, {Nine, Hearts}, {Ten, Hearts}, {Jack, Hearts}, {Queen, Hearts}, {King, Hearts}, {Ace, Hearts},
	{Six, Diamonds}, {Seven, Diamonds}, {Eight, Diamonds}, {Nine, Diamonds}, {Ten, Diamonds}, {Jack, Diamonds}, {Queen, Diamonds}, {King, Diamonds}, {Ace, Diamonds},
	{Six, Clubs}, {Seven, Clubs}, {Eight, Clubs}, {Nine, Clubs}, {Ten, Clubs}, {Jack, Clubs}, {Queen, Clubs}, {King, Clubs}, {Ace, Clubs},
	{Six, Spades}, {Seven, Spades}, {Eight, Spades}, {Nine, Spades}, {Ten, Spades}, {Jack, Spades}, {Queen, Spades}, {King, Spades}, {Ace, Spades},
}

func (c *Card) String() string {
	return fmt.Sprintf("%d %s", c.rank, c.suit)
}
