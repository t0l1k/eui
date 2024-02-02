package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/lines/game"
)

type Table struct {
	eui.View
	leftLbl, rightLbl *eui.Text
	nextBallsLayout   *eui.BoxLayout
}

func NewTable() *Table {
	t := &Table{}
	t.SetupView()
	t.leftLbl = eui.NewText("0")
	t.Add(t.leftLbl)
	t.rightLbl = eui.NewText("100")
	t.Add(t.rightLbl)
	t.nextBallsLayout = eui.NewHLayout()
	t.Add(t.nextBallsLayout)
	return t
}

func (t *Table) Setup(balls int) {
	t.nextBallsLayout.Container = nil
	for i := 0; i < balls; i++ {
		icon := NewBallIcon(BallHidden, game.BallNoColor.Color(), game.BallNoColor.Color())
		x, y, w, h := t.nextBallsLayout.Rect.GetRect()
		icon.Resize([]int{x, y, w / 3, h})
		t.nextBallsLayout.Add(eui.NewIcon(icon.GetImage()))
	}
	t.Resize(t.GetRect().GetArr())
}

func (t *Table) SetNextMoveBalls(cells []*game.Cell) {
	var bg, fg color.RGBA
	size := BallMedium
	if len(cells) == 0 {
		size = BallHidden
	}
	for i := 0; i < len(cells); i++ {
		bg = game.BallNoColor.Color()
		if size == BallHidden {
			fg = game.BallNoColor.Color()
		} else {
			fg = cells[i].Color().Color()
		}
		icon := NewBallIcon(size, bg, fg)
		icon.setup(size)
		x, y, w, h := t.nextBallsLayout.Rect.GetRect()
		icon.Resize([]int{x, y, w / len(cells), h})
		t.nextBallsLayout.Container[i].(*eui.Icon).SetIcon(icon.GetImage())
	}
}

func (t *Table) Update(dt int) {
	for _, v := range t.nextBallsLayout.Container {
		v.Update(dt)
	}
	for _, v := range t.Container {
		v.Update(dt)
	}
}

func (t *Table) Draw(surface *ebiten.Image) {
	if !t.IsVisible() {
		return
	}
	for _, v := range t.nextBallsLayout.Container {
		v.Draw(surface)
	}
	for _, v := range t.Container {
		v.Draw(surface)
	}
}

func (t *Table) Resize(rect []int) {
	t.View.Resize(rect)
	w0 := t.GetRect().W
	x, y := t.GetRect().Pos()
	w, h := int(float64(t.GetRect().W)*0.3), t.GetRect().H
	t.leftLbl.Resize([]int{x, y, w, h})
	t.rightLbl.Resize([]int{x + w0 - w, y, w, h})
	t.nextBallsLayout.Resize([]int{x + (w0-w)/2, y, w, h})
}
