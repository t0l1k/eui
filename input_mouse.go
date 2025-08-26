package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type MouseData struct {
	pos, dragStart            Point[int]
	wPos                      Point[float64]
	isBtnPressed, isDraging   bool
	pressedBtns, releasedBtns map[ebiten.MouseButton]bool
}

func NewMouseData(pos, dragStart Point[int], wPos Point[float64], pressed, released map[ebiten.MouseButton]bool, pressStatus, isDraging bool) MouseData {
	return MouseData{pos: pos, wPos: wPos, dragStart: dragStart, isBtnPressed: pressStatus, isDraging: isDraging, pressedBtns: pressed, releasedBtns: released}
}
func (m *MouseData) Pos() Point[int]                          { return m.pos }
func (m *MouseData) WPos() Point[float64]                     { return m.wPos }
func (m *MouseData) DragStart() Point[int]                    { return m.dragStart }
func (m *MouseData) IsDrag() bool                             { return m.isDraging }
func (m *MouseData) IsPressed(value ebiten.MouseButton) bool  { return m.pressedBtns[value] }
func (m *MouseData) IsReleased(value ebiten.MouseButton) bool { return m.releasedBtns[value] }

func (m *MouseData) EqBtns(other map[ebiten.MouseButton]bool) bool {
	isBtnsEq := true
	for k := range m.pressedBtns {
		if m.pressedBtns[k] != other[k] {
			isBtnsEq = false
		}
	}
	return !isBtnsEq
}
func (m *MouseData) Eq(other MouseData) bool {
	return m.pos.Eq(other.pos) && m.wPos.Eq(other.wPos) && m.dragStart.Eq(other.dragStart) && m.isBtnPressed == other.isBtnPressed && m.isDraging == other.isDraging && !m.EqBtns(other.pressedBtns)
}

type MouseListener struct {
	*Signal[Event]
	LastData MouseData
}

func NewMouseListener(fn func(Event)) *MouseListener {
	pt := NewPoint(-1, -1)
	dt := NewMouseData(pt, pt, NewPoint(0.0, 0.0), make(map[ebiten.MouseButton]bool), make(map[ebiten.MouseButton]bool), false, false)
	m := &MouseListener{
		Signal:   NewSignal[Event](),
		LastData: dt,
	}
	m.Connect(fn)
	return m
}

func (ml *MouseListener) update() {
	currentPos := NewPoint(ebiten.CursorPosition())
	x, y := ebiten.Wheel()
	whellPos := NewPoint(x, y)
	currentPressed := false
	pressedBtns := make(map[ebiten.MouseButton]bool)
	releasedBtns := make(map[ebiten.MouseButton]bool)
	for btn := range ebiten.MouseButtonMax {
		pressed := ebiten.IsMouseButtonPressed(btn)
		pressedBtns[btn] = pressed
		if pressed {
			currentPressed = true
		}
		if ml.LastData.releasedBtns[btn] && !pressed {
			releasedBtns[btn] = true
		}
	}
	isDrag := ml.LastData.isDraging
	ptS := ml.LastData.dragStart
	if currentPressed && !ml.LastData.isBtnPressed && !isDrag {
		isDrag = true
		ptS = currentPos
	}
	newData := NewMouseData(currentPos, ptS, whellPos, pressedBtns, releasedBtns, currentPressed, isDrag)
	// if !ml.Value.Eq(newData) {
	// 	log.Println("MD00:", newData, ml.Value, tmpD)
	// }
	// Если позиция изменилась, посылаем событие движения
	if !currentPos.Eq(ml.LastData.pos) {
		ml.Emit(NewEvent(EventMouseMovement, newData))
		if isDrag {
			ml.Emit(NewEvent(EventMouseDrag, newData))
		}
	}
	// Если кнопка только что нажата, посылаем событие нажатия
	if currentPressed && !ml.LastData.isBtnPressed {
		ml.Emit(NewEvent(EventMouseDown, newData))
	}
	// Если кнопка отпущена, посылаем событие отпускания
	if !currentPressed && ml.LastData.isBtnPressed {
		ml.Emit(NewEvent(EventMouseUp, newData))
	}
	if !whellPos.Eq(ml.LastData.wPos) {
		ml.Emit(NewEvent(EventMouseWheel, newData))
	}
	if !currentPressed && ml.LastData.isBtnPressed && isDrag {
		isDrag = false
		ptS = NewPoint(-1, -1)
		newData = NewMouseData(currentPos, ptS, whellPos, pressedBtns, releasedBtns, currentPressed, isDrag)
	}
	ml.LastData = newData
}
