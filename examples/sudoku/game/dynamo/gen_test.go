package dynamo

import (
	"fmt"
	"testing"

	"github.com/t0l1k/eui/examples/sudoku/game"
)

func TestFieldIsValidPathFCheckAllNotes(t *testing.T) {
	arr := []int{4, 2, 1, 3, 1, 3, 2, 4, 2, 1}
	f := NewGenSudokuField(2, game.Easy)
	for i, v := range arr {
		x, y := f.Pos(i)
		if f.IsValidPath(x, y) {
			f.Add(x, y, v)
		}
	}
	fmt.Println("Подготовка теста", f)
	gotB := f.IsValidPath(1, 2)
	wantB := false
	if gotB != wantB {
		t.Errorf("isValidPath(1,2,1) got %v want %v", gotB, wantB)
	}

	fmt.Println("Проверка пустых заметок успешна", f)
	f.ResetCell(1, 2)
	gotNotes := f.getCells()[f.Idx(1, 2)].getNotes()
	wantNotes := []int{1, 4}
	for i, v := range gotNotes {
		if v != wantNotes[i] {
			t.Errorf("ResetCell(1,2) got %v want %v", gotNotes, wantNotes)
		}
	}
	fmt.Println("Проверка обнуления 1,2 ячейки успешна", f)
	f.Add(1, 2, 4)
	gotN := f.getCells()[f.Idx(1, 2)].value
	wantN := 4
	if gotN != wantN {
		t.Errorf("Add(1,2,4) got %v want %v", gotN, wantN)
	}
	fmt.Println("Проверка добавления в ячейку 1,2 4 успешна", f)
	for i, v := range f.getCells() {
		fmt.Println(i, v, v.String())
	}
	gotB = f.IsValidPath(1, 2)
	wantB = true
	if gotB != wantB {
		t.Errorf("isValidPath(1,2,4) got %v want %v", gotB, wantB)
	}
	fmt.Println("Проверка пустого пути от ячейки 1,2 4 успешна", f)
}

func TestFieldIsValidPathF(t *testing.T) {
	arr := []int{4, 2, 1, 3, 1, 3, 2, 4, 2, 1}
	f := NewGenSudokuField(2, game.Easy)
	for i, v := range arr {
		x, y := f.Pos(i)
		if f.IsValidPath(x, y) {
			f.Add(x, y, v)
		}
	}
	gotB := f.IsValidPath(1, 2)
	wantB := false
	if gotB != wantB {
		t.Errorf("isValidPath(1,2,1) got %v want %v", gotB, wantB)
	}
}

func TestFieldIsValidPathT(t *testing.T) {
	arr := []int{4, 2, 1, 3, 1, 3, 2, 4, 2, 4}
	f := NewGenSudokuField(2, game.Easy)
	for i, v := range arr {
		x, y := f.Pos(i)
		f.Add(x, y, v)
	}
	fmt.Println(f.String())
	gotB := f.IsValidPath(1, 2)
	wantB := true
	if gotB != wantB {
		t.Errorf("isValidPath(1,2,4) got %v want %v", gotB, wantB)
	}
}

func TestFieldString(t *testing.T) {
	arr := []int{4, 2, 1, 3, 1, 3, 2, 4, 2, 4, 3, 1, 3, 1, 4, 2}
	f := NewGenSudokuField(2, game.Easy)
	for i, v := range arr {
		x, y := f.Pos(i)
		f.Add(x, y, v)
	}
	got := f.String()
	want := "sudoku 4X4\n[  4][  2][  1][  3]\n[  1][  3][  2][  4]\n[  2][  4][  3][  1]\n[  3][  1][  4][  2]\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestFieldMakeMove(t *testing.T) {
	f := NewGenSudokuField(2, game.Easy)
	x, y := 0, 0
	idx := f.Idx(x, y)
	f.Add(x, y, 4)
	got := f.cells[idx].String()
	want := "[  4]"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
	gotCellNotes := f.getCells()[idx].notes
	wantCellNotes := []int{1, 2, 3}
	for i, v := range wantCellNotes {
		if v != gotCellNotes[i] {
			t.Errorf("got %v want %v", gotCellNotes, wantCellNotes)
		}

	}
}
