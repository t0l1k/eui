package eui

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

// Умею получить от клавиатуры нажатый символ(пока только цифры и точка), backspace удаляет последний введенный символ, enter запускает прикрепленный метод. Умею мигать курсором. Только при фокусе от клавиатуры доступен ввод, активированый нажатием левой кнопки мыши(меняется виджет обрамление былым).
type InputBox struct {
	_text string
	Button
	size          int
	prompt        string
	isPrompt      bool
	timerFlashing *Timer
	keyboardState InputKeyboardState
	onReturn      func(*InputBox)
}

func NewInputBox(text string, size int, bg, fg color.Color, onReturn func(*InputBox)) *InputBox {
	i := &InputBox{
		timerFlashing: NewTimer(500),
		size:          size,
		prompt:        "|",
		onReturn:      onReturn,
	}
	i.setupBox(text, bg, fg)
	return i
}

func (inp *InputBox) setupBox(text string, bg, fg color.Color) {
	inp._text = text
	inp.SetupText(text, bg, fg)
	inp.setPrompt()
	GetUi().inputKeyboard.Attach(inp)
}

func (inp *InputBox) setPrompt() string {
	str := ""
	if len(inp._text) < inp.size {
		for i := inp.size; i > len(inp._text); i-- {
			str += " "
		}
		str += inp._text
		if inp.isPrompt {
			str += inp.prompt
		} else {
			str += " "
		}
	}
	inp.isPrompt = !inp.isPrompt
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
					fmt.Println(v)
					inp.inputNumbers(v)
				}
				fmt.Println("Pressed", v)
			}
			if len(vl.keys) == 0 {
				inp.keyboardState = KeyReleased
			}
		}
	}
}

func (inp *InputBox) inputNumbers(key ebiten.Key) {
	switch key {
	case ebiten.Key0:
		inp._text += "0"
	case ebiten.Key1:
		inp._text += "1"
	case ebiten.Key2:
		inp._text += "2"
	case ebiten.Key3:
		inp._text += "3"
	case ebiten.Key4:
		inp._text += "4"
	case ebiten.Key5:
		inp._text += "5"
	case ebiten.Key6:
		inp._text += "6"
	case ebiten.Key7:
		inp._text += "7"
	case ebiten.Key8:
		inp._text += "8"
	case ebiten.Key9:
		inp._text += "9"
	case ebiten.KeyComma:
		inp._text += ","
	case ebiten.KeyPeriod:
		inp._text += "."
	}
	inp.SetText(inp.setPrompt())
}

func (inp *InputBox) Update(dt int) {
	inp.updatePrompt(dt)
	if inp.state == ViewStateFocus {
		inp.SetState(ViewStateActive)
	}

	if inp.state == ViewStateActive {
		if inp.keyboardState == KeyBackspace {
			if len(inp._text) > 0 {
				inp._text = inp._text[:len(inp._text)-1]
				inp.SetText(inp.setPrompt())
				inp.keyboardState = KeyReleased

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
		if !inp.timerFlashing.IsOn() {
			inp.timerFlashing.On()
			inp.SetText(inp.setPrompt())
		}
		if inp.timerFlashing.IsDone() {
			inp.timerFlashing.On()
			inp.SetText(inp.setPrompt())
		}
	} else {
		inp.timerFlashing.Off()
		inp.isPrompt = false
	}
}

func (inp *InputBox) GetDigit() float64 {
	n, err := strconv.ParseFloat(inp._text, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func (inp *InputBox) SetDigit(value string) {
	inp._text = value
	inp.Text.SetText(inp.setPrompt())
	inp.dirty = true
}
