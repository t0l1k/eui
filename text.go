package eui

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Метка текста в одной строке, размер текста вычисляется исходя из размера
type Text struct {
	DrawableBase
	txt               string
	fontInit, oneFont bool
	fontSize          int
	pos               PointInt
}

func NewText(txt string) *Text {
	t := &Text{
		txt: txt,
	}
	theme := GetUi().theme
	t.Bg(theme.Get(TextBg))
	t.Fg(theme.Get(TextFg))
	t.SetText(txt)
	t.Visible(true)
	return t
}

func (t *Text) OnlyOneFontSize(value bool) {
	t.oneFont = value
	t.Dirty = true
}

func (t *Text) GetText() string { return t.txt }
func (t *Text) SetText(value string) {
	if t.txt == value {
		return
	}
	t.txt = value
	t.Dirty = true
}

func (t *Text) UpdateData(value interface{}) {
	switch v := value.(type) {
	case string:
		t.SetText(v)
	case int:
		t.SetText(strconv.Itoa(v))
	}
}

func (t *Text) Layout() {
	t.SpriteBase.Layout()
	t.Image().Fill(t.GetBg())
	var font font.Face
	if !t.oneFont || !t.fontInit {
		t.fontSize = GetFonts().calcFontSize(t.txt, t.rect)
		font = GetFonts().Get(t.fontSize)
		b := text.BoundString(font, t.txt)
		t.pos.X = (t.rect.W - b.Max.X) / 2
		t.pos.Y = t.rect.H - (t.rect.H+b.Min.Y)/2
		if !t.fontInit {
			t.fontInit = true
		}
	} else if t.oneFont {
		font = GetFonts().Get(t.fontSize)
		b := text.BoundString(font, t.txt)
		t.pos.X = (t.rect.W - b.Max.X) / 2
		t.pos.Y = t.rect.H - (t.rect.H+b.Min.Y)/2
	}
	text.Draw(t.image, t.txt, font, t.pos.X, t.pos.Y, t.fg)
	t.Dirty = false
}

func (t *Text) Draw(surface *ebiten.Image) {
	if !t.IsVisible() {
		return
	}
	if t.Dirty {
		t.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := t.GetRect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(t.Image(), op)
}

func (t *Text) Resize(rect []int) {
	t.Rect(NewRect(rect))
	t.ImageReset()
}
