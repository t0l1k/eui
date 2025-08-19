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
	txt             string
	options         *text.DrawOptions
	align           LabelAlign
	fontSize        float64
	dynamicFontSize bool
}

func NewLabel(txt string) *Label {
	l := &Label{
		Drawable:        NewDrawable(),
		txt:             txt,
		align:           LabelAlignCenter,
		fontSize:        12,
		dynamicFontSize: true,
	}
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
	fnt := GetUi().FontDefault().Get(int(l.fontSize))
	var txt string

	if l.dynamicFontSize {
		sz := GetUi().FontDefault().CalcFontSize(l.txt, l.Rect())
		if sz != int(l.fontSize) {
			fnt = GetUi().FontDefault().Get(int(sz))
			l.fontSize = float64(sz)
			l.options = nil
		}
	} else {
		txt, _ = GetUi().FontDefault().WordWrapText(l.txt, l.fontSize, l.rect.Width())

		w0, h0 := l.Rect().Size()
		w, h := text.Measure(txt, fnt, fnt.Size*1.2)
		if w > float64(w0) || h > float64(h0) {
			l.image = nil
			x, y := l.Rect().Pos()
			l.SetRect(NewRect([]int{x, y, int(w0), int(h)}))
		}
	}

	if l.options == nil {
		l.setupOptions()
	}
	x, y := l.calcAlign()
	margin := float64(l.Rect().GetLowestSize()) * 0.03
	if l.pressed {
		x += margin / 2
		y += margin / 2
	}

	l.Drawable.Layout()

	l.options.GeoM.Reset()
	l.options.GeoM.Translate(x, y)
	if l.txt != txt && len(txt) > 0 {
		text.Draw(l.Image(), txt, fnt, l.options)
	} else {
		text.Draw(l.Image(), l.txt, fnt, l.options)
	}
	l.ClearDirty()
}

func (l *Label) Text() string { return l.txt }
func (l *Label) SetText(value string) Drawabler {
	if value == l.txt {
		return l
	}
	l.txt = value
	l.MarkDirty()
	return l
}

func (l *Label) SetFontSize(value float64) Drawabler {
	l.fontSize = value
	l.SetUseDynamicFont(false)
	l.MarkDirty()
	return l
}

func (l *Label) SetUseDynamicFont(value bool) {
	l.dynamicFontSize = value
	l.MarkDirty()
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

func (s *Label) SetRect(rect Rect[int]) { s.rect = rect; s.ImageReset() }

func (l *Label) SetAlign(value LabelAlign) {
	l.align = value
	l.MarkDirty()
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
