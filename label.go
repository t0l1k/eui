package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type LabelAlign int

const (
	LabelAlignLeft LabelAlign = iota
	LabelAlignCenter
	LabelAlignRight
	LabelAlignUp
	LabelAlignDown
	LabelAlignLeftUp
)

func (a LabelAlign) String() string {
	return [...]string{
		"Align Left",
		"Align Center",
		"Align Right",
		"Align Up",
		"Align Down",
		"Align LeftUp",
	}[a]
}

type Label struct {
	*Drawable
	txt              string
	options          *text.DrawOptions
	align            LabelAlign
	font             *text.GoTextFace
	fontSize, margin float64
	dynamicFontSize  bool
	txtPos           Point[float64]
}

func NewLabel(txt string) *Label {
	l := &Label{
		Drawable:        NewDrawable(),
		txt:             txt,
		align:           LabelAlignCenter,
		fontSize:        12,
		dynamicFontSize: true,
	}
	l.font = GetUi().FontDefault().Get(int(l.fontSize))
	theme := GetUi().theme
	l.SetBg(theme.Get(TextBg))
	l.SetFg(theme.Get(TextFg))
	return l
}
func (l *Label) setupOptions() {
	l.options = &text.DrawOptions{}
	l.options.ColorScale.Reset()
	l.options.ColorScale.ScaleWithColor(l.Fg())
	l.options.LineSpacing = l.fontSize * 1.2
	l.MarkDirty()
}
func (l *Label) Layout() {
	var txt string
	if l.Image() == nil {
		l.margin = float64(l.Rect().GetLowestSize()) * 0.03
		if l.dynamicFontSize {
			sz := GetUi().FontDefault().CalcFontSize(l.txt, l.Rect())
			if sz != int(l.fontSize) {
				l.fontSize = float64(sz)
			}
		}
		l.font = GetUi().FontDefault().Get(int(l.fontSize))
		txt, _ = GetUi().FontDefault().WordWrapText(l.txt, l.fontSize, l.rect.Width())
		w0, h0 := l.Rect().Size()
		w, h := text.Measure(txt, l.font, l.font.Size*1.2)
		if w > float64(w0) || h > float64(h0) {
			x, y := l.Rect().Pos()
			l.SetRect(NewRect([]int{x, y, int(w0), int(h)}))
		}
		l.options = nil
	}
	if l.options == nil {
		l.setupOptions()
		l.txtPos.X, l.txtPos.Y = l.calcAlign()
	}
	x, y := l.txtPos.Get()
	if l.pressed {
		x += l.margin / 2
		y += l.margin / 2
	}
	l.Drawable.Layout()
	l.options.GeoM.Reset()
	l.options.GeoM.Translate(x, y)
	if l.txt != txt && len(txt) > 0 {
		text.Draw(l.Image(), txt, l.font, l.options)
	} else {
		text.Draw(l.Image(), l.txt, l.font, l.options)
	}
	l.ClearDirty()
}

func (l *Label) Text() string { return l.txt }
func (l *Label) SetText(value string) *Label {
	if value == l.txt {
		return l
	}
	l.txt = value
	l.ImageReset()
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

func (l *Label) SetAlign(value LabelAlign) *Label {
	l.align = value
	l.ImageReset()
	return l
}

func (l *Label) calcAlign() (x, y float64) {
	switch l.align {
	case LabelAlignCenter:
		x, y = float64(l.rect.Width())/2, float64(l.rect.Height())/2
		l.options.GeoM.Reset()
		l.options.GeoM.Translate(x, y)
		l.options.PrimaryAlign = text.AlignCenter
		l.options.SecondaryAlign = text.AlignCenter
	case LabelAlignDown:
		x, y = float64(l.rect.Width())/2, float64(l.rect.Height())
		l.options.GeoM.Reset()
		l.options.GeoM.Translate(x, y)
		l.options.PrimaryAlign = text.AlignCenter
		l.options.SecondaryAlign = text.AlignEnd
	case LabelAlignLeft:
		x, y = 0, float64(l.rect.Height())/2
		l.options.GeoM.Reset()
		l.options.GeoM.Translate(x, y)
		l.options.PrimaryAlign = text.AlignStart
		l.options.SecondaryAlign = text.AlignCenter
	case LabelAlignLeftUp:
		x, y = 0, 0
		l.options.GeoM.Reset()
		l.options.GeoM.Translate(x, y)
		l.options.PrimaryAlign = text.AlignStart
		l.options.SecondaryAlign = text.AlignStart
	case LabelAlignRight:
		x, y = float64(l.rect.Width()), float64(l.rect.Height())/2
		l.options.GeoM.Reset()
		l.options.GeoM.Translate(x, y)
		l.options.PrimaryAlign = text.AlignEnd
		l.options.SecondaryAlign = text.AlignCenter
	case LabelAlignUp:
		x, y = float64(l.rect.Width())/2, 0
		l.options.GeoM.Reset()
		l.options.GeoM.Translate(x, y)
		l.options.PrimaryAlign = text.AlignCenter
		l.options.SecondaryAlign = text.AlignStart
	}
	return x, y
}
