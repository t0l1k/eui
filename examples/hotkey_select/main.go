package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

const (
	title = "Выбрать клавишы управления героем"
)

type InputKey struct {
	*eui.Container
	lbl    *eui.Label
	btn    *eui.Button
	active bool
	key    *eui.Signal[ebiten.Key]
	kId    int64
}

func NewInputKey(title string) *InputKey {
	i := &InputKey{Container: eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{70, 30}, 5)), key: eui.NewSignal(func(a, b ebiten.Key) bool { return a == b })}
	i.lbl = eui.NewLabel(title)
	i.Add(i.lbl)
	i.btn = eui.NewButton("(?)", func(b *eui.Button) {
		i.active = true
		i.btn.Bg(colornames.Yellow)
	})
	i.Add(i.btn)
	i.kId = i.key.Connect(func(key ebiten.Key) {
		i.active = false
		i.btn.SetText(key.String())
		i.btn.Bg(colornames.Navy)
		i.btn.Fg(colornames.Yellow)
	})
	return i
}

func (i *InputKey) KeyPressed(kd eui.KeyboardData) {
	if !i.active {
		return
	}
	i.key.Emit(kd.GetKeysPressed()[0])
	log.Println("InputKey:KeyPressed:", kd)
}
func (i *InputKey) String() string {
	return fmt.Sprintf("%v: %v", i.lbl.Text(), i.key.Value().String())
}
func (i *InputKey) Close() { i.key.Disconnect(i.kId) }

type HotkeyDialog struct{ *eui.Container }

func NewHotkeyDialog() *HotkeyDialog {
	d := &HotkeyDialog{Container: eui.NewContainer(eui.NewLayoutVerticalPercent([]int{10, 90}, 5))}
	contInputLines := eui.NewContainer(eui.NewVBoxLayout(10))
	contTitle := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{90, 10}, 5))
	contTitle.Add(eui.NewLabel(title))
	contTitle.Add(eui.NewButton("X", func(b *eui.Button) {
		d.Hide()
		for _, v := range contInputLines.Children() {
			log.Println("hide dialog, result:", v.(*InputKey).String())
		}
	}))
	d.Add(contTitle)
	str := []string{"Выбрать управление движения влево", "Выбрать управление движения вправо", "Выбрать управление движения вверх", "Выбрать управление движения вниз"}
	for _, v := range str {
		btn := NewInputKey(v)
		contInputLines.Add(btn)
	}
	d.Add(contInputLines)
	return d
}

func main() {
	eui.Init(func() *eui.Ui {
		u := eui.GetUi()
		u.SetTitle(title)
		k := 2
		w, h := 320*k, 200*k
		u.SetSize(w, h)
		return u
	}())
	eui.Run(func() *eui.Scene {
		s := eui.NewScene(eui.NewLayoutVerticalPercent([]int{10, 90}, 5))
		dialog := NewHotkeyDialog()
		s.Add(eui.NewTopBar(title, func(b *eui.Button) {
			dialog.Show()
			log.Println("select menu: dialog visible")
		}))
		s.Add(dialog)
		return s
	}())
	eui.Quit(func() {})
}
