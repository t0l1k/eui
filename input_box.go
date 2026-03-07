package eui

import (
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

// === InputLine ===
type InputLine struct {
	*Drawable
	text        string
	textChanged *Signal[string]
	onReturn    SlotFunc[*InputLine]
	cursorPos   int
	blink       bool
	tick        time.Duration
	maxLen      int
	placeHolder string
}

func NewInputLine(onReturn SlotFunc[*InputLine]) *InputLine {
	e := &InputLine{
		Drawable:    NewDrawable(),
		onReturn:    onReturn,
		maxLen:      0,
		textChanged: NewSignal(func(a, b string) bool { return a == b }),
	}
	e.SetBg(colornames.Blue)
	e.SetFg(colornames.Yellow)
	return e
}

func (b *InputLine) WantBlur() bool { return false }
func (e *InputLine) Tick(td TickData) {
	if !e.State().IsFocused() || e.IsHidden() {
		return
	}
	e.tick += td.Duration()
	if e.tick > 500*time.Millisecond {
		e.blink = !e.blink
		e.tick = 0
		e.MarkDirty()
	}
}

func (e *InputLine) Hit(pt Point[int]) Drawabler {
	if !pt.In(e.rect) || e.IsHidden() {
		return nil
	}
	return e
}

func (e *InputLine) OnReturn() SlotFunc[*InputLine] { return e.onReturn }
func (e *InputLine) TextChanged() *Signal[string]   { return e.textChanged }
func (e *InputLine) Text() string                   { return e.text }
func (e *InputLine) SetText(value string) {
	e.text = value
	e.cursorPos = len(e.text)
	e.MarkDirty()
}
func (e *InputLine) SetPlaceholder(value string) *InputLine { e.placeHolder = value; return e }
func (e *InputLine) MaxLen() int                            { return e.maxLen }
func (e *InputLine) SetMaxLen(value int)                    { e.maxLen = value }

func (e *InputLine) KeyPressed(kd KeyboardData) {
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
		e.textChanged.Emit(e.Text())
	}
	e.MarkDirty()
}

func (e *InputLine) KeyReleased(kd KeyboardData) {
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
			e.textChanged.Emit(e.Text())
		}
	case kd.IsReleased(ebiten.KeyDelete):
		if e.cursorPos < len(e.text) && len(e.text) > 0 {
			e.text = e.text[:e.cursorPos] + e.text[e.cursorPos+1:]
			e.textChanged.Emit(e.Text())
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

func (e *InputLine) Layout() {
	h := e.rect.H
	e.Drawable.Layout()

	sz := float64(h) * 0.3
	fnt := GetUi().FontDefault().Get(int(sz))

	margin := 8.0
	w := float64(e.rect.W)
	txt := e.text
	if txt == "" && !e.state.IsFocused() {
		txt = e.placeHolder
	}

	// Определяем ширину текста до курсора
	cursorSub := ""
	if e.cursorPos > 0 && e.cursorPos <= len(txt) {
		cursorSub = txt[:e.cursorPos]
	}
	cursorW, _ := text.Measure(cursorSub, fnt, fnt.Size*1.2)
	cursorX := margin + cursorW

	visibleW := w - margin*2

	// Скроллим так, чтобы курсор всегда был видим
	scrollOffset := 0.0
	if cursorX > margin+visibleW {
		scrollOffset = cursorX - (margin + visibleW)
	} else if cursorX < margin {
		scrollOffset = cursorX - margin
	}

	// Обрезаем текст слева, если нужно
	drawText := txt
	drawOffset := 0
	for len(drawText) > 0 {
		subW, _ := text.Measure(drawText, fnt, fnt.Size*1.2)
		if float64(subW) > visibleW+scrollOffset {
			drawText = drawText[1:]
			drawOffset++
		} else {
			break
		}
	}

	opts := &text.DrawOptions{}
	opts.GeoM.Translate(margin-scrollOffset, float64(h)/2)
	opts.PrimaryAlign = text.AlignStart
	opts.SecondaryAlign = text.AlignCenter
	opts.ColorScale.ScaleWithColor(e.Fg())

	text.Draw(e.Image(), txt[drawOffset:], fnt, opts)

	// Курсор
	if e.State().IsFocused() {
		// Пересчитаем позицию курсора относительно видимого текста
		cursorSub := txt[drawOffset:e.cursorPos]
		cursorW, _ := text.Measure(cursorSub, fnt, fnt.Size*1.2)
		cursorX := margin + float64(cursorW) - scrollOffset

		cursorColor := e.Fg()
		if e.blink {
			cursorColor = e.Bg()
		}
		cursorHeight := float32(sz * 0.9)
		cursorWidth := float32(sz * 0.01) // 1% от высоты шрифта
		if cursorWidth < 1 {
			cursorWidth = 1
		}
		cursorY := float32(h)/2 - cursorHeight/2
		vector.StrokeLine(
			e.Image(),
			float32(cursorX), cursorY,
			float32(cursorX), cursorY+cursorHeight,
			cursorWidth, cursorColor, true,
		)
	}
	e.ClearDirty()
}

func (w *InputLine) Draw(screen *ebiten.Image) {
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

func (t *InputLine) Digit() (float64, error) {
	if len(t.text) == 0 {
		return 0, nil
	}
	n, err := strconv.ParseFloat(t.text, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
