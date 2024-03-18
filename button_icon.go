package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Умею показать кнопку под мышкой выделенной, или нажатой, или отпущенной(после отпускания исполняется прикрепленный метод), вторая иконка нажатая.
type ButtonIcon struct {
	DrawableBase
	btn          *Button
	icon1, icon2 *Icon
}

func NewButtonIcon(icons []*ebiten.Image, f func(*Button)) *ButtonIcon {
	b := &ButtonIcon{}
	b.SetupButtonIcon(icons, f)
	b.Visible(true)
	return b
}

func (b *ButtonIcon) SetupButtonIcon(icons []*ebiten.Image, f func(*Button)) {
	b.btn = NewButton("", f)
	b.icon1 = NewIcon(icons[0])
	b.icon2 = NewIcon(icons[1])
	b.SetImage(b.icon1.GetIcon())
}

func (b *ButtonIcon) SetFunc(f func(*Button)) {
	b.btn.onPressed = f
}

func (b *ButtonIcon) SetIcons(icons []*ebiten.Image) {
	b.icon1.SetIcon(icons[0])
	b.icon2.SetIcon(icons[1])
	b.Dirty = true
}

func (b *ButtonIcon) SetReleasedIcon(icon *ebiten.Image) {
	b.icon1.SetIcon(icon)
	b.Dirty = true
}

func (b *ButtonIcon) SetPressedIcon(icon *ebiten.Image) {
	b.icon2.SetIcon(icon)
	b.Dirty = true
}

func (b *ButtonIcon) Layout() {
	b.SpriteBase.Layout()
	b.Dirty = false
}

func (b *ButtonIcon) Update(dt int) {
	if b.disabled {
		return
	}
	b.btn.Update(dt)
	if b.btn.IsPressed() {
		b.SetImage(b.icon2.GetIcon())
	} else {
		b.SetImage(b.icon1.GetIcon())
	}
}

func (b *ButtonIcon) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	if b.Dirty {
		b.Layout()
		b.icon1.Layout()
		b.icon2.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := b.rect.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(b.Image(), op)
}

func (b *ButtonIcon) Resize(rect []int) {
	b.Rect(NewRect(rect))
	b.SpriteBase.Rect(NewRect(rect))
	b.btn.Resize(rect)
	b.icon1.Resize(rect)
	b.icon2.Resize(rect)
	b.ImageReset()
}
