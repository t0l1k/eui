package deck

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type Suit int

const (
	Clubs    Suit = iota // трефы
	Diamonds             // бубны
	Hearts               // червы
	Spades               // пики
)

func GetAllCardSuit() []Suit {
	return []Suit{Clubs, Diamonds, Hearts, Spades}
}

func (s Suit) IsEq(other Suit) bool    { return int(s) == int(other) }
func (s Suit) EqColor(other Suit) bool { return s.Color() == other.Color() }
func (s Suit) Color() color.Color {
	switch s {
	case Clubs, Spades:
		return colornames.Black
	case Diamonds, Hearts:
		return colornames.Red
	}
	return nil
}

func (s Suit) String() string {
	var ch rune
	switch s {
	case Clubs:
		ch = '\u2660'
	case Diamonds:
		ch = '\u2661'
	case Hearts:
		ch = '\u2662'
	case Spades:
		ch = '\u2663'
	}
	return string(ch)
}
