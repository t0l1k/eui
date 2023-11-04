package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Умею показать кнопку под мышкой выделенной, или нажатой, или отпущенной(после отпускания исполняется прикрепленный метод)
type Button struct {
	Text
	onPressed           func(*Button)
	buttonPressed       bool
	left, right, middle bool
}

func NewButton(text string, bg, fg color.Color, f func(*Button)) *Button {
	b := &Button{
		onPressed: f,
	}
	b.SetupButton(text, bg, fg, f)
	return b
}

func (b *Button) SetupButton(text string, bg, fg color.Color, f func(*Button)) {
	b.SetupText(text, bg, fg)
	b.Name("button")
}

func (b *Button) IsMouseDownLeft() bool {
	return b.left
}

func (b *Button) IsMouseDownRight() bool {
	return b.right
}

func (b *Button) IsMouseDownMiddle() bool {
	return b.middle
}

func (b *Button) Layout() {
	b.Text.Layout()
	var fg color.Color
	switch b.state {
	case ViewStateHover:
		fg = Yellow
	case ViewStateFocus:
		fg = Red
	case ViewStateNormal:
		fg = Black
	case ViewStateSelected:
		fg = Blue
	case ViewStateDisabled:
		fg = Purple
	case ViewStateActive:
		fg = White
	}
	_, _, w, h := b.rect.GetRectFloat()
	bold := 2
	if b.buttonPressed {
		bold = 5
	}
	vector.StrokeRect(b.image, 0, 0, w, h, float32(bold), fg, true)
	b.dirty = false
}

func (b *Button) Pressed(value bool) {
	b.buttonPressed = value
	if value {
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
}

func (b *Button) Update(dt int) {
	if b.disabled {
		return
	}
	if b.state == ViewStateFocus && !b.buttonPressed {
		b.Pressed(true)
	}
	if (b.state == ViewStateHover || b.state == ViewStateExec) && b.buttonPressed {
		if b.onPressed != nil {
			b.onPressed(b)
		}
		b.Pressed(false)
		if b.state == ViewStateExec {
			b.state = ViewStateNormal
		}
	}
	if b.state == ViewStateNormal {
		b.Pressed(false)
	}
}

func (b *Button) Draw(surface *ebiten.Image) {
	if !b.visible {
		return
	}
	if b.dirty {
		b.Layout()
		for _, c := range b.Container {
			c.Layout()
		}
	}
	op := &ebiten.DrawImageOptions{}
	x, y := b.rect.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(b.image, op)
	for _, v := range b.Container {
		v.Draw(surface)
	}
}
