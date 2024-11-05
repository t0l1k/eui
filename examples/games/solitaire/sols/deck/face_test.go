package deck

import "testing"

func TestFaceLessOneTrue(t *testing.T) {
	dataTrue := map[Card]Face{
		Card{suit: Clubs, face: Ace}:   Two,
		Card{suit: Clubs, face: Two}:   Three,
		Card{suit: Clubs, face: Three}: Four,
		Card{suit: Clubs, face: Four}:  Five,
		Card{suit: Clubs, face: Five}:  Six,
		Card{suit: Clubs, face: Six}:   Seven,
		Card{suit: Clubs, face: Seven}: Eight,
		Card{suit: Clubs, face: Eight}: Nine,
		Card{suit: Clubs, face: Nine}:  Ten,
		Card{suit: Clubs, face: Ten}:   Jack,
		Card{suit: Clubs, face: Jack}:  Queen,
		Card{suit: Clubs, face: Queen}: King,
	}
	for k, v := range dataTrue {
		got := k.face.IsOneLess(v)
		want := true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}
}

func TestFaceLessOneFalse(t *testing.T) {
	dataFalse := map[Card]Face{
		*NewCard(Clubs, Ace):         Ace,
		*NewCard(Clubs, Ace):         Three,
		Card{suit: Clubs, face: Ace}: Four,
		Card{suit: Clubs, face: Ace}: King,
		Card{suit: Clubs, face: Two}: Ace,
		Card{suit: Clubs, face: Two}: Two,
		Card{suit: Clubs, face: Two}: Three,
		Card{suit: Clubs, face: Two}: Four,
		Card{suit: Clubs, face: Two}: King,
	}
	for k, v := range dataFalse {
		got := k.face.IsOneLess(v)
		want := false
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}
}
