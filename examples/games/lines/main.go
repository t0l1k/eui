package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/lines/game"
)

const (
	GameTitle = "Очищай поле вовремя aka Color Lines"
	BNew      = "Новая Игра"

	FieldSizeSmall  = 7
	FieldSizeNormal = 9
	FieldSizeBig    = 11
)

type BallAnimType int

const (
	BallAnimNo BallAnimType = iota
	BallAnimFilledNext
	BallAnimFilled
	BallAnimJump
	BallAnimFilledAfterMove
	BallAnimDelete
)

const (
	BallHidden BallStatusType = iota
	BallSmall
	BallMedium
	BallNormal
	BallBig

	BallJumpDown
	BallJumpCenter
	BallJumpUp
)

type BallStatusType int

func (j BallStatusType) IsSmall() bool  { return j == BallSmall }
func (j BallStatusType) IsHidden() bool { return j == BallHidden }
func (j BallStatusType) IsMedium() bool { return j == BallMedium }
func (j BallStatusType) IsNormal() bool { return j == BallNormal }
func (j BallStatusType) IsBig() bool    { return j == BallBig }

func (j *BallStatusType) FilledNext() {
	switch *j {
	case BallHidden:
		*j = BallSmall
	case BallSmall:
		*j = BallMedium
	}
}

func (j *BallStatusType) Jump() {
	switch *j {
	case BallJumpDown:
		*j = BallJumpCenter
	case BallJumpCenter:
		*j = BallJumpUp
	case BallJumpUp:
		*j = BallJumpDown
	}
}

func (j *BallStatusType) Delete() {
	switch *j {
	case BallBig:
		*j = BallNormal
	case BallNormal:
		*j = BallMedium
	case BallMedium:
		*j = BallSmall
	case BallSmall:
		*j = BallHidden
	}
}

func (j BallStatusType) String() (res string) {
	return []string{
		"ball hidden",
		"ball small",
		"ball mediun",
		"ball normal",
		"ball big",
		"jump down",
		"jump center",
		"jump up",
		strconv.Itoa(int(j)) + "!",
	}[j]
}

type BallIcon struct {
	*eui.Drawable
	size   float32
	status BallStatusType
	bg, fg color.Color
}

func NewBallIcon(status BallStatusType, bg, fg color.Color) *BallIcon {
	i := &BallIcon{
		Drawable: eui.NewDrawable(),
		bg:       bg,
		fg:       fg,
		status:   status,
	}
	i.setup(status, bg, fg)
	return i
}

func (i *BallIcon) setup(status BallStatusType, bg, fg color.Color) {
	i.status = status
	switch status {
	case BallHidden:
		i.size = 0
	case BallSmall:
		i.size = 0.146
	case BallMedium:
		i.size = 0.236
	case BallNormal, BallJumpCenter, BallJumpUp, BallJumpDown:
		i.size = 0.382
	case BallBig:
		i.size = 0.5
	}
	i.bg = bg
	i.fg = fg
	i.MarkDirty()
}

func (i *BallIcon) Layout() {
	i.Drawable.Layout()
	r, g, b, _ := i.bg.RGBA()
	a := 255
	bg := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	i.Image().Fill(bg)
	if i.size > 0 {
		rad := float32(i.Rect().GetLowestSize()) * i.size
		x, y := float32(i.Rect().W/2), float32(i.Rect().H/2)
		margin := (float32(i.Rect().W) - rad*2) / 3
		switch i.status {
		case BallJumpUp:
			y = float32(i.Rect().H/2) - margin
		case BallJumpCenter:
			y = float32(i.Rect().H / 2)
		case BallJumpDown:
			y = float32(i.Rect().H/2) + margin
		}
		vector.DrawFilledCircle(i.Image(), x, y, rad, i.fg, true)
	}
	i.ClearDirty()
}

func (i *BallIcon) GetImage() *ebiten.Image {
	if i.IsDirty() {
		i.Layout()
	}
	return i.Image()
}

type CellIcon struct {
	*eui.Drawable
	btn        *eui.Button
	cell       *game.Cell
	anim       BallAnimType
	animStatus BallStatusType
	jumpDt     int
	icon       *BallIcon
	fg         color.Color
}

func NewCellIcon(cell *game.Cell, fn func(b *eui.Button)) *CellIcon {
	c := &CellIcon{Drawable: eui.NewDrawable()}
	c.cell = cell
	c.icon = NewBallIcon(BallHidden, game.BallNoColor.Color(), game.BallNoColor.Color())
	c.icon.SetRect(eui.NewRect([]int{0, 0, 1, 1}))
	c.btn = eui.NewButton("", fn)
	c.anim = BallAnimNo
	c.animStatus = BallHidden
	return c
}

func (b *CellIcon) Hit(pt eui.Point[int]) eui.Drawabler {
	if !pt.In(b.Rect()) || b.IsHidden() || b.IsDisabled() {
		return nil
	}
	log.Println("CellIcon:Hit:", b.Rect(), b.cell.String())
	return b
}
func (b *CellIcon) WantBlur() bool { return true }

func (b *CellIcon) MouseUp(md eui.MouseData) {
	b.btn.MouseUp(md)
	log.Println("CellIcon:MouseUp", b.Rect(), b.cell.String(), b.btn)
}

func (c *CellIcon) Tick(td eui.TickData) {
	dt := int(td.Duration())
	if c.anim == BallAnimFilledNext && !c.animStatus.IsMedium() {
		c.jumpDt += dt
		if c.jumpDt > 250 {
			c.animStatus.FilledNext()
			c.jumpDt = 0
			c.updateIcon(c.animStatus)
			// log.Println("cell icon anim filled next", c.animStatus.String())
		}
	}
	if c.anim == BallAnimJump {
		c.jumpDt += dt
		if c.jumpDt > 100 {
			c.animStatus.Jump()
			c.jumpDt = 0
			c.updateIcon(c.animStatus)
			// log.Println("cell icon anim jump", c.animStatus.String())
		}
	}

	if c.anim == BallAnimFilledAfterMove && !c.animStatus.IsHidden() {
		c.jumpDt += dt
		if c.jumpDt > 250 {
			c.jumpDt = 0
			c.anim = BallAnimFilled
			c.updateIcon(c.animStatus)
			c.cell.SetFilled()
			log.Println("cell icon anim show move way end", c.animStatus.String(), c.cell.IsFilledAfterMove())
		}
	}
	if c.anim == BallAnimDelete {
		if !c.animStatus.IsHidden() {
			c.jumpDt += dt
			if c.jumpDt > 50 {
				c.animStatus.Delete()
				c.jumpDt = 0
				c.updateIcon(c.animStatus)
				log.Println("cell icon anim delete", c.animStatus.String())
			}
		} else {
			c.animStatus = BallHidden
			c.updateIcon(c.animStatus)
			c.anim = BallAnimNo
			c.cell.Reset()
			log.Println("cell icon anim delete done", c.animStatus.String())
		}
	}
}

func (c *CellIcon) updateIcon(ballStatus BallStatusType) {
	bg := game.BallNoColor.Color()
	c.icon.setup(ballStatus, bg, c.fg)
	c.icon.SetRect(c.Rect())
}

func (c *CellIcon) Draw(surface *ebiten.Image) {
	if c.IsHidden() {
		return
	}
	if c.IsDirty() {
		c.icon.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := c.Rect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(c.icon.Image(), op)
}

func (c *CellIcon) Layout() {}

func (c *CellIcon) SetRect(rect eui.Rect[int]) {
	c.Drawable.SetRect(rect)
	c.btn.SetRect(rect)
	c.icon.SetRect(rect)
	c.ImageReset()
}

type Table struct {
	*eui.Container
	leftLbl, rightLbl *eui.Label
	nextBallsLayout   *eui.Container
}

func NewTable() *Table {
	t := &Table{Container: eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{15, 20, 30, 20, 15}, 3))}
	t.leftLbl = eui.NewLabel("0")
	t.Add(t.leftLbl)
	t.Add(eui.NewDrawable())
	t.nextBallsLayout = eui.NewContainer(eui.NewHBoxLayout(1))
	t.Add(t.nextBallsLayout)
	t.Add(eui.NewDrawable())
	t.rightLbl = eui.NewLabel("100")
	t.Add(t.rightLbl)
	return t
}

func (t *Table) Setup(balls int) {
	t.nextBallsLayout.ResetContainer()
	for i := 0; i < balls; i++ {
		icon := NewBallIcon(BallHidden, game.BallNoColor.Color(), game.BallNoColor.Color())
		x, y, w, h := t.nextBallsLayout.Rect().GetRect()
		icon.SetRect(eui.NewRect([]int{x, y, w / 3, h}))
		t.nextBallsLayout.Add(eui.NewIcon(icon.GetImage()))
	}
	t.SetRect(t.Rect())
	t.Layout()
}

func (t *Table) SetNextMoveBalls(cells []*game.Cell) {
	var bg, fg color.Color
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
		icon.SetRect(eui.NewRect([]int{x, y, w / len(cells), h}))
		t.nextBallsLayout.Children()[i].(*eui.Icon).SetIcon(icon.GetImage())
	}
}

type Board struct {
	*eui.Container
	table                  *Table
	gameLayout             *eui.Container
	field                  *game.Field
	showWay                bool
	showWayDt              int
	varScore, varScoreBest *eui.Signal[int]
	bestScore              int
}

func NewBoard(dim int) *Board {
	b := &Board{Container: eui.NewContainer(eui.NewLayoutVerticalPercent([]int{10, 90}, 3))}
	b.field = game.NewField(dim)
	b.table = NewTable()
	b.Add(b.table)
	r, c := b.field.Dim()
	b.gameLayout = eui.NewContainer(eui.NewGridLayout(r, c, 1))
	b.Add(b.gameLayout)
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
	b.gameLayout.ResetContainer()
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
	b.gameLayout.SetLayout(eui.NewSquareGridLayout(r, c, 1))
	b.gameLayout.Layout()
	b.field.NextMoveBalls()
	b.field.ShowFilledNext()
	b.field.NextMoveBalls()
	b.table.Setup(b.field.Conf.Balls)
}

func (b *Board) gameLogic(btn *eui.Button) {
	for i := range b.gameLayout.Children() {
		if b.gameLayout.Children()[i].(*CellIcon).btn == btn {
			cell := b.field.GetField()[i]
			cellData := cell.State.Value()
			state := cellData.State
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && b.field.InGame {
				if state == game.CellFilled || state == game.CellEmpty || state == game.CellFilledNext {
					if way := b.field.MakeMove(b.field.Pos(i)); len(way) > 0 {
						b.showWay = true
						b.showWayDt = 250
						b.setWayCells(cell.Color(), way)
						b.table.SetNextMoveBalls(b.field.GetFilledNext())
						log.Println("make move for:", cell, cell.Color().String(), cellData, state, way)
					}
				}
			}
			fmt.Println(b.field)
		}
	}
}

func (b *Board) setWayCells(col game.BallColor, way []int) {
	for _, value := range way {
		bg := game.BallNoColor.Color()
		fg := col.Color()
		cell := b.gameLayout.Children()[value].(*CellIcon)
		cell.icon.setup(BallSmall, bg, fg)
	}
	log.Println("set way colors done", way, col)
}

func (b *Board) Tick(td eui.TickData) {
	if b.field.GetScore() > b.varScoreBest.Value() {
		b.varScoreBest.Emit(b.field.GetScore())
	}
	b.varScore.Emit(b.field.GetScore())
	if b.showWay {
		if b.showWayDt > 0 {
			b.showWayDt -= int(td.Duration())
		} else {
			b.drawCellIcons()
			log.Println("restore board after anime way")
			b.showWay = false
		}
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
		cl := b.gameLayout.Children()[i].(*CellIcon)
		cl.icon.setup(size, bg, fg)
	}
}

type Dialog struct {
	*eui.Container
	title            *eui.Label
	comboSelGameDiff *eui.SpinBox[int]
	diff             int
}

func NewDialog(title string, fn func(btn *eui.Button)) *Dialog {
	d := &Dialog{Container: eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 90, 5}, 1))}

	contTitle := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{10, 90}, 1))
	contTitle.Add(eui.NewButton("X", func(b *eui.Button) {
		eui.GetUi().Pop()
	}))
	d.title = eui.NewLabel(title)
	contTitle.Add(d.title)

	data := []int{FieldSizeSmall, FieldSizeNormal, FieldSizeBig}
	d.diff = data[1]
	d.comboSelGameDiff = eui.NewSpinBox(
		"Выбор размер поля",
		data,
		1,
		func(i int) string { return fmt.Sprintf("%v", i) },
		func(a, b int) int { return a - b },
		func(a, b int) bool { return a == b },
		false,
		2,
	)
	d.comboSelGameDiff.SelectedValue.Connect(func(value int) {
		d.diff = value
	})

	d.Add(contTitle)
	d.Add(d.comboSelGameDiff)
	d.Add(eui.NewButton(BNew, fn))

	return d
}
func (d *Dialog) SetTitle(title string) { d.title.SetText(title) }

func main() {
	eui.Init(eui.GetUi().SetTitle(GameTitle).SetSize(800, 600))
	eui.Run(func() *eui.Scene {
		var (
			board  *Board
			dialog *Dialog
		)
		s := eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 95}, 5))
		s.Add(eui.NewTopBar(GameTitle, func(b *eui.Button) {
			dialog.SetTitle("Выбор новой игры")
			dialog.Show()
			board.Hide()
		}).SetUseStopwatch())
		contBoard := eui.NewContainer(eui.NewStackLayout(25))
		dialog = NewDialog("Запустить игру", func(dlg *eui.Button) {
			if dlg.Text() == BNew {
				board.NewGame(dialog.diff)
			}
			dialog.Hide()
			board.Show()
		})
		contBoard.Add(dialog)
		board = NewBoard(dialog.diff)
		contBoard.Add(board)
		s.Add(contBoard)
		return s
	}())
	eui.Quit(func() {})
}
