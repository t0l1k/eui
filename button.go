package eui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Умею показать кнопку под мышкой выделенной, или нажатой, или отпущенной(после отпускания исполняется прикрепленный метод)
type Button struct {
	*Drawable
	txt            string
	onPressed      func(*Button)
	icons          []*ebiten.Image
	useIcon, press bool
}

func NewButton(txt string, fn func(*Button)) *Button {
	b := &Button{Drawable: NewDrawable(), txt: txt, onPressed: fn}
	theme := GetUi().theme
	b.Bg(theme.Get(ButtonBg))
	b.Fg(theme.Get(ButtonFg))
	return b
}

func NewButtonIcon(icons []*ebiten.Image, fn func(*Button)) *Button {
	return &Button{Drawable: NewDrawable(), icons: icons, onPressed: fn, useIcon: true}
}

func (b *Button) Text() string         { return b.txt }
func (b *Button) SetText(value string) { b.txt = value; b.MarkDirty() }

func (b *Button) Hit(pt Point[int]) Drawabler {
	if !pt.In(b.rect) || b.IsHidden() {
		return nil
	}
	log.Println("Button:Hit:", b.txt, b.Rect(), pt)
	return b
}
func (b *Button) WantBlur() bool { return true }
func (b *Button) MouseDown(md MouseData) {
	b.press = true
	b.MarkDirty()
}
func (b *Button) MouseUp(md MouseData) {
	if b.onPressed != nil {
		b.onPressed(b)
	}
	b.press = false
	log.Println("Button:MouseReleased:", b.txt, b.Rect())
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
		if b.press {
			icon = NewIcon(b.IconUp())
		} else {
			icon = NewIcon(b.IconDown())
		}
		icon.SetRect(NewRect([]int{0, 0, w, h}))
		icon.Layout()
		icon.Draw(b.Image())
		log.Printf("btnStatus: IconUp=%v IconDown=%v useIcon=%v len:%v", b.IconUp(), b.IconDown(), b.useIcon, len(b.icons))
		log.Println("ButtonIcon:Layout:", b.Rect(), icon.Rect(), b.press, b.State(), icon.State(), icon.Image().Bounds())
	} else {
		lbl := NewText(b.txt)
		margin := int(float64(b.Rect().GetLowestSize()) * 0.03)
		lbl.SetRect(NewRect([]int{margin, margin, w - margin*2, h - margin*2}))
		lbl.Bg(b.GetBg())
		lbl.Fg(b.GetFg())
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
	if !d.State().IsHovered() {
		d.SetState(StateHover)
	}
}
func (d *Button) MouseLeave() {
	if d.State().IsHovered() {
		d.SetState(StateNormal)
	}
}
