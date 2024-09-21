package sols

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

// Прввила пасьянса, есть поле в 15 рядов по 4 карты в ряде(2 ряда пустые на старте), каждый ряд содержит 4 карты одного достоинства, когда собран(все тузы, короли и т.д. в одном ряду). Перенос на пустой ряд любой карты, или в ряд где менее 4 карт, можно положить только того же достоинства карту.
type Lauout15 map[int][]*Cell

func NewLayout15(deck *deck.DeckCards52) *Lauout15 {
	l := make(Lauout15)
	l.Reset(deck)
	return &l
}

func (l Lauout15) Reset(deck52 *deck.DeckCards52) {
	idx := 0
	for i := 0; i < 15; i++ {
		l[i] = nil
		for p := 0; p < 4; p++ {
			if idx >= deck52.Len() {
				l.AddCard(i, nil)
				continue
			}
			l.AddCard(i, deck52.Deck52()[idx])
			idx++
		}
	}
}

func (l Lauout15) GetDeck() (deck []*deck.Card) {
	for i := 0; i < 15; i++ {
		for _, cell := range l.Row(i) {
			deck = append(deck, cell.GetCard())
		}
	}
	return deck
}

func (l Lauout15) SetDeck(deck []*deck.Card) {
	idx := 0
	for i := 0; i < 15; i++ {
		for p := 0; p < 4; p++ {
			l[i][p].SetCard(deck[idx])
			idx++
		}

	}
}

func (l Lauout15) AddCard(row int, card *deck.Card) {
	cell := NewCell()
	cell.SetCard(card)
	l[row] = append(l[row], cell)
}

func (l Lauout15) RemoveLastCard(row int) {
	isEmpty, idx := l.IsRowEmpty(row)
	if !isEmpty {
		l[row][idx].SetCard(nil)
	}
}

func (l Lauout15) IsRowEmpty(row int) (bool, int) {
	for i := len(l[row]) - 1; i >= 0; i-- {
		if l[row][i].GetCard() != nil {
			return false, i
		}
	}
	return true, 0
}

func (l Lauout15) SortRows() []int {
	keys := make([]int, 0)
	for k := range l {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func (l Lauout15) Row(row int) []*Cell {
	return l[row]
}

func (l Lauout15) RowLastCard(row int) *deck.Card {
	isEmpty, idx := l.IsRowEmpty(row)
	if !isEmpty {
		return l[row][idx].GetCard()
	}
	return nil
}

func (l Lauout15) InWhichRow(card0 *deck.Card) int {
	for i, cells := range l {
		for _, cell := range cells {
			if cell.GetCard().Eq(card0) {
				return i
			}
		}
	}
	return -1
}

func (l Lauout15) IsSolved() bool {
	empty, solved := 0, 0
	for i := range l.SortRows() {
		isEmpty, idx := l.IsRowEmpty(i)
		if idx < 3 && !isEmpty {
			return false
		}
		if isEmpty {
			empty++
		}
		if idx == 3 {
			card0 := l[i][0]
			for j := 1; j < 4; j++ {
				card := l[i][j].GetCard()
				if !card0.GetCard().EqFace(card) {
					return false
				}
				solved++
			}
			solved++
		}
	}
	fmt.Println("Win!!! solved:", solved, empty)
	return solved == 52 && empty == 2
}

func (l Lauout15) AvailableMoves() (count int) {
	for row := range l.SortRows() {
		card0 := l.RowLastCard(row)
		for i := range l.SortRows() {
			if i == row {
				continue
			}
			card := l.RowLastCard(i)
			isEmpty, idx := l.IsRowEmpty(i)
			if idx > 2 && !isEmpty {
				continue
			} else if idx <= 2 && !isEmpty {
				if card0.EqFace(card) {
					for _, cell := range l[i] {
						if cell.GetCard() == nil {
							count++
						}
					}
				}
			}
		}
		for i := range l.SortRows() {
			if i == row {
				continue
			}
			isEmpty, _ := l.IsRowEmpty(i)
			if isEmpty {
				for _, cell := range l[i] {
					if cell.GetCard() == nil {
						count++
					}
				}
			}
		}
	}
	return count
}

func (l Lauout15) MakeMove(row int) bool {
	card0 := l.RowLastCard(row)
	for i := range l.SortRows() {
		if i == row {
			continue
		}
		card := l.RowLastCard(i)
		isEmpty, idx := l.IsRowEmpty(i)
		if idx > 2 && !isEmpty {
			continue
		} else if idx <= 2 && !isEmpty {
			if card0.EqFace(card) {
				for idx, v := range l[i] {
					if v.GetCard() == nil {
						l[i][idx].SetCard(card0)
						l.RemoveLastCard(row)
						return true
					}
				}
			}
		}
	}
	for i := range l.SortRows() {
		if i == row {
			continue
		}
		isEmpty, _ := l.IsRowEmpty(i)
		if isEmpty {
			for idx, v := range l[i] {
				if v.GetCard() == nil {
					l[i][idx].SetCard(card0)
					l.RemoveLastCard(row)
					return true
				}
			}
		}
	}
	return false
}

func (l Lauout15) String() string {
	str := "Sol15\n"
	for _, key := range l.SortRows() {
		cells := l[key]
		str += strconv.Itoa(key) + "{"
		for _, cell := range cells {
			if cell.GetCard() != nil {
				str += cell.GetCard().String()
			} else {
				str += fmt.Sprintf("[%2v]", "...")
			}
		}
		if (key+1)%3 == 0 && key > 0 {
			str += "\n"
		} else {
			str += "} "
		}
	}
	return str
}
