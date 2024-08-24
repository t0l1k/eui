package pos

import (
	"fmt"
	"testing"

	"github.com/t0l1k/eui/utils"
)

func TestPosition(t *testing.T) {
	t.Run("Pos new test turn", func(t *testing.T) {
		pos := NewPosititon(3)
		got := pos.turn
		want := 'x'
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("Pos new test String", func(t *testing.T) {
		pos := NewPosititon(3)
		got := pos.StringShort()
		want := "         "
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("Pos test move", func(t *testing.T) {
		pos := NewPosititon(3)
		pos.Move(1)
		got := TurnO
		want := 'o'
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
		gotS := pos.StringShort()
		wantS := " x       "
		if gotS != wantS {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("Pos test unmove", func(t *testing.T) {
		pos := NewPosititon(3)
		pos.UnMove(1)
		got := TurnX
		want := 'x'
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
		gotS := pos.StringShort()
		wantS := "         "
		if gotS != wantS {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("Pos test possible moves", func(t *testing.T) {
		pos := NewPosititon(3)
		var arr utils.IntList
		for i := 0; i < pos.size; i++ {
			arr = arr.Add(i)
		}
		fmt.Println("arr:", arr)
		arr, err := arr.Remove(1)
		if err != nil {
			panic(err)
		}
		fmt.Println("arr1:", arr)
		arr, err = arr.Remove(2)
		if err != nil {
			panic(err)
		}
		fmt.Println("arr2:", arr)

		pos.Move(1)
		pos.Move(2)
		got := pos.PossibleMoves()
		want := arr
		if !arr.Equals(got) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("Pos test is win for", func(t *testing.T) {
		got := NewPosititon(3).IsWinFor(TurnX)
		want := false
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		pos := NewPosititon(3)
		pos.Setup(TurnX, []rune{'x', 'x', 'x', TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty})
		got = pos.IsWinFor(TurnX)
		want = true
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos = NewPosititon(3)
		pos = pos.Setup(TurnX, []rune{TurnX, TurnEmpty, TurnEmpty, TurnX, TurnEmpty, TurnEmpty, TurnX, TurnEmpty, TurnEmpty})
		got = pos.IsWinFor(TurnX)
		want = true
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos = NewPosititon(3)
		pos = pos.Setup(TurnX, []rune{TurnX, TurnEmpty, TurnEmpty, TurnEmpty, TurnX, TurnEmpty, TurnEmpty, TurnEmpty, TurnX})
		got = pos.IsWinFor(TurnX)
		want = true
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos = NewPosititon(3)
		pos = pos.Setup(TurnX, []rune{TurnEmpty, TurnEmpty, TurnO, TurnEmpty, TurnO, TurnEmpty, TurnO, TurnEmpty, TurnEmpty})
		got = pos.IsWinFor(TurnO)
		want = true
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}
	})

	t.Run("Pos test minimax", func(t *testing.T) {
		pos := NewPosititon(3)

		got := pos.Minimax()
		want := 0
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos.Setup(TurnX, []rune{TurnX, TurnX, TurnX, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty})
		got = pos.Minimax()
		want = 6
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos.Setup(TurnO, []rune{TurnO, TurnO, TurnO, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty})
		got = pos.Minimax()
		want = -6
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos.Setup(TurnO, []rune{TurnO, TurnX, TurnO, TurnO, TurnX, TurnO, TurnX, TurnO, TurnX})
		got = pos.Minimax()
		want = 0
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos.Setup(TurnX, []rune{TurnX, TurnX, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty})
		got = pos.Minimax()
		want = 6
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos.Setup(TurnO, []rune{TurnO, TurnO, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty})
		got = pos.Minimax()
		want = -6
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

	})

	t.Run("Pos test bestmove", func(t *testing.T) {
		pos := NewPosititon(3)
		pos.Setup(TurnX, []rune{TurnX, TurnX, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty})
		got := pos.BestMove()
		want := 2
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos.Setup(TurnO, []rune{TurnO, TurnO, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty})
		got = pos.BestMove()
		want = 2
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}
	})

	t.Run("Pos test game end", func(t *testing.T) {
		pos := NewPosititon(3)
		got := pos.IsGameEnd()
		want := false
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos.Setup(TurnX, []rune{TurnX, TurnX, TurnX, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty, TurnEmpty})
		got = pos.IsGameEnd()
		want = true
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

		pos.Setup(TurnO, []rune{TurnO, TurnX, TurnO, TurnO, TurnX, TurnO, TurnX, TurnO, TurnX})
		got = pos.IsGameEnd()
		want = true
		if got != want {
			t.Errorf("got %v want %v %v", got, want, pos)
		}

	})

}
