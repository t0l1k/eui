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
	lbl    *eui.Text
	btn    *eui.Button
	active bool
	Value  ebiten.Key
	kId    int64
}

func NewInputKey(title string) *InputKey {
	i := &InputKey{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	i.lbl = eui.NewText(title)
	i.Add(i.lbl)
	i.btn = eui.NewButton("(?)", func(b *eui.Button) {
		if b.IsPressed() {
			i.active = true
			i.btn.Bg(colornames.Yellow)
		}
	})
	i.Add(i.btn)
	i.kId = eui.GetUi().GetInputKeyboard().Connect(i.UpdateInput)
	return i
}

func (i *InputKey) UpdateInput(ev eui.Event) {
	kd := ev.Value.(eui.KeyboardData)
	if i.active {
		if ev.Type == eui.EventKeyReleased {
			i.btn.SetText(kd.GetKeysReleased()[0].String())
			i.Value = kd.GetKeysReleased()[0]
		}
	}
}

func (i *InputKey) Update(dt int) {
	i.Container.Update(dt)
	if i.btn.GetState() == eui.ViewStateNormal {
		i.active = false
		i.btn.Bg(colornames.Silver)
	}
}

func (i *InputKey) Resize(rect eui.Rect) {
	i.SetRect(rect)
	w0, h0 := i.Rect().Size()
	x0, y0 := i.Rect().Pos()
	i.btn.Resize(eui.NewRect([]int{x0, y0, h0 * 2, h0}))
	i.lbl.Resize(eui.NewRect([]int{x0 + h0*2, y0, w0 - h0*2, h0}))
	i.ImageReset()
}

func (i *InputKey) String() string {
	return fmt.Sprintf("%v: %v", i.lbl.GetText(), i.Value)
}

func (i *InputKey) Close() { eui.GetUi().GetInputKeyboard().Disconnect(i.kId) }

type HotkeyDialog struct {
	*eui.Container
	layV    *eui.Container
	title   *eui.Text
	btnHide *eui.Button
}

func NewHotkeyDialog() *HotkeyDialog {
	d := &HotkeyDialog{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	d.title = eui.NewText(title)
	d.Add(d.title)
	d.btnHide = eui.NewButton("X", func(b *eui.Button) {
		d.Visible(false)
		for _, v := range d.layV.Childrens() {
			log.Println("hide dialog, result:", v.(*InputKey).String())
		}
	})
	d.Add(d.btnHide)
	d.layV = eui.NewContainer(eui.NewVBoxLayout(1))
	str := []string{"Выбрать управление движения влево", "Выбрать управление движения вправо", "Выбрать управление движения вверх", "Выбрать управление движения вниз"}
	for _, v := range str {
		btn := NewInputKey(v)
		d.layV.Add(btn)
	}
	d.Add(d.layV)
	return d
}

func (d *HotkeyDialog) Visible(value bool) {
	d.Traverse(func(c eui.Drawabler) { c.Visible(value) }, false)
}

func (d *HotkeyDialog) Resize(rect eui.Rect) {
	d.SetRect(rect)
	w0, h0 := d.Rect().Size()
	x0, y0 := d.Rect().Pos()
	hTop := int(float64(h0) * 0.1) // topbar height
	d.title.Resize(eui.NewRect([]int{x0, y0, w0 - hTop, hTop}))
	d.btnHide.Resize(eui.NewRect([]int{x0 + w0 - hTop, y0, hTop, hTop}))
	x, y := x0, y0+hTop
	w, h := w0, h0-hTop
	d.layV.Resize(eui.NewRect([]int{x, y, w, h}))
	d.ImageReset()
}

type SceneSelectHotkey struct {
	*eui.Scene
	dialog *HotkeyDialog
	topBar *eui.TopBar
}

func NewSceneSelectHotkey() *SceneSelectHotkey {
	s := &SceneSelectHotkey{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar(title, func(b *eui.Button) {
		s.dialog.Visible(true)
		log.Println("selrct menu: dialog visible")
	})
	s.Add(s.topBar)
	s.dialog = NewHotkeyDialog()
	s.Add(s.dialog)
	return s
}

func (s *SceneSelectHotkey) Resize() {
	w0, h0 := eui.GetUi().Size()
	hTop := int(float64(h0) * 0.05) // topbar height
	s.topBar.Resize(eui.NewRect([]int{0, 0, w0, hTop}))
	s.dialog.Resize(eui.NewRect([]int{hTop, hTop * 2, w0 - hTop*2, h0 - hTop*3}))
}

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle(title)
	k := 2
	w, h := 320*k, 200*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneSelectHotkey())
	eui.Quit(func() {})
}
