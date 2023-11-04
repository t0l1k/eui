package eui

type InputState string

const (
	ViewStateNormal   = "normal"  // Виджет без фокуса
	ViewStateHover    = "hovered" // Виджет под мышью
	ViewStateFocus    = "focused" // Виджет нажата кнопка мышь или касание экрана
	ViewStateExec     = "exec"    // Виджет исполнить метод после отпускания касания экрана
	ViewStateActive   = "active"  // Виджет активен для ввода от клавиатуры
	ViewStateSelected = "selected"
	ViewStateDisabled = "disabled" // Виджет неактивен
)

func (v InputState) String() string {
	return string(v)
}

type InputKeyboardState string

const (
	KeyPressed   = "key pressed"
	KeyReleased  = "key released"
	KeyBackspace = "backspace" // Удалить символ
	KeyEnter     = "return"    // Выполнить метод
	KeyEscape    = "escape"
)
