package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Label struct {
	*Drawable
	txt              string
	options          *text.DrawOptions
	hAlign, vAlign   text.Align
	font             *text.GoTextFace
	fontSize, margin float64
	dynamicFontSize  bool
	wrappedTxt       string
	txtPos           Point[float64]
}

func NewLabel(txt string) *Label {
	l := &Label{
		Drawable:        NewDrawable(),
		txt:             txt,
		hAlign:          text.AlignCenter,
		vAlign:          text.AlignCenter,
		fontSize:        12,
		dynamicFontSize: true,
	}
	l.font = GetUi().FontDefault().Get(int(l.fontSize))
	theme := GetUi().theme
	l.SetBg(theme.Get(TextBg))
	l.SetFg(theme.Get(TextFg))
	return l
}
func (l *Label) Layout() {
	if !l.IsDirty() && l.Image() != nil {
		return
	}

	// 1. Подготовка параметров шрифта
	if l.dynamicFontSize {
		l.fontSize = float64(GetUi().FontDefault().CalcFontSize(l.txt, l.Rect(), true))
	}

	// 2. Подготовка холста (Drawable.Layout создаст/очистит l.image)
	l.Drawable.Layout()

	w, h := l.Rect().Size()
	renderRect := NewRect([]int{0, 0, w, h})

	if l.pressed {
		l.margin = float64(l.Rect().GetLowestSize()) * 0.03
		renderRect.X += int(l.margin / 2)
		renderRect.Y += int(l.margin / 2)
	}

	GetUi().FontDefault().DrawString(l.Image(), l.txt, int(l.fontSize), renderRect, l.hAlign, l.vAlign, l.Fg(), true)

	l.ClearDirty()
}

func (l *Label) Text() string { return l.txt }
func (l *Label) SetText(value string) *Label {
	if value == l.txt {
		return l
	}
	l.txt = value
	l.MarkDirty()
	return l
}

func (l *Label) FontSize() float64 { return l.fontSize }
func (l *Label) SetFontSize(value float64) *Label {
	l.fontSize = value
	l.SetUseDynamicFont(false)
	l.ImageReset()
	return l
}

func (l *Label) SetUseDynamicFont(value bool) *Label {
	l.dynamicFontSize = value
	l.ImageReset()
	return l
}

func (l *Label) Draw(screen *ebiten.Image) {
	if l.IsHidden() {
		return
	}
	if l.IsDirty() {
		l.Layout()
	}
	l.Drawable.Draw(screen)
}

func (l *Label) SetAlign(horizontal, vertical text.Align) *Label {
	l.hAlign = horizontal
	l.vAlign = vertical
	l.ImageReset()
	return l
}
