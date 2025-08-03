package eui

// Базовый виджет умею хранить состояние от указателя мыши, при наведении(Hover) при нажатие на виджет(Focus) при покадании курсора мыши(Normal)
type View struct {
	*Drawable
	state                        InputState
	isDragging                   bool
	dragStartPoint, dragEndPoint Point[int]
}

func NewView() *View {
	v := &View{Drawable: NewDrawable()}
	v.SetupView()
	return v
}

func (v *View) SetupView() {
	theme := GetUi().theme
	v.Bg(theme.Get(ViewBg))
	v.SetState(ViewStateNormal)
	GetUi().inputMouse.Attach(v)
	GetUi().inputTouch.Attach(v)
}

func (v *View) Enable() {
	if !v.disabled {
		return
	}
	v.disabled = false
	v.MarkDirty()
}

func (v *View) Disable() {
	if v.disabled {
		return
	}
	v.disabled = true
	v.state = ViewStateNormal
	v.MarkDirty()
}

func (v *View) GetState() InputState {
	return v.state
}

func (v *View) SetState(state InputState) {
	if v.state == state {
		return
	}
	v.state = state
	v.MarkDirty()
}

func (v *View) UpdateInput(value interface{}) {
	if v.disabled {
		return
	}
	switch vl := value.(type) {
	case MouseData:
		x, y, b := vl.position.X, vl.position.Y, vl.button
		inRect := v.rect.InRect(x, y)
		if inRect {
			if b == buttonReleased {
				if v.state == ViewStateNormal {
					v.SetState(ViewStateHover)
				}
				if v.state == ViewStateFocus {
					v.SetState(ViewStateHover)
				}
				if v.isDragging {
					v.dragEndPoint = Point[int]{x, y}
					v.isDragging = false
				}
			}
			if b == buttonPressed {
				if v.state == ViewStateHover {
					v.SetState(ViewStateFocus)
				}
				if !v.isDragging {
					v.isDragging = true
					v.dragStartPoint = Point[int]{x, y}
				}
			}
			v.dragEndPoint = Point[int]{x, y}
		} else if v.state != ViewStateNormal {
			v.SetState(ViewStateNormal)
		}
	case []TouchData:
		for _, vt := range vl {
			x, y, b := vt.pos.X, vt.pos.Y, vt.id
			inRect := v.rect.InRect(x, y)
			if inRect {
				if b >= 0 {
					if v.state == ViewStateNormal {
						v.SetState(ViewStateFocus)
					}
					if !v.isDragging {
						v.isDragging = true
						v.dragStartPoint = Point[int]{x, y}
					}
				}
				v.dragEndPoint = Point[int]{x, y}
				if b == -1 {
					if v.state == ViewStateFocus {
						v.SetState(ViewStateExec)
					}
					if v.isDragging {
						v.dragEndPoint = Point[int]{x, y}
						v.isDragging = false
					}
				}
			} else {
				v.SetState(ViewStateNormal)
			}
		}
	}
}

func (v *View) SetRect(rect Rect[int]) {
	v.Drawable.SetRect(rect)
	v.ImageReset()
}

func (v *View) Close() {
	GetUi().GetInputMouse().Detach(v)
	GetUi().GetInputTouch().Detach(v)
}
