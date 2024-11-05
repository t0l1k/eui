package deck

import (
	"fmt"
	"image/color"

	"github.com/t0l1k/eui/colors"
)

type Card struct {
	suit Suit
	face Face
}

func NewCard(suit Suit, face Face) *Card {
	c := &Card{suit: suit, face: face}
	return c
}

func (s *Card) GetCard() (value *Card) { return s }

func (s *Card) SetCard(value *Card) {
	s.face = value.face
	s.suit = value.suit
}

func (s *Card) Eq(other *Card) bool {
	return (s == nil && other == nil) || s.face.IsEq(other.face) && s.suit.IsEq(other.suit)
	// return !(other == nil || s == nil) && s.face == other.face && s.suit == other.suit
}

func (s *Card) IsEqFace(other *Card) bool {
	return !(other == nil || s == nil) && s.face == other.face
}

func (s *Card) EqSuit(other *Card) bool {
	return !(other == nil || s == nil) && s.suit.IsEq(other.suit)
}

func (s *Card) EqColor(other *Card) bool {
	return !(s == nil || other == nil) && s.suit.EqColor(other.suit)
}

func (s *Card) IsOneLess(other *Card) bool {
	return !(s == nil || other == nil) && s.face.IsOneLess(other.face)
}

func (s *Card) IsOneHigh(other *Card) bool {
	return !(s == nil || other == nil) && s.face.IsOneHigh(other.face)
}

func (s Card) Color() (col color.Color) {
	switch s.suit {
	case Hearts, Diamonds:
		col = colors.Red
	case Clubs, Spades:
		col = colors.Black
	}
	return col
}

func (c Card) StringShort() string {
	return fmt.Sprintf("%2v", c.face.String()) + c.suit.String()
}

func (c *Card) String() string {
	if c == nil {
		return "[...]"
	}
	return fmt.Sprintf("[%2v%v]", c.face.String(), c.suit.String())
}
