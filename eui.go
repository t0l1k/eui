package eui

import (
	"image"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var uiInstance *Ui = nil

func init() {
	uiInstance = GetUi()
}

// Инициализация и настройка размеров окна
func Init(f *Ui) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	ebiten.SetWindowTitle(f.title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(f.size.X, f.size.Y)
	log.Println("Set Window Size:", f.size.X, f.size.Y)
	ebiten.SetFullscreen(f.settings.Get(UiFullscreen).(bool))
}

// Переход в вечный цикл...
func Run(sc Scene) {
	GetUi().Push(sc)
	if err := ebiten.RunGame(GetUi()); err != nil {
		log.Fatal(err)
	}
}

// Что-то перед выходом из приложения сделать
func Quit() {}

// Одиночка
func GetUi() (u *Ui) {
	if uiInstance == nil {
		tm := time.Now()
		u = &Ui{
			start:         time.Now(),
			tick:          tm.Nanosecond() / 1e6,
			scenes:        []Scene{},
			inputMouse:    NewMouseInput(),
			inputTouch:    NewTouchInput(),
			inputKeyboard: NewKeyboardInput(),
			theme:         DefaultTheme(),
			settings:      DefaultSettings(),
		}
		log.Printf("App init done")
	} else {
		u = uiInstance
	}
	return u
}

// Умею иницализировать для Ebitenengine приложение, затем запускается сцена в которой отслеживаются через подписку от клавиатуры, мыши, касания экрана события и передаются структуре View(только мышь), InputBox(клавиатуры, пока только цифры). Остальным виджетам это текстовая метка(Text) и изображение(Image) встаивается структура View, которая определяет события от мыши(InputState) и меняет состояние виджета это при наведении, в фокусе, покидание курсором виджета. Все происходит в сцене, где есть обновление и дельта от последнего обновления и рисование, затем уже сцена передает внутри себя виджетам события. После создания сцены по умолчанию создается раскладка по вертикали и при переопределении метода Resize уже в нем производится раскладка виджетов внутри сцены.
type Ui struct {
	title         string
	scenes        []Scene
	currentScene  Scene
	theme         *Theme
	settings      *Setting
	tick          int
	start         time.Time
	size          image.Point
	inputMouse    *MouseInput
	inputTouch    *TouchInput
	inputKeyboard *KeyboardInput
}

func (u *Ui) GetStartTime() time.Time {
	return u.start
}

func (u *Ui) GetInputTouch() *TouchInput {
	return u.inputTouch
}

func (u *Ui) GetInputMouse() *MouseInput {
	return u.inputMouse
}

func (u *Ui) GetInputKeyboard() *KeyboardInput {
	return u.inputKeyboard
}

func (u *Ui) GetTitle() string {
	return u.title
}

func (u *Ui) SetTitle(value string) {
	u.title = value
}

func (u *Ui) SetFullscreen(value bool) {
	u.settings.Set(UiFullscreen, value)
}

func (u *Ui) SetSize(w, h int) {
	u.size.X = w
	u.size.Y = h
}

func (u *Ui) IsMainScene() bool {
	return len(u.scenes) == 0

}

func (u *Ui) Size() (int, int) {
	return u.size.X, u.size.Y
}

func (u *Ui) GetTheme() *Theme {
	return u.theme
}

func (u *Ui) GetSettings() *Setting {
	return u.settings
}

// Отсюда можно следить за изменением размера окна, при изменении обновляются размеры текущей сцены
func (u *Ui) Layout(w, h int) (int, int) {
	if w != u.size.X || h != u.size.Y {
		u.size.X, u.size.Y = w, h
		u.currentScene.Resize()
		for _, scene := range u.scenes {
			scene.Resize()
		}
		log.Println("Resize app done, new size:", w, h)
	}
	return w, h
}

func (u *Ui) Update() error {
	tick := u.getTick()
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) || ebiten.IsWindowBeingClosed() {
		err := u.Pop()
		if err != nil {
			os.Exit(0)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF11) {
		u.ToggleFullscreen()
	}
	u.inputMouse.update(tick)
	u.inputTouch.update(tick)
	u.inputKeyboard.update(tick)
	u.currentScene.Update(tick)
	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	screen.Fill(u.theme.Get(SceneBg))
	u.currentScene.Draw(screen)
}

func (a *Ui) ToggleFullscreen() {
	fullscreen := a.settings.Get(UiFullscreen).(bool)
	fullscreen = !fullscreen
	ebiten.SetFullscreen(fullscreen)
	a.settings.Set(UiFullscreen, fullscreen)
	log.Println("Toggle FullScreen", a.size)
}

// Определяю дельту от последнего обновления в миллисекундах
func (u *Ui) getTick() (ticks int) {
	tm := time.Now()
	dt := tm.Nanosecond() / 1e6
	if dt < u.tick {
		ticks = 999 - u.tick + dt
	} else {
		ticks = dt - u.tick
	}
	u.tick = dt
	return ticks
}

// Добавить сцену и сделать текущей
func (u *Ui) Push(sc Scene) {
	u.scenes = append(u.scenes, sc)
	u.currentScene = sc
	u.currentScene.Entered()
	log.Println("Scene push", u.scenes)
}

// Закрыть текущую сцену если первая выход, иначе последнюю сделать текущей
func (u *Ui) Pop() error {
	if len(u.scenes) > 0 {
		u.currentScene.Quit()
		idx := len(GetUi().scenes) - 1
		u.scenes = GetUi().scenes[:idx]
		log.Println("App Pop Scene Quit done.", u.scenes)
		if u.IsMainScene() {
			log.Printf("App Quit.")
			os.Exit(0)
		}
		u.currentScene = GetUi().scenes[len(GetUi().scenes)-1]
		u.currentScene.Entered()
		log.Println("App Pop New Scene Entered.", u.scenes)
	}
	return nil
}
