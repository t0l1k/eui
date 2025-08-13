package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type MouseData struct {
	pos, dragStart          Point[int]
	wPos                    Point[float64]
	isBtnPressed, isDraging bool
	buttons                 map[ebiten.MouseButton]bool
}

func NewMouseData(pos, dragStart Point[int], wPos Point[float64], buttons map[ebiten.MouseButton]bool, pressStatus, isDraging bool) MouseData {
	return MouseData{pos: pos, wPos: wPos, dragStart: dragStart, isBtnPressed: pressStatus, isDraging: isDraging, buttons: buttons}
}
func (m *MouseData) Pos() Point[int]                         { return m.pos }
func (m *MouseData) WPos() Point[float64]                    { return m.wPos }
func (m *MouseData) DragStart() Point[int]                   { return m.dragStart }
func (m *MouseData) IsDrag() bool                            { return m.isDraging }
func (m *MouseData) IsPressed(value ebiten.MouseButton) bool { return m.buttons[value] }

func (m *MouseData) EqBtns(other map[ebiten.MouseButton]bool) bool {
	isBtnsEq := true
	for k := range m.buttons {
		if m.buttons[k] != other[k] {
			isBtnsEq = false
		}
	}
	return !isBtnsEq
}
func (m *MouseData) Eq(other MouseData) bool {
	return m.pos.Eq(other.pos) && m.wPos.Eq(other.wPos) && m.dragStart.Eq(other.dragStart) && m.isBtnPressed == other.isBtnPressed && m.isDraging == other.isDraging && !m.EqBtns(other.buttons)
}

type MouseListener struct {
	*Signal[Event]
	LastData MouseData
}

func NewMouseListener(fn func(Event)) *MouseListener {
	pt := NewPoint(-1, -1)
	dt := NewMouseData(pt, pt, NewPoint(0.0, 0.0), make(map[ebiten.MouseButton]bool), false, false)
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
	btns := make(map[ebiten.MouseButton]bool)
	for btn := range ebiten.MouseButtonMax {
		status := ebiten.IsMouseButtonPressed(btn)
		btns[btn] = status
		if status {
			currentPressed = true
		}
	}
	isDrag := ml.LastData.isDraging
	// tmpD := isDrag
	ptS := ml.LastData.dragStart
	if currentPressed && !ml.LastData.isBtnPressed && !isDrag {
		isDrag = true
		ptS = currentPos
	}
	newData := NewMouseData(currentPos, ptS, whellPos, btns, currentPressed, isDrag)
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
		newData = NewMouseData(currentPos, ptS, whellPos, btns, currentPressed, isDrag)
	}
	ml.LastData = newData
}
