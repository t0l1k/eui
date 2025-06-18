package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TouchData struct {
	pos PointInt
	id  ebiten.TouchID
}

func (t TouchData) Equal(a TouchData) bool {
	return t.id == a.id && t.pos.Equal(a.pos)
}

// Умею передать подписчикам события касания экрана это текущее местоположение на экране и сколько нажатий
type TouchInput struct {
	touches  []TouchData
	listener []Inputer
}

func NewTouchInput(fn SlotFunc[Event]) *TouchInput { return &TouchInput{} }

func (s *TouchInput) Attach(o Inputer) {
	s.listener = append(s.listener, o)
}

func (s *TouchInput) Detach(o Inputer) {
	for i, observer := range s.listener {
		if observer == o {
			s.listener = append(s.listener[:i], s.listener[i+1:]...)
			break
		}
	}
}

func (s *TouchInput) Notify() {
	for _, observer := range s.listener {
		observer.UpdateInput(s.touches)
	}
}

func (s *TouchInput) set(value []TouchData) {
	s.touches = value
	s.Notify()
}

func (s *TouchInput) update(_ int) {
	s.touches = s.touches[:0]
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		s.touches = append(s.touches, TouchData{
			id: id,
			pos: PointInt{
				X: x,
				Y: y,
			},
		})
	}
	releasedId := inpututil.AppendJustReleasedTouchIDs(nil)
	if len(releasedId) > 0 {
		s.touches = s.touches[:0]
		for _, id := range releasedId {
			x, y := inpututil.TouchPositionInPreviousTick(id)
			s.touches = append(s.touches, TouchData{
				id: -1,
				pos: PointInt{
					X: x,
					Y: y,
				},
			})
		}
		s.set(s.touches)
	} else {
		if len(s.touches) > 0 {
			s.set(s.touches)
		}
	}
}
