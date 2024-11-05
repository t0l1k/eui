package fc

import (
	"fmt"
	"testing"

	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

func TestFreeFieldShow(t *testing.T) {
	f := NewFreecell(deck.NewDeckCards52())
	got := f.GetDeck()
	want := 60
	if len(got) != want {
		t.Errorf("test getDeck %v %v %v %v", got, want, len(got), f)
	}
}

func TestFreeCellInit(t *testing.T) {
	f := NewFreecell(deck.NewDeckCards52().Shuffle())
	fmt.Println("test init", f)
	got := f.Column(ColHouses[0])[0].IsEmpty()
	want := true
	if got != want {
		t.Error("Test Frecell Game At Init House is not nil ", got, want, f)
	}
	gotL := len(f.Column(sols.Col1))
	wantL := 7
	if gotL != wantL {
		t.Error("Test Frecell Game At Init column 1 len is not 7", f, got)
	}

	gotL = len(f.Column(sols.Col8))
	wantL = 6
	if gotL != wantL {
		t.Error("Test Frecell Game At Init column 8 len is not 6", f, got)
	}
}

func TestFreeCellGetColumnLastCard(t *testing.T) {
	f := NewFreecell(deck.NewDeckCards52())
	col := f.Column(sols.Col1)
	idx, card := f.ColumnLast(sols.Col1)
	card0 := col[idx].GetCard()
	if !card0.Eq(card.GetCard()) {
		t.Error("column 1:", sols.Col1, col, len(col), card0, idx, card)
	}

	col = f.Column(sols.ColFC1)
	idx, card = f.ColumnLast(sols.ColFC1)
	card0 = col[idx].GetCard()
	if !card0.Eq(card.GetCard()) {
		t.Error("column fc1:", sols.ColFC1, col, len(col), card0, idx, card.GetCard(), card0.Eq(card.GetCard()))
	}

}

func TestFreeCellMakeMove(t *testing.T) {
	t.Run("Ход на столб 1", func(t *testing.T) {
		f := NewFreecell(deck.NewDeckCards52())
		f.MakeMove(sols.Col1)
		fmt.Println("test move at init field is\n", f)
		gotT := f.MakeMove(sols.Col1)
		fmt.Println("test move after 1 move field is\n", f)
		got := f.Column(sols.ColFC1)
		if !gotT || got == nil && len(got) == 9 {
			t.Error("Test Frecell Game At Make Move Free cell 1 is not nil ", f, got, got == nil)
		}
	})

	t.Run("Ход на столб 1 1", func(t *testing.T) {
		f := NewFreecell(deck.NewDeckCards52())
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		fmt.Println("test move after 2 move field is\n", f)
		gotT := f.Column(sols.ColFc2)[0].IsEmpty()
		wantT := false
		if gotT == wantT {
			t.Error("Test Frecell Game At Make Move Free cell 1 is not nil ", f, gotT, wantT)
		}
	})

	t.Run("Ход на столб 1 1 1", func(t *testing.T) {
		f := NewFreecell(deck.NewDeckCards52())
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		fmt.Println("test move after 2 move field is\n", f)
		gotT := f.Column(sols.ColFc2)[0].IsEmpty()
		wantT := false
		if gotT != wantT {
			t.Error("Test Frecell Game At Make Move Free cell 1 is not nil ", f, gotT, wantT)
		}
	})

	t.Run("Ход на столб 1 1 1 1 1 1", func(t *testing.T) {
		f := NewFreecell(deck.NewDeckCards52())
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		gotT := f.Column(sols.ColFc2)[0].IsEmpty()
		wantT := false
		if gotT != wantT {
			t.Error("Test Frecell Game At Make Move Free cell 1 is not nil ", f, gotT, wantT)
		}
	})

	t.Run("Ход на столб 1 1 8", func(t *testing.T) {
		f := NewFreecell(deck.NewDeckCards52())
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col1)
		f.MakeMove(sols.Col8)
		fmt.Println("test move after 3 move field is\n", f)
		_, gotC := f.ColumnLast(sols.Col2)
		carfd := deck.NewCard(deck.Spades, deck.Queen)
		wantC := sols.NewCell().SetCard(carfd)
		fmt.Printf("after 3 move gotC:%v wantC:%v card:%v\n%v", gotC, wantC, carfd, f)
		if gotC == wantC {
			t.Error("Test Frecell Game At Make Move Free cell 1 is not nil ", f)
		}
	})
}

func TestFaceLessOneFalse(t *testing.T) {
	f := NewFreecell(deck.NewDeckCards52())
	dataTrue := []sols.Column{
		sols.Col1,
		sols.Col2,
		sols.Col7,
		sols.Col8,
		sols.Col1,
		sols.Col2,
		sols.Col7,
		sols.Col8,
		sols.Col1,
		sols.Col2,
		sols.Col7,
		sols.Col8,
		sols.Col1,
		sols.Col2,
		sols.Col7,
		sols.Col8,
		sols.Col1,
		sols.Col2,
		sols.Col7,
		sols.Col8,
		sols.Col1,
		// sols.Col2,
		// sols.Col7,
		// sols.Col8,
		// sols.Col1,
		// sols.Col2,
		// sols.Col7,
		// sols.Col8,
	}
	for i, v := range dataTrue {
		got := f.MakeMove(v)
		want := false
		if got != want {
			_, cell := f.ColumnLast(v)
			t.Errorf("#%v %v value:%v got %v want %v \n%v", i, v, cell, got, want, f)
		}
	}
}
