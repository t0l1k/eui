package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ButtonIcon struct {
	rect                            *Rect
	icons                           []*ebiten.Image
	Image                           *ebiten.Image
	Dirty, Visible, focus, disabled bool
	mouseDown                       bool
	onPressed                       func(b *ButtonIcon)
	left, right, middle             bool
}

func NewButtonIcon(icons []*ebiten.Image, rect []int, f func(b *ButtonIcon)) *ButtonIcon {
	return &ButtonIcon{
		icons:     icons,
		rect:      NewRect(rect),
		Image:     nil,
		Dirty:     true,
		Visible:   true,
		focus:     false,
		mouseDown: false,
		disabled:  false,
		onPressed: f,
	}
}

func (b *ButtonIcon) Layout() {
	w, h := b.rect.Size()
	if b.Image == nil {
		b.Image = ebiten.NewImage(w, h)
	} else {
		b.Image.Clear()
	}
	icon := NewIcon(b.icons[0], []int{0, 0, w, h})

	if !b.disabled && !b.focus && !b.mouseDown { //default state
		icon.Draw(b.Image)
		ebitenutil.DrawRect(b.Image, 0, 0, float64(w), float64(h), color.RGBA{32, 32, 32, 32})
	} else if b.disabled { //disabled state
		icon.Draw(b.Image)
		ebitenutil.DrawRect(b.Image, 0, 0, float64(w), float64(h), color.RGBA{32, 32, 32, 64})
	} else if !b.disabled && b.focus && !b.mouseDown { //hover state
		icon.Draw(b.Image)
		ebitenutil.DrawRect(b.Image, 0, 0, float64(w), float64(h), color.RGBA{0, 32, 0, 32})
	} else if !b.disabled && b.focus && b.mouseDown { //pressed state
		icon = NewIcon(b.icons[1], []int{0, 0, w, h})
		icon.Draw(b.Image)
		ebitenutil.DrawRect(b.Image, 0, 0, float64(w), float64(h), color.RGBA{32, 0, 0, 64})
	}
	b.Dirty = false
}

func (b *ButtonIcon) SetIconRelesed(icon *ebiten.Image) {
	if b.icons[0] == icon {
		return
	}
	b.icons[0] = icon
	b.Dirty = true
}

func (b *ButtonIcon) SetIconPressed(icon *ebiten.Image) {
	if b.icons[1] == icon {
		return
	}
	b.icons[1] = icon
	b.Dirty = true
}

func (b *ButtonIcon) SetFocus(value bool) {
	if b.focus == value {
		return
	}
	b.focus = value
	b.Dirty = true
}

func (b *ButtonIcon) IsDisabled() bool {
	return b.disabled
}

func (b *ButtonIcon) Enable() {
	b.disabled = false
	b.Dirty = true
}

func (b *ButtonIcon) Disable() {
	b.disabled = true
	b.Dirty = true
}

func (b *ButtonIcon) SetMouseDown(value bool) {
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

func (b *ButtonIcon) IsMouseDownLeft() bool   { return b.left }
func (b *ButtonIcon) IsMouseDownRight() bool  { return b.right }
func (b *ButtonIcon) IsMouseDownMiddle() bool { return b.middle }

func (b *ButtonIcon) Update(dt int) {
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

func (b *ButtonIcon) Draw(surface *ebiten.Image) {
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

func (b *ButtonIcon) Resize(rect []int) {
	b.rect = NewRect(rect)
	b.Dirty = true
	b.Image = nil
}

func (b *ButtonIcon) Close() {
	b.Image.Dispose()
}
