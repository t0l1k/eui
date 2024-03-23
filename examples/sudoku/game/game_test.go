package game

import (
	"testing"

	"github.com/t0l1k/eui"
)

func TestCellAdd3(t *testing.T) {
	c := NewCell(3, eui.NewPointInt(3, 3))
	c.Add(3)
	got := c.String()
	want := "\n  3\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCellAddNote3(t *testing.T) {
	c := NewCell(2, eui.NewPointInt(3, 0))
	c.AddNote(3)
	got := c.String()
	want := "12\n3!4\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFieldPrintShort(t *testing.T) {
	c := NewField(2)
	c.cells[0].SetValue(1)
	got := c.String()
	want := "sudoku 4X4\n  1.........\n............\n............\n............\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFieldIsValidHorizontal(t *testing.T) {
	c := NewField(2)
	c.cells[0].SetValue(1)
	gotB := c.isValidHor(0, 2)
	wantB := true
	if gotB != wantB {
		t.Errorf("isValiHor(0,2) got %v want %v", gotB, wantB)
	}

	gotB = c.isValidHor(0, 1)
	wantB = false
	if gotB != wantB {
		t.Errorf("isValiHor(0,1) got %v want %v", gotB, wantB)
	}

	c.cells[1].SetValue(2)
	got := c.String()
	want := "sudoku 4X4\n  1  2......\n............\n............\n............\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

}

func TestFieldIsValidVertical(t *testing.T) {
	c := NewField(2)
	c.cells[0].SetValue(1)
	gotB := c.isValidVer(0, 2)
	wantB := true
	if gotB != wantB {
		t.Errorf("isValiVer(0,2) got %v want %v", gotB, wantB)
	}

	gotB = c.isValidVer(0, 1)
	wantB = false
	if gotB != wantB {
		t.Errorf("isValiVer(0,1) got %v want %v", gotB, wantB)
	}

	gotB = c.isValidVer(0, 3)
	wantB = true
	if gotB != wantB {
		t.Errorf("isValiVer(0,3) got %v want %v", gotB, wantB)
	}

	c.cells[4].SetValue(2)
	got := c.String()
	want := "sudoku 4X4\n  1.........\n  2.........\n............\n............\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFieldIsValidHorAndVer(t *testing.T) {
	c := NewField(2)
	c.cells[0].SetValue(1)
	gotB := c.IsValidMove(1, 0, 2)
	wantB := true
	if gotB != wantB {
		t.Errorf("isValid(1,0,2) got %v want %v", gotB, wantB)
	}

	gotB = c.IsValidMove(1, 0, 1)
	wantB = false
	if gotB != wantB {
		t.Errorf("isValid(1,0,1) got %v want %v", gotB, wantB)
	}

	gotB = c.IsValidMove(1, 0, 3)
	wantB = true
	if gotB != wantB {
		t.Errorf("isValid(1,0,3) got %v want %v", gotB, wantB)
	}

	c.cells[1].SetValue(2)
	c.cells[4].SetValue(2)
	got := c.String()
	want := "sudoku 4X4\n  1  2......\n  2.........\n............\n............\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFieldIsValidRect(t *testing.T) {
	c := NewField(2)
	c.cells[0].SetValue(1)
	c.cells[1].SetValue(2)

	gotB := c.IsValidMove(0, 1, 3)
	wantB := true
	if gotB != wantB {
		t.Errorf("isValidRect(0,1,3) got %v want %v", gotB, wantB)
	}

	c.cells[2].SetValue(3)
	c.cells[3].SetValue(4)
	c.cells[4].SetValue(3)

	gotB = c.isValidRect(1, 1, 4)
	wantB = true
	if gotB != wantB {
		t.Errorf("isValidRect(1,1,4) got %v want %v", gotB, wantB)
	}

	gotB = c.IsValidMove(1, 1, 1)
	wantB = false
	if gotB != wantB {
		t.Errorf("isValidRect(1,1,1) got %v want %v", gotB, wantB)
	}

	c.cells[5].SetValue(4)
	c.cells[6].SetValue(1)
	c.cells[7].SetValue(2)
	c.cells[8].SetValue(2)
	c.cells[9].SetValue(3)

	gotB = c.IsValidMove(2, 2, 4)
	wantB = true
	if gotB != wantB {
		t.Errorf("isValidRect(2,2,4) got %v want %v", gotB, wantB)
	}

	c.cells[10].SetValue(4)
	c.cells[11].SetValue(1)
	c.cells[12].SetValue(4)
	c.cells[13].SetValue(1)

	gotB = c.IsValidMove(2, 3, 2)
	wantB = true
	if gotB != wantB {
		t.Errorf("isValidRect(2,3,2) got %v want %v", gotB, wantB)
	}

	gotB = c.IsValidMove(2, 3, 4)
	wantB = false
	if gotB != wantB {
		t.Errorf("isValidRect(2,3,4) got %v want %v", gotB, wantB)
	}

	got := c.String()
	want := "sudoku 4X4\n  1  2  3  4\n  3  4  1  2\n  2  3  4  1\n  4  1......\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

// func TestFieldNextRand(t *testing.T) {
// go panicOnTimeout(24 * time.Second)
// 	c := NewField(2)
// 	for {
// 		n := c.nextRand()
// 		if n == 0 || n == c.size+1 {
// 			t.Errorf("for 0 or %v got %v ", c.size+1, n)
// 		}
// 	}
// }

// func TestFieldNew(t *testing.T) {
// 	// go panicOnTimeout(1 * time.Second)
// 	c := NewField(2)
// 	c.New()
// }

// func panicOnTimeout(d time.Duration) {
// 	<-time.After(d)
// 	panic("Test timed out")
// }
