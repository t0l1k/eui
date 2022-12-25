package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	rect                            *Rect
	Image                           *ebiten.Image
	Dirty, Visible, focus, disabled bool
	bg, fg                          color.Color
	text                            string
	mouseDown                       bool
	onPressed                       func(b *Button)
	left, right, middle             bool
}

func NewButton(text string, rect []int, bg, fg color.Color, f func(b *Button)) *Button {
	return &Button{
		text:      text,
		rect:      NewRect(rect),
		Image:     nil,
		Dirty:     true,
		Visible:   true,
		focus:     false,
		mouseDown: false,
		bg:        bg,
		fg:        fg,
		onPressed: f,
	}
}

func (b *Button) Layout() {
	w, h := b.rect.Size()
	if b.Image == nil {
		b.Image = ebiten.NewImage(w, h)
	} else {
		b.Image.Clear()
	}
	lbl := NewLabel(b.text, []int{0, 0, w, h}, b.bg, b.fg)
	defer lbl.Close()
	if b.disabled || !b.focus && !b.mouseDown {
		lbl.SetBg(b.bg)
		lbl.SetFg(b.fg)
	} else if !b.disabled && b.focus && !b.mouseDown {
		lbl.SetBg(b.fg)
		lbl.SetFg(b.bg)
		lbl.SetRect(false)
	} else if !b.disabled && b.focus && b.mouseDown {
		lbl.SetRect(true)
	}
	lbl.Draw(b.Image)
	b.Dirty = false
}

func (b *Button) GetText() string {
	return b.text
}

func (b *Button) SetText(value string) {
	if b.text == value {
		return
	}
	b.text = value
	b.Dirty = true
}

func (b *Button) SetBg(value color.Color) {
	if b.bg == value {
		return
	}
	b.bg = value
	b.Dirty = true
}

func (b *Button) SetFg(value color.Color) {
	if b.fg == value {
		return
	}
	b.fg = value
	b.Dirty = true
}

func (b *Button) SetFocus(value bool) {
	if b.focus == value {
		return
	}
	b.focus = value
	b.Dirty = true
}

func (b *Button) IsDisabled() bool {
	return b.disabled
}

func (b *Button) Enable() {
	b.disabled = false
	b.Dirty = true
}

func (b *Button) Disable() {
	b.disabled = true
	b.Dirty = true
}

func (b *Button) SetMouseDown(value bool) {
	if b.mouseDown == value {
		return
	}
	b.mouseDown = value
	if b.mouseDown {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.left = true
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			b.right = true
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			b.middle = true
		}
	} else {
		b.left = false
		b.right = false
		b.middle = false
	}
	b.Dirty = true
}

func (b *Button) IsMouseDownLeft() bool   { return b.left }
func (b *Button) IsMouseDownRight() bool  { return b.right }
func (b *Button) IsMouseDownMiddle() bool { return b.middle }

func (b *Button) Update(dt int) {
	if b.disabled {
		return
	}
	x, y := ebiten.CursorPosition()
	if b.rect.InRect(x, y) {
		b.SetFocus(true)
	} else {
		b.SetFocus(false)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		if b.focus {
			b.SetMouseDown(true)
		} else {
			b.SetMouseDown(false)
		}
	} else {
		if b.mouseDown && !b.disabled {
			if b.onPressed != nil {
				b.onPressed(b)
			}
		}
		b.SetMouseDown(false)
	}
}

func (b *Button) Draw(surface *ebiten.Image) {
	if b.Dirty {
		b.Layout()
	}
	if b.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := b.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(b.Image, op)
	}
}

func (b *Button) Resize(rect []int) {
	b.rect = NewRect(rect)
	b.Dirty = true
	b.Image = nil
}

func (b *Button) Close() {
	b.Image.Dispose()
}
