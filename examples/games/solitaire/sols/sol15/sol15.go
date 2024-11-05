package sol15

import (
	"fmt"
	"log"
	"strconv"

	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

// Прввила пасьянса, есть поле в 15 рядов по 4 карты в ряде(2 ряда пустые на старте), каждый ряд содержит 4 карты одного достоинства, когда собран(все тузы, короли и т.д. в одном ряду). Перенос на пустой ряд любой карты, или в ряд где менее 4 карт, можно положить только того же достоинства карту.
type Lauout15 map[sols.Column][]*sols.Cell

func NewLayout15(deck *deck.DeckCards52) *Lauout15 {
	l := make(Lauout15)
	l.Reset(deck)
	return &l
}

func (l Lauout15) Reset(deck52 *deck.DeckCards52) {
	idx := 0
	var i sols.Column
	for i = 0; i < 15; i++ {
		for _, cell := range l[i] {
			cell.Reset()
		}
		l[i] = nil
		for p := 0; p < 4; p++ {
			if idx >= deck52.Len() {
				l.addCard(i, nil)
				continue
			}
			l.addCard(i, deck52.Deck52()[idx])
			idx++
		}
	}
}

func (l Lauout15) GetDeck() (deck []*deck.Card) {
	for i := 0; i < 15; i++ {
		for _, cell := range l.Column(sols.Column(i)) {
			deck = append(deck, cell.GetCard())
		}
	}
	return deck
}

func (l Lauout15) SetDeck(deck []*deck.Card) {
	idx := 0
	var i sols.Column
	for i = 0; i < 15; i++ {
		for p := 0; p < 4; p++ {
			l[i][p].SetCard(deck[idx])
			idx++
		}

	}
}

func (l Lauout15) addCard(row sols.Column, card *deck.Card) {
	cell := sols.NewCell()
	cell.SetCard(card)
	l[row] = append(l[row], cell)
}

func (l Lauout15) removeColumnLastCard(row sols.Column) {
	isEmpty, idx := l.isColumnEmpty(row)
	if !isEmpty {
		l[row][idx].SetCard(nil)
	}
}

func (l Lauout15) isColumnEmpty(row sols.Column) (bool, int) {
	for i := len(l[row]) - 1; i >= 0; i-- {
		if !l[row][i].IsEmpty() {
			return false, i
		}
	}
	return true, 0
}

func (l Lauout15) sortColumns() []sols.Column {
	keys := make([]sols.Column, 0)
	for i := 0; i < 15; i++ {
		keys = append(keys, sols.Column(i))
	}
	return keys
}

func (l Lauout15) Column(row sols.Column) []*sols.Cell {
	return l[row]
}

func (l Lauout15) columnLastCard(row sols.Column) *deck.Card {
	isEmpty, idx := l.isColumnEmpty(row)
	if !isEmpty {
		return l[row][idx].GetCard()
	}
	return nil
}

func (l Lauout15) Index(card0 *deck.Card) (sols.Column, int) {
	var idx int
	for key, cells := range l {
		for _, cell := range cells {
			if !cell.IsEmpty() && card0.GetCard().Eq(cell.GetCard()) {
				for i, v := range l.Column(key) {
					if v.IsEmpty() {
						continue
					}
					if v.GetCard().Eq(card0) {
						idx = 3 - i // supermove
						break
					}
				}
				return key, idx
			}
		}
	}
	return -1, -1
}

func (l Lauout15) IsSolved() bool {
	empty, solved := 0, 0
	for _, i := range l.sortColumns() {
		isEmpty, idx := l.isColumnEmpty(i)
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
				if !card0.GetCard().IsEqFace(card) {
					return false
				}
				solved++
			}
			solved++
		}
	}
	log.Println("Win!!! solved:", solved, empty)
	return solved == 52 && empty == 2
}

func (l Lauout15) AvailableMoves() (count int) {
	for _, row := range l.sortColumns() {
		card0 := l.columnLastCard(row)
		for _, i := range l.sortColumns() {
			if i == row {
				continue
			}
			card := l.columnLastCard(i)
			isEmpty, idx := l.isColumnEmpty(i)
			if idx > 2 && !isEmpty {
				continue
			} else if idx <= 2 && !isEmpty {
				if card0.IsEqFace(card) {
					for _, cell := range l[i] {
						if cell.IsEmpty() {
							count++
						}
					}
				}
			}
		}
		for _, i := range l.sortColumns() {
			if i == row {
				continue
			}
			isEmpty, _ := l.isColumnEmpty(i)
			if isEmpty {
				for _, cell := range l[i] {
					if cell.IsEmpty() {
						count++
					}
				}
			}
		}
	}
	return count
}

// Сделать Ход это сначала проверка всех рядов в которых есть такая же карта и в первую свободную ячейку совершается ход, иначе в пустой ряд ход, иначе пропуск хода
func (l Lauout15) MakeMove(row sols.Column) bool {
	card0 := l.columnLastCard(row)
	for _, i := range l.sortColumns() {
		if i == row {
			continue
		}
		card := l.columnLastCard(i)
		isEmpty, idx := l.isColumnEmpty(i)
		if idx > 2 && !isEmpty {
			continue
		} else if idx <= 2 && !isEmpty {
			if card0.IsEqFace(card) {
				for idx, v := range l[i] {
					if v.IsEmpty() {
						l[i][idx].SetCard(card0)
						l.removeColumnLastCard(row)
						return true
					}
				}
			}
		}
	}
	for _, i := range l.sortColumns() {
		if i == row {
			continue
		}
		isEmpty, _ := l.isColumnEmpty(i)
		if isEmpty {
			for idx, v := range l[i] {
				if v.IsEmpty() {
					l[i][idx].SetCard(card0)
					l.removeColumnLastCard(row)
					return true
				}
			}
		}
	}
	return false
}

func (l Lauout15) String() string {
	str := "Sol15\n"
	for _, key := range l.sortColumns() {
		cells := l[key]
		str += strconv.Itoa(int(key)) + "{"
		for _, cell := range cells {
			if !cell.IsEmpty() {
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
