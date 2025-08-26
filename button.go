package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Умею показать кнопку под мышкой выделенной, или нажатой, или отпущенной(после отпускания исполняется прикрепленный метод)
type Button struct {
	*Drawable
	txt        string
	onReleased func(*Button)
	icons      []*ebiten.Image
	useIcon    bool
}

func NewButton(txt string, fn func(*Button)) *Button {
	b := &Button{Drawable: NewDrawable(), txt: txt, onReleased: fn}
	theme := GetUi().theme
	b.SetBg(theme.Get(ButtonBg))
	b.SetFg(theme.Get(ButtonFg))
	b.SetShadow(1, color.RGBA{0, 0, 0, 128})
	return b
}

func NewButtonIcon(icons []*ebiten.Image, fn func(*Button)) *Button {
	return &Button{Drawable: NewDrawable(), icons: icons, onReleased: fn, useIcon: true}
}

func (b *Button) Text() string         { return b.txt }
func (b *Button) SetText(value string) { b.txt = value; b.MarkDirty() }

func (b *Button) Hit(pt Point[int]) Drawabler {
	if !pt.In(b.rect) || b.IsHidden() || b.IsDisabled() {
		return nil
	}
	return b
}
func (b *Button) WantBlur() bool  { return true }
func (b *Button) IsPressed() bool { return b.pressed }

func (b *Button) MouseDown(md MouseData) {
	b.pressed = true
	b.MarkDirty()
}
func (b *Button) MouseUp(md MouseData) {
	if b.onReleased != nil {
		b.onReleased(b)
	}
	b.pressed = false
}

func (b *Button) SetIcons(icons []*ebiten.Image) { b.icons = icons; b.MarkDirty() }
func (b *Button) IconUp() *ebiten.Image          { return b.icons[0] }
func (b *Button) IconDown() *ebiten.Image        { return b.icons[1] }
func (b *Button) SetUpIcon(icon *ebiten.Image)   { b.icons[0] = icon; b.MarkDirty() }
func (b *Button) SetDownIcon(icon *ebiten.Image) { b.icons[1] = icon; b.MarkDirty() }

func (b *Button) Layout() {
	b.Drawable.Layout()

	w, h := b.Rect().Size()

	if b.useIcon {
		var icon *Icon
		if b.pressed {
			icon = NewIcon(b.IconDown())
		} else {
			icon = NewIcon(b.IconUp())
		}
		icon.SetRect(NewRect([]int{0, 0, w, h}))
		icon.Layout()
		icon.Draw(b.Image())
	} else {
		lbl := NewLabel(b.txt)
		margin := int(float64(b.Rect().GetLowestSize()) * 0.03)
		x := margin
		y := margin
		if b.pressed {
			x += margin / 2
			y += margin / 2
		}
		lbl.SetRect(NewRect([]int{x, y, w - margin*2, h - margin*2}))
		lbl.SetBg(b.Bg())
		lbl.SetFg(b.Fg())
		lbl.Layout()
		lbl.Draw(b.Image())
		vector.StrokeRect(b.Image(), 0, 0, float32(w), float32(h), float32(margin), b.state.Color(), true)
	}
	b.ClearDirty()
}

func (b *Button) Draw(surface *ebiten.Image) {
	if b.IsHidden() {
		return
	}
	if b.IsDirty() {
		b.Layout()
	}
	b.Drawable.Draw(surface)
}

func (d *Button) MouseEnter() {
	if d.IsDisabled() {
		return
	}
	if !d.State().IsHovered() {
		d.SetState(StateHover)
	}
}
func (d *Button) MouseLeave() {
	if d.IsDisabled() {
		return
	}
	if d.State().IsHovered() {
		d.SetState(StateNormal)
	}
}
