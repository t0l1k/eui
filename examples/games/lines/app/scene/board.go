package scene

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/lines/game"
)

type Board struct {
	eui.DrawableBase
	table                  *Table
	gameLayout             *eui.GridLayoutRightDown
	field                  *game.Field
	showWay                bool
	showWayDt              int
	varScore, varScoreBest *eui.Signal[int]
	bestScore              int
}

func NewBoard(dim int) *Board {
	b := &Board{}
	b.field = game.NewField(dim)
	b.table = NewTable()
	b.table.Visible(false)
	b.Add(b.table)
	r, c := b.field.Dim()
	b.gameLayout = eui.NewGridLayoutRightDown(float64(r), float64(c))
	b.varScore = eui.NewSignal(func(a, b int) bool { return a == b })
	b.varScore.Connect(func(data int) {
		b.table.leftLbl.SetText(strconv.Itoa(data))
	})
	b.bestScore = 10
	b.varScoreBest = eui.NewSignal(func(a, b int) bool { return a == b })
	b.varScoreBest.ConnectAndFire(func(data int) {
		b.table.rightLbl.SetText(strconv.Itoa(data))
	}, b.bestScore)
	return b
}

func (b *Board) NewGame(dim int) {
	b.gameLayout.ResetContainerBase()
	b.field.NewGame(dim)
	for _, cell := range b.field.GetField() {
		cell.Reset()
		btn := NewCellIcon(cell, b.gameLogic)
		cell.State.Connect(func(v *game.CellData) {
			c := btn
			switch v.State {
			case game.CellEmpty:
				c.anim = BallAnimNo
				c.animStatus = BallHidden
				c.fg = v.Color.Color()
				c.updateIcon(BallHidden)
			case game.CellFilledNext:
				c.anim = BallAnimFilledNext
				c.animStatus = BallHidden
				c.fg = v.Color.Color()
				c.updateIcon(BallSmall)
			case game.CellFilled:
				c.anim = BallAnimFilled
				c.animStatus = BallMedium
				c.fg = v.Color.Color()
				c.updateIcon(BallNormal)
			case game.CellMarkedForMove:
				c.anim = BallAnimJump
				c.animStatus = BallJumpCenter
				c.fg = v.Color.Color()
				c.updateIcon(BallNormal)
			case game.CellFilledAfterMove:
				c.anim = BallAnimFilledAfterMove
				c.animStatus = BallJumpCenter
				c.fg = v.Color.Color()
				c.updateIcon(BallNormal)
			case game.CellMarkedForDelete:
				c.anim = BallAnimDelete
				c.animStatus = BallBig
				c.fg = v.Color.Color()
				c.updateIcon(BallBig)
			}

		})
		b.gameLayout.Add(btn)
	}
	r, c := b.field.Dim()
	b.gameLayout.SetDim(float64(r), float64(c))
	b.field.NextMoveBalls()
	b.field.ShowFilledNext()
	b.field.NextMoveBalls()
	b.table.Setup(b.field.Conf.Balls)
	b.table.Visible(true)
}

func (b *Board) gameLogic(btn *eui.Button) {
	for i := range b.gameLayout.GetContainer() {
		if b.gameLayout.GetContainer()[i].(*CellIcon).btn == btn {
			cell := b.field.GetField()[i]
			cellData := cell.State.Value()
			state := cellData.State
			if btn.IsMouseDownLeft() && b.field.InGame {
				if state == game.CellFilled || state == game.CellEmpty || state == game.CellFilledNext {
					if way := b.field.MakeMove(b.field.Pos(i)); len(way) > 0 {
						b.showWay = true
						b.showWayDt = 250
						b.setWayCells(cell.Color(), way)
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
		cell := b.gameLayout.GetContainer()[value].(*CellIcon)
		cell.icon.setup(BallSmall, bg, fg)
	}
	log.Println("set way colors done", way, col)
}

func (b *Board) Update(dt int) {
	if b.field.GetScore() > b.varScoreBest.Value() {
		b.varScoreBest.Emit(b.field.GetScore())
	}
	b.varScore.Emit(b.field.GetScore())
	if b.showWay {
		if b.showWayDt > 0 {
			b.showWayDt -= dt
		} else {
			b.drawCellIcons()
			log.Println("restore board after anime way")
			b.showWay = false
		}
	}
	for _, v := range b.GetContainer() {
		v.Update(dt)
	}
	for _, v := range b.gameLayout.GetContainer() {
		v.Update(dt)
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
		cl := b.gameLayout.GetContainer()[i].(*CellIcon)
		cl.icon.setup(size, bg, fg)
	}
}

func (b *Board) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range b.gameLayout.GetContainer() {
		v.Draw(surface)
	}
}

func (b *Board) Resize(rect []int) {
	b.Rect(eui.NewRect(rect))
	b.SpriteBase.Resize(rect)
	b.gameLayout.SetCellMargin(float64(b.GetRect().GetLowestSize()) * 0.008)
	w0, h0 := b.GetRect().Size()
	x0, y0 := b.GetRect().Pos()
	dim := b.field.Conf.Dim
	cellSize := getCellSize(b.GetRect(), dim)
	b.gameLayout.Resize([]int{x0, y0 + cellSize, w0, h0 - cellSize})
	b.table.Resize([]int{x0 + (w0-cellSize*dim)/2, y0, cellSize * dim, cellSize})
	b.Dirty = true
	log.Println("board resize done")
}

func getCellSize(rect eui.Rect, dim int) (size int) {
	r := dim
	c := dim
	for r*size < rect.W && c*size < rect.H {
		size += 1
	}
	return size
}
