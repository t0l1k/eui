package eui

import (
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

// === TextInputLine ===
type TextInputLine struct {
	*Drawable
	text      string
	onReturn  func(*TextInputLine)
	cursorPos int
	blink     bool
	tick      time.Duration
	maxLen    int
}

func NewTextInputLine(onReturn func(*TextInputLine)) *TextInputLine {
	e := &TextInputLine{
		Drawable: NewDrawable(),
		onReturn: onReturn,
		maxLen:   0,
	}
	e.Bg(colornames.Navy)
	e.Fg(colornames.Yellow)
	return e
}

func (b *TextInputLine) WantBlur() bool { return false }
func (e *TextInputLine) Tick(td TickData) {
	if !e.State().IsFocused() {
		return
	}
	e.tick += td.Duration()
	if e.tick > 500*time.Millisecond {
		e.blink = !e.blink
		e.tick = 0
		e.MarkDirty()
	}
}

func (e *TextInputLine) Hit(pt Point[int]) Drawabler {
	if !pt.In(e.rect) {
		return nil
	}
	return e
}

func (e *TextInputLine) Text() string { return e.text }
func (e *TextInputLine) SetText(value string) {
	e.text = value
	e.cursorPos = len(e.text)
	e.MarkDirty()
}

func (e *TextInputLine) MaxLen() int         { return e.maxLen }
func (e *TextInputLine) SetMaxLen(value int) { e.maxLen = value }

func (e *TextInputLine) KeyPressed(kd KeyboardData) {
	if !e.State().IsFocused() {
		return
	}
	// Вставка символов в позицию курсора
	for _, r := range kd.GetRunes() {
		if e.maxLen > 0 && len(e.text) >= e.maxLen {
			break
		}
		if e.cursorPos < 0 {
			e.cursorPos = 0
		}
		if e.cursorPos > len(e.text) {
			e.cursorPos = len(e.text)
		}
		e.text = e.text[:e.cursorPos] + string(r) + e.text[e.cursorPos:]
		e.cursorPos++
	}
	e.MarkDirty()
}

func (e *TextInputLine) KeyReleased(kd KeyboardData) {
	if !e.State().IsFocused() {
		return
	}
	switch {
	case kd.IsReleased(ebiten.KeyEnter):
		if e.onReturn != nil {
			e.onReturn(e)
		}
	case kd.IsReleased(ebiten.KeyBackspace):
		if e.cursorPos > 0 && len(e.text) > 0 {
			e.text = e.text[:e.cursorPos-1] + e.text[e.cursorPos:]
			e.cursorPos--
		}
	case kd.IsReleased(ebiten.KeyDelete):
		if e.cursorPos < len(e.text) && len(e.text) > 0 {
			e.text = e.text[:e.cursorPos] + e.text[e.cursorPos+1:]
		}
	case kd.IsReleased(ebiten.KeyArrowLeft):
		if e.cursorPos > 0 {
			e.cursorPos--
		}
	case kd.IsReleased(ebiten.KeyArrowRight):
		if e.cursorPos < len(e.text) {
			e.cursorPos++
		}
	}
	e.MarkDirty()
}

func (e *TextInputLine) Layout() {
	e.Drawable.Layout()

	fontSize := float64(e.rect.GetLowestSize()) * 0.3
	font := GetFonts().Get(int(fontSize))

	w, h := e.Rect().Size()
	margin := float64(e.Rect().GetLowestSize()) * 0.03

	x := margin
	textBounds := text.BoundString(font, e.text)
	y := (e.rect.H+textBounds.Dy())/2 - textBounds.Max.Y

	cursorSub := ""
	if e.cursorPos > 0 && e.cursorPos <= len(e.text) {
		cursorSub = e.text[:e.cursorPos]
	}
	cursorBounds := text.BoundString(font, cursorSub)
	cursorX := float64(x) + float64(cursorBounds.Dx()) + margin

	visibleWidth := float64(w) - margin*4

	scrollOffset := 0.0
	if cursorX > float64(x)+visibleWidth {
		scrollOffset = cursorX - (float64(x) + visibleWidth)
	} else if cursorX < float64(x) {
		scrollOffset = cursorX - float64(x)
	}

	if e.blink && e.State().IsFocused() {
		cursorTop := float32(y + textBounds.Min.Y)
		cursorBottom := float32(y + textBounds.Max.Y)

		cursorWidth := float32(fontSize * 0.05)
		if cursorWidth < 1 {
			cursorWidth = 1
		}
		vector.StrokeLine(e.Image(), float32(cursorX-scrollOffset), cursorTop, float32(cursorX-scrollOffset), cursorBottom, cursorWidth, e.GetFg(), true)
	}

	text.Draw(e.Image(), e.text, font, int(x-scrollOffset), y, e.GetFg())

	vector.StrokeRect(e.Image(), 0, 0, float32(w), float32(h), float32(margin), e.state.Color(), true)

	e.ClearDirty()
}

func (w *TextInputLine) Draw(screen *ebiten.Image) {
	if w.IsHidden() {
		return
	}
	if w.IsDirty() {
		w.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.rect.X), float64(w.rect.Y))
	screen.DrawImage(w.Image(), op)
}

func (t *TextInputLine) Digit() (float64, error) {
	if len(t.text) == 0 {
		return 0, nil
	}
	n, err := strconv.ParseFloat(t.text, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
