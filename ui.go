package eui

import (
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Умею иницализировать для Ebitenengine приложение, затем запускается сцена в которой отслеживаются через подписку от клавиатуры, мыши, касания экрана события и передаются структуре View, InputBox(клавиатуры, пока только цифры). Остальным виджетам это текстовая метка(Text) и изображение(Image) встаивается структура View, которая определяет события от мыши(InputState) и меняет состояние виджета это при наведении, в фокусе, покидание курсором виджета. Все происходит в сцене, где есть обновление и дельта от последнего обновления и рисование, затем уже сцена передает внутри себя виджетам события. После создания сцены по умолчанию создается раскладка по вертикали и при переопределении метода Resize уже в нем производится раскладка виджетов внутри сцены.
type Ui struct {
	title          string
	scenes         []Sceneer
	currentScene   Sceneer
	theme          *Theme
	settings       *Setting
	tick           int
	start          time.Time
	size           PointInt
	inputMouse     *MouseInput
	inputTouch     *TouchInput
	inputKeyboard  *KeyboardInput
	resizeListener *ResizeListener
	tickListener   *TickListener
	modal          Drawabler
}

func (u *Ui) GetStartTime() time.Time          { return u.start }
func (u *Ui) GetInputTouch() *TouchInput       { return u.inputTouch }
func (u *Ui) GetInputMouse() *MouseInput       { return u.inputMouse }
func (u *Ui) GetInputKeyboard() *KeyboardInput { return u.inputKeyboard }
func (u *Ui) GetTitle() string                 { return u.title }
func (u *Ui) SetTitle(value string) *Ui        { u.title = value; return u }
func (u *Ui) SetFullscreen(value bool)         { u.settings.Set(UiFullscreen, value) }
func (u *Ui) Size() (int, int)                 { return u.size.X, u.size.Y }
func (u *Ui) SetSize(w, h int) *Ui             { u.size = NewPointInt(w, h); return u }
func (u *Ui) IsMainScene() bool                { return len(u.scenes) == 0 }
func (u *Ui) GetTheme() *Theme                 { return u.theme }
func (u *Ui) SetTheme(value *Theme) *Ui        { u.theme = value; return u }
func (u *Ui) GetSettings() *Setting            { return u.settings }

// Отсюда можно следить за изменением размера окна, при изменении обновляются размеры текущей сцены
func (u *Ui) Layout(w, h int) (int, int) {
	if w != u.size.X || h != u.size.Y {
		u.resizeListener.Emit(NewEvent(EventResize, NewRect([]int{0, 0, w, h})))
		log.Println("Emit:Resize:")
	}
	return w, h
}

func (u *Ui) HandleEvent(ev Event) {
	switch ev.Type {
	case EventTick:
		// tc := ev.Value.(TickData)
		// log.Println("Ui:HandleEvent:Tick", tc.String())
	case EventResize:
		r := ev.Value.(Rect)
		u.SetSize(r.W, r.H)
		u.currentScene.Resize()
		for _, scene := range u.scenes {
			scene.Resize()
		}
		log.Println("Resize app done, new size:", r)
	case EventKeyReleased:
		kd := ev.Value.(KeyboardData)
		if kd.IsReleased(ebiten.KeyF12) {
			u.ToggleFullscreen()
		}
		if kd.IsReleased(ebiten.KeyEscape) {
			err := u.Pop()
			if err != nil {
				os.Exit(0)
			}
		}
	}
	if !(ev.Type == EventTick) {
		log.Println("Ui:HandleEvent:", ev)
	}
}

func (u *Ui) Update() error {
	tick := u.getTick()
	u.tickListener.update(tick)
	u.inputMouse.update(tick)
	u.inputTouch.update(tick)
	u.inputKeyboard.update(tick)
	u.currentScene.Update(tick)
	if u.modal != nil {
		u.modal.Update(tick)
	}
	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	screen.Fill(u.theme.Get(SceneBg))
	u.currentScene.Draw(screen)
	if u.modal != nil && u.modal.IsVisible() {
		u.modal.Draw(screen)
	}
}

func (u *Ui) ShowModal(d Drawabler) {
	u.modal = d
	u.modal.Visible(true)
}

func (u *Ui) HideModal() {
	if u.modal != nil {
		u.modal.Visible(false)
		u.modal.Close()
		u.modal = nil
	}
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
func (u *Ui) Push(sc Sceneer) {
	u.scenes = append(u.scenes, sc)
	u.currentScene = sc
	u.resizeListener.Emit(NewEvent(EventResize, NewRect([]int{0, 0, u.size.X, u.size.Y})))
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
