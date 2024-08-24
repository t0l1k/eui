package utils

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestList(t *testing.T) {
	t.Run("List new test max 12345", func(t *testing.T) {
		l := NewIntList()
		for i := -5; i <= 5; i++ {
			n := rand.Intn(6)
			l = l.Add(n)
		}
		got := l.Max()
		want := 5
		if got != want {
			t.Errorf("got %v want %v from %v", got, want, l)
		}
	})

	t.Run("List new test min 54321", func(t *testing.T) {
		l := NewIntList()
		for i := 5; i >= -5; i-- {
			l = l.Add(i)
		}
		got := l.Min()
		want := -5
		if got != want {
			t.Errorf("got %v want %v from %v", got, want, l)
		}
	})

	t.Run("List Add", func(t *testing.T) {
		l := NewIntList()
		l = l.Add(1)
		l = l.Add(2)
		fmt.Println("add 1", l)
		got := l.IsContain(1)
		want := true
		if got != want {
			t.Errorf("got %v want %v from %v", got, want, l)
		}
	})

	t.Run("List Remove", func(t *testing.T) {
		l := NewIntList().Add(1).Add(2).Add(3)
		fmt.Println("add 1 2 3", l)
		l, _ = l.Remove(3)
		fmt.Println("remove 3 from  (1 2)", l)
		got := l.IsContain(3)
		want := false
		if got != want {
			t.Errorf("got %v want %v from %v", got, want, l)
		}
	})
	t.Run("List Pop", func(t *testing.T) {
		l := NewIntList().Add(1).Add(2).Add(3)
		fmt.Println("add 1 2 3", l)
		l = l.Pop()
		fmt.Println("pop from  (1 2)", l)
		got := l.IsContain(3)
		want := false
		if got != want {
			t.Errorf("got %v want %v from %v", got, want, l)
		}
	})
}
