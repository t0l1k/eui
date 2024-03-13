package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

const (
	title = "Выбрать клавишы управления героем"
)

type InputKey struct {
	eui.DrawableBase
	lbl    *eui.Text
	btn    *eui.Button
	active bool
	Value  ebiten.Key
}

func NewInputKey(title string) *InputKey {
	i := &InputKey{}
	i.lbl = eui.NewText(title)
	i.Add(i.lbl)
	i.btn = eui.NewButton("(?)", func(b *eui.Button) {
		if b.IsPressed() {
			i.active = true
			i.btn.Bg(eui.Yellow)
		}
	})
	i.Add(i.btn)
	eui.GetUi().GetInputKeyboard().Attach(i)
	return i
}

func (i *InputKey) UpdateInput(value interface{}) {
	switch v := value.(type) {
	case eui.KeyboardData:
		if i.active {
			i.btn.SetText(v.GetKeys()[0].String())
			i.Value = v.GetKeys()[0]
		}
	}
}

func (i *InputKey) Update(dt int) {
	i.DrawableBase.Update(dt)
	if i.btn.GetState() == eui.ViewStateNormal {
		i.active = false
		i.btn.Bg(eui.Silver)
	}
}

func (i *InputKey) Resize(rect []int) {
	i.Rect(eui.NewRect(rect))
	w0, h0 := i.GetRect().Size()
	x0, y0 := i.GetRect().Pos()
	i.btn.Resize([]int{x0, y0, h0 * 2, h0})
	i.lbl.Resize([]int{x0 + h0*2, y0, w0 - h0*2, h0})
	i.ImageReset()
}

func (i *InputKey) String() string {
	return fmt.Sprintf("%v: %v", i.lbl.GetText(), i.Value)
}

func (i *InputKey) Close() { eui.GetUi().GetInputKeyboard().Detach(i) }

type HotkeyDialog struct {
	eui.DrawableBase
	layV    *eui.BoxLayout
	title   *eui.Text
	btnHide *eui.Button
}

func NewHotkeyDialog() *HotkeyDialog {
	d := &HotkeyDialog{}
	d.title = eui.NewText(title)
	d.Add(d.title)
	d.btnHide = eui.NewButton("X", func(b *eui.Button) {
		d.Visible = false
		for _, v := range d.layV.GetContainer() {
			log.Println("hide dialog, result:", v.(*InputKey).String())
		}
	})
	d.Add(d.btnHide)
	d.layV = eui.NewVLayout()
	str := []string{"Выбрать управление движения влево", "Выбрать управление движения вправо", "Выбрать управление движения вверх", "Выбрать управление движения вниз"}
	for _, v := range str {
		btn := NewInputKey(v)
		d.layV.Add(btn)
	}
	return d
}

func (d *HotkeyDialog) Update(dt int) {
	if !d.Visible {
		return
	}
	d.DrawableBase.Update(dt)
	for _, v := range d.layV.GetContainer() {
		v.Update(dt)
	}
}

func (d *HotkeyDialog) Draw(surface *ebiten.Image) {
	if !d.Visible {
		return
	}
	d.DrawableBase.Draw(surface)
	for _, v := range d.layV.GetContainer() {
		v.Draw(surface)
	}
}

func (d *HotkeyDialog) Resize(rect []int) {
	d.Rect(eui.NewRect(rect))
	w0, h0 := d.GetRect().Size()
	x0, y0 := d.GetRect().Pos()
	hTop := int(float64(h0) * 0.1) // topbar height
	d.title.Resize([]int{x0, y0, w0 - hTop, hTop})
	d.btnHide.Resize([]int{x0 + w0 - hTop, y0, hTop, hTop})
	x, y := x0, y0+hTop
	w, h := w0, h0-hTop
	d.layV.Resize([]int{x, y, w, h})
	d.ImageReset()
}

type SceneSelectHotkey struct {
	eui.SceneBase
	dialog *HotkeyDialog
	topBar *eui.TopBar
}

func NewSceneSelectHotkey() *SceneSelectHotkey {
	s := &SceneSelectHotkey{}
	s.topBar = eui.NewTopBar(title, func(b *eui.Button) {
		s.dialog.Visible = true
		log.Println("selrct menu: dialog visible")
	})
	s.Add(s.topBar)
	s.dialog = NewHotkeyDialog()
	s.Add(s.dialog)
	s.Resize()
	return s
}

func (s *SceneSelectHotkey) Resize() {
	w0, h0 := eui.GetUi().Size()
	hTop := int(float64(h0) * 0.05) // topbar height
	s.topBar.Resize([]int{0, 0, w0, hTop})
	s.dialog.Resize([]int{hTop, hTop * 2, w0 - hTop*2, h0 - hTop*3})
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
	eui.Quit()
}
