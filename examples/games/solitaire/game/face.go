package game

type Face int

const (
	Ace Face = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func GetAllCardFace() []Face {
	return []Face{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}
}

func (f Face) String() string {
	s := ""
	switch f {
	case Ace:
		s = "A"
	case Two:
		s = "2"
	case Three:
		s = "3"
	case Four:
		s = "4"
	case Five:
		s = "5"
	case Six:
		s = "6"
	case Seven:
		s = "7"
	case Eight:
		s = "8"
	case Nine:
		s = "9"
	case Ten:
		s = "10"
	case Jack:
		s = "J"
	case Queen:
		s = "Q"
	case King:
		s = "K"
	}
	return s
}
