package app

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

var actsStr = []string{"Отменить", "Удалить", "Заметка"}

type BottomBar struct {
	eui.DrawableBase
	layoutActs, layoutNums *eui.BoxLayout
	dim                    int
	show                   bool
	fn                     func(*eui.Button)
}

func NewBottomBar(fn func(*eui.Button)) *BottomBar {
	b := &BottomBar{}
	b.fn = fn
	b.layoutActs = eui.NewHLayout()
	b.layoutNums = eui.NewHLayout()
	b.Bg(eui.Red)
	b.Fg(eui.Silver)
	return b
}

func (b *BottomBar) Setup(dim int) {
	b.dim = dim
	b.layoutActs.ResetContainerBase()
	b.layoutNums.ResetContainerBase()
	for _, v := range actsStr {
		btn := eui.NewButton(v, b.fn)
		b.layoutActs.Add(btn)
	}
	size := dim * dim
	for i := 0; i < size; i++ {
		btn := eui.NewButton(strconv.Itoa(i+1), b.fn)
		b.layoutNums.Add(btn)
	}
	b.Resize(b.GetRect().GetArr())
}

func (b *BottomBar) IsVisible() bool    { return b.show }
func (b *BottomBar) Visible(value bool) { b.show = value }

func (b *BottomBar) Update(dt int) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layoutActs.GetContainer() {
		v.Update(dt)
	}
	for _, v := range b.layoutNums.GetContainer() {
		v.Update(dt)
	}
}

func (b *BottomBar) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layoutActs.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range b.layoutNums.GetContainer() {
		v.Draw(surface)
	}
}

func (b *BottomBar) Resize(rect []int) {
	b.Rect(eui.NewRect(rect))
	b.SpriteBase.Resize(rect)
	w, h := b.GetRect().Size()
	x, y := b.GetRect().Pos()
	b.layoutActs.Resize([]int{x, y, w, h / 2})
	y += h / 2
	b.layoutNums.Resize([]int{x, y, w, h / 2})
	b.ImageReset()
}
