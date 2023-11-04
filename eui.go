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
	ebiten.SetFullscreen(false)
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
	tick          int
	start         time.Time
	size          image.Point
	inputMouse    *MouseInput
	inputTouch    *TouchInput
	inputKeyboard *KeyboardInput
	fullScreen    bool
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
	u.fullScreen = value
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

// Отсюда можно следить за изменением размера окна, при изменении обновляются размеры текущей сцены
func (u *Ui) Layout(w, h int) (int, int) {
	if w != u.size.X || h != u.size.Y {
		u.size.X, u.size.Y = w, h
		u.currentScene.Resize()
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
	u.currentScene.Draw(screen)
}

func (a *Ui) ToggleFullscreen() {
	a.fullScreen = !a.fullScreen
	ebiten.SetFullscreen(a.fullScreen)
	for _, scene := range a.scenes {
		scene.Resize()
	}
	log.Println("Toggle FullScreen")
}

// Определяю дельту от последнего обновления в миллисекундах
func (u *Ui) getTick() int {
	tm := time.Now()
	dt := tm.Nanosecond() / 1e6
	ticks := dt - u.tick
	if dt < u.tick {
		ticks = 999 - u.tick + dt
	}
	u.tick = dt
	return ticks
}

// Добавить сцену и сделать текущей
func (u *Ui) Push(sc Scene) {
	u.scenes = append(u.scenes, sc)
	u.currentScene = sc
	u.currentScene.Entered()
	log.Println("Scene push")
}

// Закрыть текущую сцену если первая выход, иначе последнюю сделать текущей
func (u *Ui) Pop() error {
	if len(u.scenes) > 0 {
		u.currentScene.Quit()
		idx := len(GetUi().scenes) - 1
		u.scenes = GetUi().scenes[:idx]
		log.Printf("App Pop Scene Quit done.")
		if u.IsMainScene() {
			log.Printf("App Quit.")
			os.Exit(0)
		}
		u.currentScene = GetUi().scenes[len(GetUi().scenes)-1]
		u.currentScene.Entered()
		log.Printf("App Pop New Scene Entered.")
	}
	return nil
}
