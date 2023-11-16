package eui

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Метка текста в одной строке, размер текста вычисляется исходя из размера
type Text struct {
	View
	text              string
	fontInit, oneFont bool
	fontSize          int
	pos               PointInt
}

func NewText(text string) *Text {
	t := &Text{
		text: text,
	}
	t.SetupText(text)
	return t
}

func (t *Text) SetupText(text string) {
	t.SetupView()
	theme := GetUi().theme
	t.Bg(theme.Get(TextBg))
	t.Fg(theme.Get(TextFg))
	t.SetText(text)
}

func (t *Text) OnlyOneFontSize(value bool) {
	t.oneFont = value
	t.dirty = true
}

func (t *Text) GetText() string {
	return t.text
}

func (t *Text) SetText(value string) {
	if t.text == value {
		return
	}
	t.text = value
	t.dirty = true
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
	t.View.Layout()
	var font font.Face
	if !t.oneFont || !t.fontInit {
		t.fontSize = GetFonts().calcFontSize(t.text, t.rect)
		font = GetFonts().Get(t.fontSize)
		b := text.BoundString(font, t.text)
		t.pos.X = (t.rect.W - b.Max.X) / 2
		t.pos.Y = t.rect.H - (t.rect.H+b.Min.Y)/2
		if !t.fontInit {
			t.fontInit = true
		}
	} else if t.oneFont {
		font = GetFonts().Get(t.fontSize)
		b := text.BoundString(font, t.text)
		t.pos.X = (t.rect.W - b.Max.X) / 2
		t.pos.Y = t.rect.H - (t.rect.H+b.Min.Y)/2
	}
	text.Draw(t.image, t.text, font, t.pos.X, t.pos.Y, t.fg)
	t.dirty = false
}

func (t *Text) Draw(surface *ebiten.Image) {
	if !t.visible {
		return
	}
	if t.dirty {
		t.Layout()
		for _, c := range t.Container {
			c.Layout()
		}
	}
	op := &ebiten.DrawImageOptions{}
	x, y := t.rect.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(t.image, op)
	for _, v := range t.Container {
		v.Draw(surface)
	}
}
