package eui

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui/res"
)

var uiInstance *Ui = nil

func init() { uiInstance = GetUi() }

// Инициализация и настройка размеров окна
func Init(u *Ui) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	ebiten.SetWindowTitle(u.title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(u.size.X, u.size.Y)
	log.Println("Set Window Size:", u.size.X, u.size.Y)
	ebiten.SetFullscreen(u.settings.Get(UiFullscreen).(bool))
	u.tickListener = NewTickListener(u.HandleEvent, 10*time.Millisecond)
	u.inputKeyboard = NewKeyboardListener(u.HandleEvent)
	u.inputMouse = NewMouseListener(u.HandleEvent)
	u.resizeListener = NewResizeListener(u.HandleEvent)
	u.focusManager = &FocusManager{}
	u.resource = NewResourceManager().LoadFont(FontDefault, res.DejaVuSans_ttf, 36)
}

// Переход в вечный цикл...
func Run(sc Sceneer) {
	GetUi().Push(sc)
	if err := ebiten.RunGame(GetUi()); err != nil {
		log.Fatal(err)
	}
}

// Что-то перед выходом из приложения сделать
func Quit(fn func()) {
	if fn != nil {
		fn()
	}
}

// Одиночка
func GetUi() (u *Ui) {
	if uiInstance == nil {
		tm := time.Now()
		u = &Ui{
			start:    time.Now(),
			tick:     tm.Nanosecond() / 1e6,
			scenes:   []Sceneer{},
			theme:    DefaultTheme(),
			settings: DefaultSettings(),
		}
		log.Printf("Ui instance init done")
	} else {
		u = uiInstance
	}
	return u
}
