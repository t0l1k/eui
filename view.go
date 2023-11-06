package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Базовый виджет умею хранить состояние от указателя мыши, при наведении(Hover) при нажатие на виджет(Focus) при покадании курсора мыши(Normal)
type View struct {
	name string
	BoxLayout
	state                        InputState
	rect                         *Rect
	image                        *ebiten.Image
	dirty, visible, disabled     bool
	bg, fg                       color.Color
	isDragging                   bool
	dragStartPoint, dragEndPoint PointInt
}

func NewView() *View {
	v := &View{}
	v.SetupView()
	return v
}

func (v *View) SetupView() {
	v.horizontal = true
	theme := GetUi().theme
	v.Bg(theme.Get(ViewBg))
	v.Name("view")
	v.Parent(nil)
	v.SetState(ViewStateNormal)
	v.Visible(true)
	GetUi().inputMouse.Attach(v)
	GetUi().inputTouch.Attach(v)
}

func (v *View) GetImage() *ebiten.Image {
	return v.image
}

func (v *View) Image(image *ebiten.Image) {
	v.image = image
	v.dirty = true
}

func (v *View) GetRect() *Rect {
	return v.rect
}
func (v *View) Rect(rect []int) {
	v.rect = NewRect(rect)
	v.dirty = true
}

func (v *View) IsDisabled() bool {
	return v.disabled
}

func (v *View) Enable() {
	if !v.disabled {
		return
	}
	v.disabled = false
	v.dirty = true
}

func (v *View) Disable() {
	if v.disabled {
		return
	}
	v.disabled = true
	v.state = ViewStateNormal
	v.dirty = true
}

func (v *View) IsDirty() bool {
	return v.dirty
}

func (v *View) Dirty(value bool) {
	v.dirty = value
}

func (v *View) IsVisible() bool {
	return v.visible
}

func (v *View) Visible(value bool) {
	v.visible = value
	v.dirty = true
}

func (v *View) Bg(bg color.Color) {
	if v.bg == bg {
		return
	}
	v.bg = bg
	v.dirty = true
}

func (v *View) Fg(fg color.Color) {
	if v.fg == fg {
		return
	}
	v.fg = fg
	v.dirty = true
}

func (v *View) Name(name string) {
	v.name = name
}

func (v *View) SetState(state InputState) {
	if v.state == state {
		return
	}
	if v.parent != nil {
		vrb, ok := v.parent.(*View)
		if ok {
			vrb.SetState(ViewStateNormal)
		}
	}
	v.state = state
	v.dirty = true
}

func (v *View) Layout() {
	w0, h0 := v.rect.Size()
	if v.image == nil {
		v.image = ebiten.NewImage(w0, h0)
	} else {
		v.image.Clear()
	}
	v.image.Fill(v.bg)
	v.dirty = false
}

func (v *View) UpdateInput(value interface{}) {
	if !v.visible || v.disabled {
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
					v.dragEndPoint = PointInt{x, y}
					v.isDragging = false
				}
			}
			if b == buttonPressed {
				if v.state == ViewStateHover {
					v.SetState(ViewStateFocus)
				}
				if !v.isDragging {
					v.isDragging = true
					v.dragStartPoint = PointInt{x, y}
				}
				v.dragEndPoint = PointInt{x, y}
			}
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
						v.dragStartPoint = PointInt{x, y}
					}
					v.dragEndPoint = PointInt{x, y}
				}
				if b == -1 {
					if v.state == ViewStateFocus {
						v.SetState(ViewStateExec)
					}
					if v.isDragging {
						v.dragEndPoint = PointInt{x, y}
						v.isDragging = false
					}
				}
			} else {
				v.SetState(ViewStateNormal)
			}
		}
	}
}

func (v *View) Update(dt int) {
	if !v.visible || v.disabled {
		return
	}
	for _, c := range v.Container {
		c.Update(dt)
	}
}

func (v *View) Draw(surface *ebiten.Image) {
	if !v.visible {
		return
	}
	if v.dirty {
		v.Layout()
		for _, c := range v.Container {
			c.Layout()
		}
	}
	op := &ebiten.DrawImageOptions{}
	x, y := v.rect.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(v.image, op)
	for _, v := range v.Container {
		v.Draw(surface)
	}
}

func (v *View) Resize(rect []int) {
	v.rect = NewRect(rect)
	v.BoxLayout.Resize(rect)
	v.dirty = true
	v.image = nil
}

func (v *View) Close() {
	for _, c := range v.Container {
		c.Close()
	}
}
