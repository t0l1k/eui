package deck

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

func GetAllCardSuit() []Suit {
	return []Suit{Clubs, Diamonds, Hearts, Spades}
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
