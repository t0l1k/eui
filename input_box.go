package eui

import (
	"strconv"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

// Умею получить от клавиатуры нажатый символ(пока только английский), backspace удаляет последний введенный символ, enter запускает прикрепленный метод. Умею мигать курсором. Только при фокусе от клавиатуры доступен ввод, активированый нажатием левой кнопки мыши(меняется виджет обрамление былым). Есть проверка только на цифры, выбрать при создании экземпляра метод с настроенной проверкой на цифры.
type InputBox struct {
	View
	_text         string
	btn           *Button
	size          int
	cursor        *Cursor
	timerFlashing *Timer
	keyboardState InputKeyboardState
	onReturn      func(*InputBox)
	onlyDigits    bool
}

func NewInputBox(text string, size int, onReturn func(*InputBox)) *InputBox {
	i := &InputBox{
		timerFlashing: NewTimer(500),
		size:          size,
		onReturn:      onReturn,
		onlyDigits:    false,
	}
	i.setupBox(text)
	return i
}

func NewDigitInputBox(text string, size int, onReturn func(*InputBox)) *InputBox {
	i := &InputBox{
		timerFlashing: NewTimer(500),
		size:          size,
		onReturn:      onReturn,
		onlyDigits:    true,
	}
	i.setupBox(text)
	return i
}

func (inp *InputBox) setupBox(text string) {
	inp.SetupView()
	inp.btn = NewButton(text, func(b *Button) {})
	inp._text = ""
	inp.btn.SetText(text)
	theme := GetUi().theme
	inp.bg = theme.Get(InputBoxBg)
	inp.fg = theme.Get(InputBoxFg)
	inp.btn.Fg(inp.fg)
	inp.btn.Bg(inp.bg)
	inp.cursor = NewCursor(inp.bg, inp.fg)
	GetUi().inputKeyboard.Attach(inp)
}

func (inp *InputBox) setPrompt() string {
	str := ""
	for i := inp.size; i > len(inp._text); i-- {
		str += " "
	}
	str += inp._text
	if !inp.cursor.IsVisible() {
		inp.cursor.Visible(true)
	} else {
		inp.cursor.Visible(false)
	}
	return str
}

func (inp *InputBox) UpdateInput(value interface{}) {
	switch vl := value.(type) {
	case KeyboardData:
		if inp.state == ViewStateActive {
			for _, v := range vl.keys {
				if v == ebiten.KeyBackspace {
					inp.keyboardState = KeyBackspace
				} else if v == ebiten.KeyEnter {
					inp.keyboardState = KeyEnter
				} else if v == ebiten.KeyEscape {
					inp.keyboardState = KeyEscape
				} else {
					inp.keyboardState = KeyPressed
					inp.parseInput(vl.GetRunes())
				}
			}
			if len(vl.keys) == 0 {
				inp.keyboardState = KeyReleased
			}
		}
	}
}

func (inp *InputBox) parseInput(chars []rune) {
	if len(inp._text) >= inp.size {
		return
	}
	for _, v := range chars {
		if unicode.IsDigit(v) || v == '.' {
			continue
		} else if inp.onlyDigits {
			inp.btn.Bg(colornames.Red)
			break
		}
	}
	value := string(chars)
	inp._text += value
	inp.btn.SetText(inp.setPrompt())
}

func (inp *InputBox) Update(dt int) {
	inp.btn.Update(dt)
	inp.cursor.Update(dt)
	inp.updatePrompt(dt)
	if inp.state == ViewStateFocus {
		inp.SetState(ViewStateActive)
	}

	if inp.state == ViewStateActive {
		if inp.keyboardState == KeyBackspace {
			if len(inp._text) > 0 {
				inp._text = inp._text[:len(inp._text)-1]
				inp.btn.SetText(inp.setPrompt())
				inp.keyboardState = KeyReleased
				if inp.onlyDigits {
					_, err := strconv.ParseFloat(inp._text, 64)
					if err == nil || len(inp._text) == 0 {
						inp.btn.Bg(inp.bg)
					}
				}
			}
		}
		if inp.keyboardState == KeyEnter {
			if inp.onReturn != nil {
				inp.onReturn(inp)
			}
			inp.keyboardState = KeyReleased
		}
	}
}

func (inp *InputBox) updatePrompt(dt int) {
	inp.timerFlashing.Update(dt)
	if inp.state == ViewStateActive || inp.state == ViewStateFocus {
		if !inp.timerFlashing.IsOn() || inp.timerFlashing.IsDone() {
			inp.timerFlashing.On()
			inp.btn.SetText(inp.setPrompt())
		}
	} else {
		inp.timerFlashing.Off()
		inp.cursor.Visible(false)
	}
}

func (inp *InputBox) GetText() string      { return inp._text }
func (inp *InputBox) SetText(value string) { inp.btn.SetText(value) }

func (inp *InputBox) GetDigit() (float64, error) {
	if len(inp._text) == 0 {
		return 0, nil
	}
	n, err := strconv.ParseFloat(inp._text, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (inp *InputBox) SetDigit(value string) {
	inp._text = value
	inp.btn.SetText(inp.setPrompt())
	inp.Dirty = true
}

func (inp *InputBox) Draw(surface *ebiten.Image) {
	inp.btn.Draw(surface)
	inp.cursor.Draw(surface)
}

func (inp *InputBox) Resize(rect []int) {
	inp.View.Resize(rect)
	sz := inp.size
	w := inp.GetRect().W / sz
	w1 := int(float64(w) * 0.2)
	h := inp.GetRect().H
	x, y := inp.GetRect().Pos()
	inp.btn.Resize([]int{x, y, w * (inp.size), h})
	inp.cursor.Resize([]int{x + w*(inp.size) - w1, y, w1, h})
}

func (inp *InputBox) Close() { GetUi().inputKeyboard.Detach(inp) }
