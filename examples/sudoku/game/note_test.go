package game

import "testing"

func TestNoteRemoveOneSeven(t *testing.T) {
	n := NewNote(3)
	got := n.String()
	want := "123\n456\n789\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
	n.RemoveNote(1)
	got = n.String()
	want = ".23\n456\n789\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
	n.RemoveNote(7)
	got = n.String()
	want = ".23\n456\n.89\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestNoteAdd1(t *testing.T) {
	n := NewNote(3)
	n.AddNote(1)
	got := n.String()
	want := "1!23\n456\n789\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
	n.RemoveNote(1)
	got = n.String()
	want = ".23\n456\n789\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
	n.AddNote(1)
	got = n.String()
	want = "123\n456\n789\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestNoteAdd379(t *testing.T) {
	n := NewNote(3)
	n.RemoveNote(3)
	n.RemoveNote(7)
	n.RemoveNote(9)
	got := n.String()
	want := "12.\n456\n.8.\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
	n.AddNote(3)
	n.AddNote(7)
	n.AddNote(9)
	got = n.String()
	want = "123\n456\n789\n"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
