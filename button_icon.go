package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Умею показать кнопку под мышкой выделенной, или нажатой, или отпущенной(после отпускания исполняется прикрепленный метод), вторая иконка нажатая.
type ButtonIcon struct {
	*Drawable
	btn          *Button
	icon1, icon2 *Icon
}

func NewButtonIcon(icons []*ebiten.Image, f func(*Button)) *ButtonIcon {
	b := &ButtonIcon{Drawable: NewDrawable()}
	b.SetupButtonIcon(icons, f)
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
	b.MarkDirty()
}

func (b *ButtonIcon) SetReleasedIcon(icon *ebiten.Image) {
	b.icon1.SetIcon(icon)
	b.MarkDirty()
}

func (b *ButtonIcon) SetPressedIcon(icon *ebiten.Image) {
	b.icon2.SetIcon(icon)
	b.MarkDirty()
}

func (b *ButtonIcon) Layout() {
	b.Drawable.Layout()
	b.ClearDirty()
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
	if b.IsHidden() {
		return
	}
	if b.IsDirty() {
		b.Layout()
		b.icon1.Layout()
		b.icon2.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := b.rect.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(b.Image(), op)
}

func (b *ButtonIcon) SetRect(rect Rect[int]) {
	b.Drawable.SetRect(rect)
	b.btn.SetRect(rect)
	b.icon1.SetRect(rect)
	b.icon2.SetRect(rect)
	b.ImageReset()
}
