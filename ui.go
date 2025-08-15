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
	size           Point[int]
	inputMouse     *MouseListener
	inputKeyboard  *KeyboardInput
	resizeListener *ResizeListener
	tickListener   *TickListener
	modal          Drawabler
	focusManager   *FocusManager
	resource       *ResourceManager
}

func (u *Ui) GetStartTime() time.Time          { return u.start }
func (u *Ui) GetInputMouse() *MouseListener    { return u.inputMouse }
func (u *Ui) GetInputKeyboard() *KeyboardInput { return u.inputKeyboard }
func (u *Ui) GetTitle() string                 { return u.title }
func (u *Ui) SetTitle(value string) *Ui        { u.title = value; return u }
func (u *Ui) SetFullscreen(value bool)         { u.settings.Set(UiFullscreen, value) }
func (u *Ui) Size() (int, int)                 { return u.size.X, u.size.Y }
func (u *Ui) SetSize(w, h int) *Ui             { u.size = NewPoint(w, h); return u }
func (u *Ui) IsMainScene() bool                { return len(u.scenes) == 0 }
func (u *Ui) GetTheme() *Theme                 { return u.theme }
func (u *Ui) SetTheme(value *Theme) *Ui        { u.theme = value; return u }
func (u *Ui) GetSettings() *Setting            { return u.settings }
func (u *Ui) FontDefault() *Font               { return u.resource.FontDefault() }

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
		tc := ev.Value.(TickData)
		u.currentScene.Traverse(func(d Drawabler) {
			if t, ok := d.(interface{ Tick(TickData) }); ok {
				t.Tick(tc)
			}
		}, false)
	case EventResize:
		r := ev.Value.(Rect[int])
		u.SetSize(r.W, r.H)
		u.currentScene.SetRect(r)
		for _, scene := range u.scenes {
			scene.SetRect(r)
		}
		log.Println("Resize app done, new size:", r)
	}
	u.HandleKeybordEvent(ev)
	u.HandleMouseEvent(ev)
	if !(ev.Type == EventTick || ev.Type == EventMouseMovement) {
		log.Println("Ui:HandleEvent:", ev)
	}
}

func (u *Ui) HandleKeybordEvent(ev Event) {
	switch ev.Type {
	case EventKeyPressed:
		kd := ev.Value.(KeyboardData)
		u.currentScene.Traverse(func(d Drawabler) {
			if kh, ok := d.(interface{ KeyPressed(KeyboardData) }); ok {
				kh.KeyPressed(kd)
			}
		}, false)
		d := u.currentScene
		if kh, ok := d.(interface{ KeyPressed(KeyboardData) }); ok {
			kh.KeyPressed(kd)
		}
		log.Println("Ui:HandleEvent:EventKeyPressed", ev)
	case EventKeyReleased:
		kd := ev.Value.(KeyboardData)
		if kd.IsReleased(ebiten.KeyF12) {
			u.ToggleFullscreen()
		}
		u.currentScene.Traverse(func(d Drawabler) {
			if kh, ok := d.(interface{ KeyReleased(KeyboardData) }); ok {
				kh.KeyReleased(kd)
			}
		}, false)
		d := u.currentScene
		if kh, ok := d.(interface{ KeyReleased(KeyboardData) }); ok {
			kh.KeyReleased(kd)
		}
		log.Println("Ui:HandleEvent:EventKeyReleased", ev)
	}
}

func (u *Ui) HandleMouseEvent(ev Event) {
	switch ev.Type {
	case EventMouseDown:
		var pressed Drawabler
		md := ev.Value.(MouseData)
		u.currentScene.Traverse(func(d Drawabler) {
			if mh, ok := d.(interface{ Hit(Point[int]) Drawabler }); ok {
				if mh.Hit(md.pos) != nil && pressed == nil {
					pressed = d
					return
				}
			}
		}, false)
		if pressed != nil {
			u.focusManager.SetFocused(pressed)
			if mp, ok := pressed.(interface{ MouseDown(MouseData) }); ok {
				mp.MouseDown(md)
			}
			log.Println("Ui:HandleMouseEvent:01", md)
		} else {
			u.focusManager.Blur()
		}
	case EventMouseUp:
		md := ev.Value.(MouseData)
		if u.focusManager.Focused() != nil {
			d := u.focusManager.Focused()
			if mp, ok := d.(interface{ MouseUp(MouseData) }); ok {
				mp.MouseUp(md)
			}
			if mp, ok := d.(interface{ WantBlur() bool }); ok {
				if mp.WantBlur() {
					u.focusManager.Blur()
				}
			}
			log.Println("Ui:HandleMouseEvent:02", md)
		}
	case EventMouseWheel:
		var wheel Drawabler
		md := ev.Value.(MouseData)
		u.currentScene.Traverse(func(d Drawabler) {
			if mh, ok := d.(interface{ Hit(Point[int]) Drawabler }); ok {
				if mh.Hit(md.pos) != nil && wheel == nil {
					wheel = d
				}
			}

		}, false)
		if wheel != nil {
			if m, ok := wheel.(interface{ MouseWheel(MouseData) }); ok {
				m.MouseWheel(md)
			}
			log.Println("Ui:HandleMouseEvent:04", wheel.Rect())
		}
	case EventMouseMovement, EventMouseDrag:
		var hovered Drawabler
		md := ev.Value.(MouseData)
		u.currentScene.Traverse(func(d Drawabler) {
			if mh, ok := d.(interface{ Hit(Point[int]) Drawabler }); ok {
				if mh.Hit(md.pos) != nil && hovered == nil {
					hovered = d
				}
			}
		}, false)
		if u.focusManager.Focused() == nil {
			if u.focusManager.Hovered() != hovered {
				if u.focusManager.Hovered() != nil {
					if m, ok := hovered.(interface{ MouseLeave() }); ok {
						m.MouseLeave()
					}
					log.Println("Ui:HandleMouseEvent:05")
				}
				if hovered != nil {
					if m, ok := hovered.(interface{ MouseEnter() }); ok {
						m.MouseEnter()
					}
					log.Println("Ui:HandleMouseEvent:06", hovered.Rect())
				}
				u.focusManager.SetHovered(hovered)
			}
		}
		if u.focusManager.Hovered() != nil || u.focusManager.Focused() != nil {
			if hovered != nil {
				if ev.Type == EventMouseDrag {
					if m, ok := hovered.(interface{ MouseDrag(MouseData) }); ok {
						m.MouseDrag(md)
					}
					log.Println("Ui:HandleMouseEvent:07", hovered.Rect())
				} else {
					if m, ok := hovered.(interface{ MouseMotion(MouseData) }); ok {
						m.MouseMotion(md)
					}
					log.Println("Ui:HandleMouseEvent:08", hovered.Rect())
				}
			}
		}
	}
}

func (u *Ui) Update() error {
	tick := u.getTick()
	u.tickListener.update()
	u.inputMouse.update()
	u.inputKeyboard.update()
	u.currentScene.Update(tick)
	if u.modal != nil {
		u.modal.Update(tick)
	}
	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	screen.Fill(u.theme.Get(SceneBg))
	u.currentScene.Draw(screen)
	if u.modal != nil && !u.modal.IsHidden() {
		u.modal.Draw(screen)
	}
}

func (u *Ui) ShowModal(d Drawabler) {
	u.modal = d
	u.modal.Show()
}

func (u *Ui) HideModal() {
	if u.modal != nil {
		u.modal.Hide()
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
