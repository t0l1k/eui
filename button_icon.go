package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Умею показать кнопку под мышкой выделенной, или нажатой, или отпущенной(после отпускания исполняется прикрепленный метод), вторая иконка нажатая.
type ButtonIcon struct {
	Icon
	icons               []*ebiten.Image
	onPressed           func(*ButtonIcon)
	buttonPressed       bool
	left, right, middle bool
}

func NewButtonIcon(icons []*ebiten.Image, f func(*ButtonIcon)) *ButtonIcon {
	b := &ButtonIcon{
		onPressed: f,
	}
	b.SetupButtonIcon(icons, f)
	return b
}

func (b *ButtonIcon) SetupButtonIcon(icons []*ebiten.Image, f func(*ButtonIcon)) {
	b.icons = icons
	b.SetupIcon(icons[0])
	b.Name("ButtonIcon")
}

func (b *ButtonIcon) SetReleasedIcon(icon *ebiten.Image) {
	if b.icons[0] == icon {
		return
	}
	b.icons[0] = icon
	b.dirty = true
}

func (b *ButtonIcon) SetPressedIcon(icon *ebiten.Image) {
	if b.icons[1] == icon {
		return
	}
	b.icons[1] = icon
	b.dirty = true
}

func (b *ButtonIcon) IsMouseDownLeft() bool {
	return b.left && b.buttonPressed || b.buttonPressed && b.state == ViewStateExec
}

func (b *ButtonIcon) IsMouseDownRight() bool {
	return b.right
}

func (b *ButtonIcon) IsMouseDownMiddle() bool {
	return b.middle
}

func (b *ButtonIcon) Layout() {
	b.Icon.Layout()
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

func (b *ButtonIcon) Pressed(value bool) {
	b.buttonPressed = value
	if value {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.left = true
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			b.right = true
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			b.middle = true
		}
		b.SetIcon(b.icons[1])
	} else {
		b.left = false
		b.right = false
		b.middle = false
		b.SetIcon(b.icons[0])
	}
}

func (b *ButtonIcon) Update(dt int) {
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

func (b *ButtonIcon) Draw(surface *ebiten.Image) {
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
