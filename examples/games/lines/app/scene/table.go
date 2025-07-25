package scene

import (
	"image/color"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/lines/game"
)

type Table struct {
	*eui.Container
	leftLbl, rightLbl *eui.Text
	nextBallsLayout   *eui.Container
}

func NewTable() *Table {
	t := &Table{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	t.leftLbl = eui.NewText("0")
	t.Add(t.leftLbl)
	t.rightLbl = eui.NewText("100")
	t.Add(t.rightLbl)
	t.nextBallsLayout = eui.NewContainer(eui.NewHBoxLayout(1))
	return t
}

func (t *Table) Setup(balls int) {
	t.nextBallsLayout.ResetContainer()
	for i := 0; i < balls; i++ {
		icon := NewBallIcon(BallHidden, game.BallNoColor.Color(), game.BallNoColor.Color())
		x, y, w, h := t.nextBallsLayout.Rect().GetRect()
		icon.Resize(eui.NewRect([]int{x, y, w / 3, h}))
		t.nextBallsLayout.Add(eui.NewIcon(icon.GetImage()))
	}
	t.Resize(t.Rect())
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
		defer icon.Close()
		icon.setup(size, bg, fg)
		x, y, w, h := t.nextBallsLayout.Rect().GetRect()
		icon.Resize(eui.NewRect([]int{x, y, w / len(cells), h}))
		t.nextBallsLayout.Childrens()[i].(*eui.Icon).SetIcon(icon.GetImage())
	}
}

func (t *Table) Resize(rect eui.Rect[int]) {
	t.SetRect(rect)
	w0 := t.Rect().W
	x, y := t.Rect().Pos()
	w, h := int(float64(t.Rect().W)*0.3), t.Rect().H
	t.leftLbl.Resize(eui.NewRect([]int{x, y, w, h}))
	t.rightLbl.Resize(eui.NewRect([]int{x + w0 - w, y, w, h}))
	t.nextBallsLayout.Resize(eui.NewRect([]int{x + (w0-w)/2, y, w, h}))
}
