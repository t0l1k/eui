package scene

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/lines/game"
)

type Board struct {
	eui.View
	table                  *Table
	gameLayout             *eui.GridLayoutRightDown
	field                  *game.Field
	showWay                bool
	showWayDt              int
	varScore, varScoreBest *eui.IntVar
	bestScore              int
}

func NewBoard(dim int) *Board {
	b := &Board{}
	b.SetupView()
	b.field = game.NewField(dim)
	b.table = NewTable()
	b.table.Visible(false)
	b.Add(b.table)
	b.gameLayout = eui.NewGridLayoutRightDown(b.field.Dim())
	b.Add(b.gameLayout)
	b.varScore = eui.NewIntVar(0)
	b.varScore.Attach(b.table.leftLbl)
	b.bestScore = 10
	b.varScoreBest = eui.NewIntVar(0)
	b.varScoreBest.Attach(b.table.rightLbl)
	b.varScoreBest.SetValue(b.bestScore)
	return b
}

func (b *Board) NewGame(dim int) {
	b.gameLayout.Container = nil
	b.field.NewGame(dim)
	for _, cell := range b.field.GetField() {
		cell.Reset()
		btn := NewCellIcon(cell, b.gameLogic)
		cell.State.Attach(btn)
		b.gameLayout.Add(btn)
	}
	b.gameLayout.SetDim(b.field.Dim())
	b.field.NextMoveBalls()
	b.field.ShowFilledNext()
	b.field.NextMoveBalls()
	b.table.Setup(b.field.Conf.Balls)
	b.table.Visible(true)
}

func (b *Board) gameLogic(btn *eui.ButtonIcon) {
	for i := range b.gameLayout.Container {
		if b.gameLayout.Container[i].(*CellIcon).btn == btn {
			cell := b.field.GetField()[i]
			cellData := cell.State.Value().(*game.CellData)
			state := cellData.State
			if btn.IsMouseDownLeft() {
				if state == game.CellFilled || state == game.CellEmpty || state == game.CellFilledNext {
					if col, way := b.field.MakeMove(b.field.Pos(i)); len(way) > 0 {
						b.showWay = true
						b.showWayDt = 250
						b.setWayCells(col, way)
						b.table.SetNextMoveBalls(b.field.GetFilledNext())
						log.Println("make move for:", cell, cell.Color().String(), cellData, state, way)
					}
				}
			} else if btn.IsMouseDownRight() {
				log.Printf("data:(%v)(empty:%v)(filled:%v)(filledNext:%v)(filledAfterMove;%v)(markedForMove:%v)(markedFoeDel:%v)", cellData, cell.IsEmpty(), cell.IsFilled(), cell.IsFilledNext(), cell.IsFilledAfterMove(), cell.IsMarkedForMove(), cell.IsMarkedForDelete())
			}
			fmt.Println(b.field)
		}
	}
}

func (b *Board) setWayCells(col game.BallColor, way []int) {
	for _, value := range way {
		bg := game.BallNoColor.Color()
		fg := col.Color()
		icon := NewBallIcon(BallSmall, bg, fg)
		cell := b.gameLayout.Container[value].(*CellIcon)
		icon.Resize(cell.Rect.GetArr())
		image := icon.GetImage()
		cell.btn.SetIcons([]*ebiten.Image{image, image})
	}
	log.Println("set way colors done", way)
}

func (b *Board) Update(dt int) {
	if b.field.GetScore() > b.varScoreBest.Value().(int) {
		b.varScoreBest.SetValue(b.field.GetScore())
	}
	b.varScore.SetValue(b.field.GetScore())
	if b.showWay {
		if b.showWayDt > 0 {
			b.showWayDt -= dt
		} else {
			b.drawCellIcons()
			log.Println("restore board after anime way")
			b.showWay = false
		}
	}
	for _, v := range b.Container {
		v.Update(dt)
	}
	for _, v := range b.gameLayout.Container {
		v.Update(dt)
	}
	if b.IsDirty() && !b.showWay {
		b.drawCellIcons()
	}
}

func (b *Board) drawCellIcons() {
	for i, cell := range b.field.GetField() {
		bg := game.BallNoColor.Color()
		fg := cell.Color().Color()
		var size BallStatusType
		switch {
		case cell.IsEmpty():
			size = BallHidden
		case cell.IsFilledNext():
			size = BallMedium
		case cell.IsFilled():
			size = BallNormal
		case cell.IsFilledAfterMove():
			size = BallNormal
		case cell.IsMarkedForMove():
			size = BallNormal
		}
		icon := NewBallIcon(size, bg, fg)
		cell := b.gameLayout.Container[i].(*CellIcon)
		icon.Resize(cell.Rect.GetArr())
		image := icon.GetImage()
		cell.btn.SetIcons([]*ebiten.Image{image, image})
	}
}

func (b *Board) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.Container {
		v.Draw(surface)
	}
	for _, v := range b.gameLayout.Container {
		v.Draw(surface)
	}
}

func (b *Board) Resize(rect []int) {
	b.View.Resize(rect)
	b.gameLayout.SetCellMargin(int(float64(b.GetRect().GetLowestSize()) * 0.008))
	w0, h0 := b.GetRect().Size()
	x0, y0 := b.GetRect().Pos()
	dim := b.field.Conf.Dim
	cellSize := getCellSize(b.GetRect(), dim)
	b.gameLayout.Resize([]int{x0, y0 + cellSize, w0, h0 - cellSize})
	b.table.Resize([]int{x0 + (w0-cellSize*dim)/2, y0, cellSize * dim, cellSize})
	b.Dirty(true)
	log.Println("board resize done")
}

func getCellSize(rect *eui.Rect, dim int) (size int) {
	r := dim
	c := dim
	for r*size < rect.W && c*size < rect.H {
		size += 1
	}
	return size
}
