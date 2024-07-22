package game

import (
	"math/rand"
)

type DeckCards52 struct {
	deck []*Card
}

func NewDeckCards52() *DeckCards52 {
	d := &DeckCards52{}
	d.Reset()
	return d
}

func (d *DeckCards52) Reset() {
	d.deck = nil
	for _, face := range GetAllCardFace() {
		for _, suit := range GetAllCardSuit() {
			d.deck = append(d.deck, NewCard(suit, face))
		}
	}
}

func (d *DeckCards52) Len() int { return len(d.deck) }

func (d *DeckCards52) Shuffle() *DeckCards52 {
	for i := 0; i < 100; i++ {
		nr := rand.Intn(52)
		tmp := d.deck[nr]
		d.deck[nr] = d.deck[0]
		d.deck[0] = tmp
	}
	return d
}

func (d *DeckCards52) String() string {
	str := ""
	for _, card := range d.deck {
		str += card.String()
	}
	return str
}
