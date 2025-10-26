package main

import (
	"log"
	"testing"
)

func TestDim(t *testing.T) {
	t.Run("edge", func(t *testing.T) {
		var tt = []struct {
			str  string
			x, y int
			want bool
		}{
			{"2x2 0,0 false", 0, 0, false},
			{"2x2 0,1 false", 0, 1, false},
			{"2x2 1,1 false", 1, 1, false},
			{"2x2 1,0 false", 1, 0, false},
			{"2x2 2,2 true", 2, 2, true},
			{"2x2 1,3 true", 1, 3, true},
			{"2x2 -1,0 true", -1, 0, true},
			{"2x2 1,5 true", 1, 5, true},
			{"2x2 3,1 true", 3, 1, true},
		}
		d := NewDim(2, 2)
		for _, v := range tt {
			t.Run(v.str, func(t *testing.T) {
				got := d.IsEdge(v.x, v.y)
				want := v.want
				if got != want {
					t.Errorf("%v: got:[%v,%v]%v,%v", v.str, v.x, v.y, got, v.want)
				}
			})
		}
	})
}

func TestFieldInit(t *testing.T) {
	g := NewGame(NewDim(2, 2))
	got := g.field.Cell(3).Value()
	want := 0
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestMove(t *testing.T) {
	t.Run("Move Up", func(t *testing.T) {
		g := NewGame(NewDim(2, 2))
		g.Move(Up)
		log.Println("Up:", g)
		got := g.field.Cell(g.dim.Idx(1, 1)).Value()
		want := 2
		if got != want {
			t.Errorf("got %v want %v \n%v", got, want, g)
		}
	})
}

func TestWin(t *testing.T) {
	t.Run("Win", func(t *testing.T) {
		var tt = []struct {
			str  string
			got  []int
			want bool
		}{
			{"2x2 0123", []int{0, 1, 2, 3}, false},
			{"2x2 1023", []int{1, 0, 2, 3}, false},
			{"2x2 1203", []int{1, 2, 0, 3}, false},
			{"2x2 1230", []int{1, 2, 3, 0}, true},
		}
		for _, v := range tt {
			t.Run(v.str, func(t *testing.T) {
				g := NewGame(NewDim(2, 2))
				for i, cell := range g.field {
					cell.Move(v.got[i])
				}
				got := g.IsWin()
				want := v.want
				if got != want {
					t.Errorf("%v: got:%v want:%v", v.str, got, v.want)
				}
			})
		}
	})
}
