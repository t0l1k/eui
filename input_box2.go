package eui

import (
	"image/color"
	"strconv"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui/colors"
)

type InputBox2 struct {
	DrawableBase
	lbl                   *Text
	textVar               *SubjectBase
	hasFocus, alwaysFocus bool
	onReturn              func(*InputBox2)
	onlyDigits            bool
	bg, fg                color.Color
	state                 InputState
}

func NewInputBox2(onReturn func(*InputBox2)) *InputBox2 {
	i := &InputBox2{}
	i.onReturn = onReturn
	i.textVar = NewSubject()
	i.lbl = NewText("")
	i.textVar.Attach(i.lbl)
	i.textVar.SetValue("")
	i.Add(i.lbl)
	theme := GetUi().theme
	i.bg = theme.Get(InputBoxBg)
	i.fg = theme.Get(InputBoxFg)
	GetUi().inputMouse.Attach(i)
	i.alwaysFocus = false
	return i
}

func NewDigitInputBox2(txt string, ln int, onReturn func(*InputBox2)) *InputBox2 {
	i := NewInputBox2(onReturn)
	i.onlyDigits = true
	i.lbl.SetText(txt)
	return i
}

func (i *InputBox2) GetText() string      { return i.textVar.Value().(string) }
func (i *InputBox2) Reset()               { i.textVar.Reset() }
func (i *InputBox2) SetText(value string) { i.textVar.SetValue(value) }

func (i *InputBox2) IsFocused() bool  { return i.hasFocus }
func (i *InputBox2) SetFocus()        { i.hasFocus = true }
func (i *InputBox2) Blur()            { i.hasFocus = false }
func (i *InputBox2) SetAlwaysFocus()  { i.alwaysFocus = true }
func (i *InputBox2) BlurAlwaysFocus() { i.alwaysFocus = false }

func (i *InputBox2) GetDigit() (float64, error) {
	if len(i.textVar.Value().(string)) == 0 {
		return 0, nil
	}
	n, err := strconv.ParseFloat(i.textVar.Value().(string), 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (i *InputBox2) checkRepeat(key ebiten.Key) bool {
	const (
		delay    = 50
		interval = 5
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 || d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (i *InputBox2) Update(dt int) {
	if !i.IsFocused() {
		return
	}
	var runes []rune
	runes = ebiten.AppendInputChars(runes[:0])
	if len(runes) > 0 {
		txt := i.textVar.Value().(string)
		txt += string(runes)
		i.textVar.SetValue(txt)
	}
	if i.checkRepeat(ebiten.KeyBackspace) {
		if len(i.textVar.Value().(string)) > 0 {
			txt := i.textVar.Value().(string)
			txt = txt[:len(txt)-1]
			i.textVar.SetValue(txt)
		}
	}
	if i.onlyDigits {
		i.lbl.Bg(i.bg)
		for _, v := range i.textVar.Value().(string) {
			if !unicode.IsDigit(v) && !(v == '.' || v == ',') {
				i.lbl.Bg(colors.Red)
			}
		}
	}
	if i.checkRepeat(ebiten.KeyEnter) {
		if i.onReturn != nil {
			i.onReturn(i)
		}
	}
}

func (i *InputBox2) GetState() InputState {
	return i.state
}

func (i *InputBox2) SetState(state InputState) {
	if i.state == state {
		return
	}
	i.state = state
	i.ImageReset()
}

func (i *InputBox2) UpdateInput(value interface{}) {
	if i.disabled || i.alwaysFocus {
		return
	}
	switch vl := value.(type) {
	case MouseData:
		x, y, b := vl.position.X, vl.position.Y, vl.button
		inRect := i.rect.InRect(x, y)
		if inRect {
			if b == buttonReleased {
				if i.state == ViewStateNormal {
					i.SetState(ViewStateHover)
				}
				if i.state == ViewStateFocus {
					i.SetState(ViewStateHover)
				}
			}
			if b == buttonPressed {
				if i.state == ViewStateHover {
					i.SetState(ViewStateFocus)
					i.SetFocus()
				}
			}
		} else if i.state != ViewStateNormal {
			i.SetState(ViewStateNormal)
			i.Blur()
		}
	}
}

func (i *InputBox2) Resize(rect []int) {
	i.Rect(NewRect(rect))
	i.lbl.Resize(rect)
	i.ImageReset()
}

func (i *InputBox2) Close() {
	GetUi().inputMouse.Detach(i)
}
