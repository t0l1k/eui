package game

import (
	"fmt"
	"image/color"

	"github.com/t0l1k/eui"
)

type Card struct {
	suit Suit
	face Face
}

func NewCard(suit Suit, face Face) *Card {
	c := &Card{suit: suit, face: face}
	return c
}

func (s *Card) SetCard(value *Card) {
	s.face = value.face
	s.suit = value.suit
}

func (s *Card) Eq(other *Card) bool {
	if other == nil || s == nil {
		return false
	}
	return s.face == other.face && s.suit == other.suit
}

func (s *Card) EqFace(other *Card) bool {
	if other == nil || s == nil {
		return false
	}
	return s.face == other.face
}

func (s *Card) EqSuit(other *Card) bool {
	if other == nil || s == nil {
		return false
	}
	return s.suit == other.suit
}

func (s Card) Color() (col color.Color) {
	switch s.suit {
	case Hearts, Diamonds:
		col = eui.Red
	case Clubs, Spades:
		col = eui.Black
	}
	return col
}

func (c Card) StringShort() string {
	return fmt.Sprintf("%2v", c.face.String()) + c.suit.String()
}

func (c Card) String() string {
	return fmt.Sprintf("[%2v%v]", c.face.String(), c.suit.String())
}
