package fc

import (
	"fmt"

	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

type Freecell map[sols.Column][]*sols.Cell

var (
	ColHouses = []sols.Column{sols.ColH1, sols.ColH2, sols.ColH3, sols.ColH4}
	ColFree   = []sols.Column{sols.ColFC1, sols.ColFc2, sols.ColFC3, sols.ColFC4}
	Cols      = []sols.Column{sols.Col1, sols.Col2, sols.Col3, sols.Col4, sols.Col5, sols.Col6, sols.Col7, sols.Col8}
)

func NewFreecell(deck *deck.DeckCards52) *Freecell {
	f := make(Freecell)
	f.Reset(deck)
	return &f
}

func (f Freecell) Reset(deck *deck.DeckCards52) {
	for _, cellName := range ColHouses {
		f[cellName] = nil
		f.addCard(cellName, nil)
	}
	for _, cellName := range ColFree {
		f[cellName] = nil
		f.addCard(cellName, nil)
	}
	idx := 0
	for _, name := range Cols {
		for _, v := range f[name] {
			v.Reset()
		}
		f[name] = nil
	}
	for idx < deck.Len() {
		for _, name := range Cols {
			if idx < deck.Len() {
				card := deck.Deck52()[idx]
				f.addCard(name, card)
				idx++
			}
		}
	}
}

func (f Freecell) GetDeck() (values []*deck.Card) {
	for _, colName := range ColFree {
		cards := f[colName]
		for _, card := range cards {
			values = append(values, card.GetCard())
		}
	}
	for _, colName := range ColHouses {
		cards := f[colName]
		for _, card := range cards {
			values = append(values, card.GetCard())
		}
	}
	for _, colName := range Cols {
		cards := f[colName]
		for _, card := range cards {
			values = append(values, card.GetCard())
		}
	}
	return values
}

func (f Freecell) SetDeck(values []*deck.Card) {
	for _, card := range values {
		for _, v := range ColFree {
			f.addCard(v, card)
		}
		for _, v := range ColHouses {
			f.addCard(v, card)
		}
		for _, v := range Cols {
			f.addCard(v, card)
		}
	}
}

func (l Freecell) Index(card0 *deck.Card) (sols.Column, int) {
	for key, cells := range l {
		for _, cell := range cells {
			if !cell.IsEmpty() && card0.GetCard().Eq(cell.GetCard()) {
				for i, v := range l.Column(key) {
					if v.IsEmpty() {
						continue
					}
					if v.GetCard().Eq(card0) {
						return key, len(cells) - i - 1
					}
				}
			}
		}
	}
	return -1, -1
}

func (f Freecell) AvailableMoves() int {
	return 1
}

func (f Freecell) addCard(column sols.Column, value *deck.Card) {
	cell := sols.NewCell()
	cell.SetCard(value)
	f[column] = append(f[column], cell)
}

func (f Freecell) updateCard(from, to sols.Column, card *deck.Card) {
	var (
		idx0, idx1   int
		card0, card1 *sols.Cell
	)
	idx0, card0 = f.ColumnLast(from)
	if len(f[to]) > 0 {
		idx1, card1 = f.ColumnLast(to)
		fmt.Println("updateCard00:", from, to, card, card0, card1, idx0, idx1, f[from], f[to], f[from][idx0].IsEmpty(), f[to][idx1].IsEmpty())
	} else {
		f.addCard(to, card)
		idx1, card1 = f.ColumnLast(to)
		fmt.Println("updateCard02:", from, to, card, card0, card1, idx0, idx1, f[from], f[from][idx0].IsEmpty())
	}
	cell := f[to][idx1]
	cell.SetCard(card)
	f[from][idx0].Reset()
	f[from] = f[from][:len(f[from])-1]
	fmt.Println("updateCard03:", from, to, card, card0, card1, idx0, idx1, f[from], f[to])
}

func (f Freecell) MakeMove(name0 sols.Column) bool {
	fmt.Printf("00 Ход на %v %v\n", name0, f[name0])
	for _, colName := range Cols {
		fmt.Printf("01 Ход на %v %v\n", name0, f[name0])
		if name0.IsEq(colName) {
			fmt.Printf("02 Ход на %v %v\n", name0, f[name0])
			if f.moveFromCols(name0) {
				fmt.Printf("03 Ход на %v %v\n", name0, f[name0])
				fmt.Println("move done", f)
				return true
			}
		}
	}
	return false
}

func (f Freecell) isValidCardForMove(a, b *deck.Card) bool {
	fmt.Printf("Ход уместен для:{%v %v} \n", a, b)
	return a.GetCard().IsOneLess(b.GetCard()) && !a.GetCard().EqColor(b.GetCard())
}

func (f Freecell) isValidForMoveToFCColumn(name sols.Column) bool {
	for _, v := range ColFree {
		if v.IsEq(name) && f[name][0].IsEmpty() {
			fmt.Printf("Ход уместен:{%v} \n", name)
			return true
		}
	}
	return false
}

func (f Freecell) moveFromCols(name0 sols.Column) bool {
	idx0, card0 := f.ColumnLast(name0)
	fmt.Printf("000 Ход на %v %v %v %v\n", name0, f[name0], idx0, card0)
	for _, colName := range Cols {
		if colName.IsValidMove(name0) {
			idx, card := f.ColumnLast(colName)
			fmt.Printf("001 Ход на {%v %v} %v %v | %v %v %v %v\n", name0, f[name0], idx0, card0, colName, idx, card, f.isValidCardForMove(card0.GetCard(), card.GetCard()))
			if f.isValidCardForMove(card0.GetCard(), card.GetCard()) {
				fmt.Printf("002 Ход на {%v %v} %v %v | %v %v %v %v\n", name0, f[name0], idx0, card0, colName, idx, card, f.isValidCardForMove(card0.GetCard(), card.GetCard()))
				f.updateCard(name0, colName, card0.GetCard())
				// f[name0][idx0].Reset()
				// f[name0] = f[name0][:len(f[name0])-1]
				fmt.Printf("003 Ход на {%v %v} %v %v | %v %v %v %v || %v\n", name0, f[name0], idx0, card0, colName, idx, card, f.isValidCardForMove(card0.GetCard(), card.GetCard()), f[colName])
				return true
			}
		}
	}
	for _, name1 := range ColFree {
		if f.isValidForMoveToFCColumn(name1) {
			idx, card := f.ColumnLast(name1)
			fmt.Printf("004 Ход на {%v %v} %v %v | %v %v %v %v\n", name0, f[name0], idx0, card0, name1, idx, card, f.isValidForMoveToFCColumn(name1))
			// f[name1] = f[name1][:len(f[name1])-1]
			f.updateCard(name0, name1, card0.GetCard())
			// f[name0][idx0].Reset()
			// f[name0] = f[name0][:len(f[name0])-1]
			fmt.Printf("005 Ход на {%v %v} %v %v | %v %v %v %v || %v\n%v", name0, f[name0], idx0, card0, name1, idx, card, f.isValidForMoveToFCColumn(name1), f[name1], f)
			// panic("Ход на свободную ячейку")
			return true
		}
	}
	return false
}

// func (f Freecell) checkMoveFromFC(name0 sols.Column) bool {
// 	for i, name1 := range ColFree {
// 		if name0 == name1 {
// 			fmt.Printf("move from freecell:%v %v %v %v \n", i, name1, f[name1][0].GetCard(), f[name1][0].IsEmpty())

// 			for _, name2 := range Cols {
// 				idx2, card := f.ColumnLast(name2)
// 				fmt.Println("move from freecell to:", idx2, name2)
// 				if card == nil {
// 					fmt.Println("empty column", idx2, name2)
// 					continue
// 				} else if card.GetCard().IsOneHigh(f[name1][0].GetCard()) {
// 					fmt.Println("found high card", f[name0][0].GetCard(), card.GetCard(), card.GetCard().IsOneHigh(f[name1][0].GetCard()))
// 					panic("move on freecell01")
// 				}
// 			}

// 			for _, name2 := range Cols {
// 				idx2, card := f.ColumnLast(name2)
// 				fmt.Println("move from freecell to:", idx2, name2)
// 				if card == nil {
// 					fmt.Println("empty column for move", idx2, name2, f[name0][0])
// 					card0 := f[name0][0].GetCard()
// 					f.addCard(name2, card0)
// 					// f[name2][idx2].SetCard(card0)
// 					f[name0] = nil
// 					return true
// 					// panic("move on freecell02")
// 				}
// 			}

// 		}
// 	}
// 	return false
// }

// func (f Freecell) checkHouse() {
// 	var (
// 		lowerFace deck.Face
// 		searchFor []*deck.Card
// 	)
// 	for _, column := range ColHouses {
// 		if f[column] == nil {
// 			searchFor = append(searchFor, nil)
// 			continue
// 		}
// 		searchFor = append(searchFor, f[column][len(f[column])-1].GetCard())
// 		fmt.Println("check house:", column, f[column], len(f[column]), f[column][0], searchFor, lowerFace)
// 	}
// 	for suit, v := range searchFor {
// 		if v == nil {
// 			lowerFace = deck.Ace
// 			searchFor[suit] = deck.NewCard(deck.Suit(suit), deck.Ace)
// 			fmt.Println("check house02:", suit, v, searchFor, lowerFace)
// 		}
// 		fmt.Println("check house03:", suit, v, searchFor, lowerFace)
// 	}
// 	for idx0, column := range Cols {
// 		idx, card0 := f.ColumnLast(column)
// 		if idx == -1 {
// 			continue
// 		}
// 		for i, card := range searchFor {
// 			fmt.Println("found 04:", i, idx0, idx, ColHouses[i], card, card0.GetCard(), card.Eq(card0.GetCard()))
// 			if card0.GetCard().Eq(card) {
// 				fmt.Println("found 05:", idx0, idx, i, card, card0, card0.GetCard().Eq(card), ColHouses[i], column, f[column][idx])
// 				f.addCard(ColHouses[i], card)
// 				f[column][idx].Reset()
// 				f[column] = f[column][:len(f[column])-1]
// 				fmt.Println("found 06:", idx0, idx, i, card, card0, card0.GetCard().Eq(card), f[column], f[ColHouses[i]])
// 			}
// 		}
// 	}
// }

// func (f Freecell) checkMoveToFC(name0 sols.Column) bool {
// 	for i, name := range ColFree {
// 		fmt.Println("move on free col00:", i, name, name0, f[name], f[name] == nil)
// 		if f[name] == nil {
// 			_, card := f.ColumnLast(name0)
// 			fmt.Println("move on free col01", i, name, name0, f[name], f[name0], card)
// 			f.addCard(name, card.GetCard())
// 			f[name0] = f[name0][:len(f[name0])-1]
// 			fmt.Println("move on free col02", i, name, name0, f[name], f[name0])
// 			f.checkHouse()
// 			return true
// 		}
// 	}
// 	return false
// }

func (f Freecell) Column(name sols.Column) []*sols.Cell { return f[name] }
func (f Freecell) ColumnLast(name sols.Column) (idx int, value *sols.Cell) {
	if len(f[name]) > 0 {
		if cell := f[name][len(f[name])-1]; cell.Value() == nil {
			for i := len(f[name]) - 1; i > 0; i-- {
				if cell := f[name][i]; cell.Value() != nil {
					idx = i
					value = f[name][idx]
					return idx, value
				}
			}
		}
	}
	return len(f[name]) - 1, f[name][len(f[name])-1]
}

func (f Freecell) IsSolved() bool { return false }

func (f Freecell) String() (result string) {
	var (
		idx, row int
	)
	result += "Freecell\n"
	for _, colName := range ColFree {
		if f[colName][row].IsEmpty() {
			result += "[...]"
		} else {
			result += f[colName][row].GetCard().String()
			idx += 1
		}
	}
	result += "|"
	for _, colName := range ColHouses {
		if f[colName][row].IsEmpty() {
			result += "[. .]"
		} else {
			result += f[colName][row].GetCard().String()
			idx += 1
		}
	}
	result += "\n"
	for i := 0; i < 5*8+1; i++ {
		result += "."
	}
	result += "\n"
	for idx < 52 {
		for _, colName := range Cols {
			// fmt.Println("prn00:", colName, idx, row, len(f[colName]), result)
			if len(f[colName]) <= row {
				result += "[ . ]"
				// fmt.Println("prn01:", colName, idx, row, len(f[colName]), result)
			} else if f[colName][row].IsEmpty() {
				result += "[ ..]"
				// fmt.Println("prn01a:", colName, idx, row, len(f[colName]), result)
			} else {
				result += f[colName][row].GetCard().String()
				idx += 1
				// fmt.Println("prn02:", colName, idx, row, len(f[colName]), result)
			}
			if colName == sols.Col8 {
				result += "\n"
				row += 1
				// fmt.Println("prn03:", colName, idx, row, len(f[colName]), result)
			}
			// fmt.Println("prn04:", colName, idx, row, len(f[colName]), result)
		}
	}
	// fmt.Println("prn done", result)
	return result
}
