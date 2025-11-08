package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/minesweeper/res"
	"github.com/t0l1k/eui/utils"
	"golang.org/x/image/colornames"
)

const (
	Title      = "Игра Сапер"
	Easy       = "Легко"
	Medium     = "Средне"
	Hard       = "Сложно"
	Custom     = "Настрой сложность"
	QuitDialog = "<"
	PlayAgain  = "Игра новая"
	PlayRepeat = "Игру повторить"
)

type CellState rune

const (
	CellClosed       CellState = 'C'
	CellOpened       CellState = 'O'
	CellFlagged      CellState = 'F'
	CellQuestioned   CellState = 'Q'
	CellFirstMined   CellState = 'f'
	CellSaved        CellState = 'v'
	CellBlown        CellState = 'b'
	CellWrongFlagged CellState = 'w'
)

type Cell struct {
	State *eui.Signal[CellState]
	mined bool
	count int
}

func NewCell() *Cell   { return &Cell{State: eui.NewSignal(func(a, b CellState) bool { return a == b })} }
func (c *Cell) Reset() { c.State.Emit(CellClosed) }
func (c *Cell) Open()  { c.State.Emit(CellOpened) }
func (c *Cell) Mark() {
	switch c.State.Value() {
	case CellClosed:
		c.State.Emit(CellFlagged)
	case CellFlagged:
		c.State.Emit(CellQuestioned)
	case CellQuestioned:
		c.State.Emit(CellClosed)
	}
}
func (c *Cell) IsClosed() bool     { return c.State.Value() == CellClosed }
func (c *Cell) IsOpen() bool       { return c.State.Value() == CellOpened }
func (c *Cell) IsFlagged() bool    { return c.State.Value() == CellFlagged }
func (c *Cell) IsQuestioned() bool { return c.State.Value() == CellQuestioned }
func (c *Cell) IsMined() bool      { return c.mined }
func (c *Cell) Count() int         { return c.count }
func (c *Cell) String() string {
	s := ""
	switch c.State.Value() {
	case CellClosed:
		s += " "
	case CellFlagged:
		s += "F"
	case CellQuestioned:
		s += "Q"
	case CellFirstMined:
		s += "f"
	case CellSaved:
		s += "v"
	case CellBlown:
		s += "b"
	case CellWrongFlagged:
		s += "w"
	case CellOpened:
		if c.mined {
			s += "*"
		} else {
			switch c.count {
			case 0:
				s += " "
			default:
				s += strconv.Itoa(c.count)
			}
		}
	default:
		s += "!"
	}
	return s
}

type Field []*Cell

func NewField() Field              { return make([]*Cell, 0) }
func (f Field) Cell(idx int) *Cell { return f[idx] }

type Board struct {
	State                     *eui.Signal[GameState]
	Field                     Field
	Conf                      *GameData
	CountMarked, CountFlagged *eui.Signal[int]
	Sw                        *eui.Stopwatch
	first                     eui.Point[int]
}

func NewGame() *Board {
	b := &Board{
		State:        eui.NewSignal(func(a, b GameState) bool { return a == b }),
		Field:        NewField(),
		Sw:           eui.NewStopwatch(),
		CountMarked:  eui.NewSignal(func(a, b int) bool { return a == b }),
		CountFlagged: eui.NewSignal(func(a, b int) bool { return a == b }),
	}
	return b
}

func (b *Board) New(r, c, m int) {
	b.Conf = NewGameData(r, c, m)
	b.State.Emit(GameStart)
	b.Conf.duration.Emit(0)
	b.Sw.Reset()
	b.CountFlagged.Emit(0)
	b.CountMarked.Emit(0)
	b.Field = NewField()
	for range r * c {
		cell := NewCell()
		cell.Reset()
		b.Field = append(b.Field, cell)
	}
}

func (b *Board) Reset() {
	b.State.Emit(GamePlay)
	b.Conf.duration.Emit(0)
	b.Sw.Reset()
	b.CountFlagged.Emit(0)
	b.CountMarked.Emit(b.Conf.dim.mines)
	for _, v := range b.Field {
		v.Reset()
	}
	b.Open(b.first.X, b.first.Y)
	b.Sw.Start()
}

func (b *Board) Shuffle(x0, y0 int) {
	b.first = eui.NewPoint(x0, y0)
	for mines := 0; mines < b.Conf.dim.mines; {
		x, y := rand.Intn(b.Conf.dim.row), rand.Intn(b.Conf.dim.column)
		if x != x0 && y != y0 {
			cell := b.Field[b.Idx(x, y)]
			if !cell.mined {
				cell.mined = true
				mines++
				b.CountMarked.Emit(b.CountMarked.Value() + 1)
			}
		}
	}
	for idx, cell := range b.Field {
		if !cell.mined {
			x, y := b.Pos(idx)
			b.GetNeighbours(x, y, func(i1, i2 int) {
				newCell := b.Field[b.Idx(i1, i2)]
				if newCell.mined {
					b.Field[idx].count++
				}
			})
		}
	}
	b.State.Emit(GamePlay)
	b.Sw.Start()
	log.Println("Board:Reset:", b.String2())
}

func (b *Board) Open(x, y int) {
	if b.isFieldEdge(x, y) {
		return
	}
	cell := b.Field.Cell(b.Idx(x, y))
	if cell.IsFlagged() || cell.IsOpen() {
		return
	}
	cell.Open()
	b.CountMarked.Emit(b.CountMarked.Value() + 1)
	if cell.mined {
		cell.State.Emit(CellFirstMined)
		b.State.Emit(GameGameOver)
		b.Sw.Stop()
		return
	}
	if b.Winned() {
		b.State.Emit(GameWin)
		b.Sw.Stop()
		return
	}
	if cell.count > 0 {
		return
	}
	b.GetNeighbours(x, y, func(i1, i2 int) {
		b.Open(i1, i2)
	})
}

func (b *Board) AutoMark(x, y int) {
	if b.Cell(x, y).IsClosed() {
		return
	}
	var countFlaggs, countClosed int
	b.GetNeighbours(x, y, func(i1, i2 int) {
		cell := b.Cell(i1, i2)
		if cell.IsFlagged() {
			countFlaggs++
		} else if cell.IsClosed() || cell.IsQuestioned() {
			countClosed++
		}
	})
	if countFlaggs+countClosed == b.Cell(x, y).Count() {
		b.GetNeighbours(x, y, func(i1, i2 int) {
			cell := b.Cell(i1, i2)
			if cell.IsClosed() || cell.IsQuestioned() {
				cell.State.Emit(CellFlagged)
				b.CountFlagged.Emit(b.CountFlagged.Value() + 1)
			}
		})
	} else if countFlaggs == b.Cell(x, y).Count() {
		b.GetNeighbours(x, y, func(i1, i2 int) {
			cell := b.Cell(i1, i2)
			if cell.IsClosed() || cell.IsQuestioned() {
				b.Open(i1, i2)
			}
		})
	}
}

func (b *Board) Mark(x, y int) {
	b.Cell(x, y).Mark()
	if b.Cell(x, y).IsFlagged() {
		b.CountFlagged.Emit(b.CountFlagged.Value() + 1)
	} else if b.Cell(x, y).IsOpen() || b.Cell(x, y).IsQuestioned() {
		b.CountFlagged.Emit(b.CountFlagged.Value() - 1)
	}
}

func (b *Board) Pos(idx int) (int, int) { return idx % b.Conf.dim.row, idx / b.Conf.dim.row }
func (b *Board) Idx(x, y int) int       { return y*b.Conf.dim.row + x }
func (b *Board) Cell(x, y int) *Cell    { return b.Field[b.Idx(x, y)] }
func (b *Board) isFieldEdge(x, y int) bool {
	return x < 0 || x > b.Conf.dim.row-1 || y < 0 || y > b.Conf.dim.column-1
}

func (b *Board) GetNeighbours(x, y int, fn func(int, int)) {
	for dy := -1; dy < 2; dy++ {
		for dx := -1; dx < 2; dx++ {
			nx := x + dx
			ny := y + dy
			if !b.isFieldEdge(nx, ny) {
				fn(nx, ny)
			}
		}
	}
}

func (b *Board) Winned() bool           { return b.CountMarked.Value() == b.Conf.dim.row*b.Conf.dim.column }
func (b *Board) CountFlags() (int, int) { return b.CountFlagged.Value(), b.Conf.dim.mines }

func (g *Board) MarkGameAfterEnd() {
	for _, cell := range g.Field {
		switch cell.State.Value() {
		case CellFlagged:
			if !cell.IsMined() {
				cell.State.Emit(CellWrongFlagged)
			} else {
				cell.State.Emit(CellSaved)
			}
		case CellClosed, CellQuestioned:
			if g.State.Value() == GameWin {
				cell.State.Emit(CellSaved)
			} else if cell.IsMined() {
				cell.State.Emit(CellBlown)
			}
		}
	}
	g.Conf.Set(g.Sw.Duration(), g.State.Value() == GameWin)
}

func (b *Board) String() string {
	str := ""
	str += strconv.Itoa(b.Conf.dim.row) + ","
	str += strconv.Itoa(b.Conf.dim.column) + ",("
	str += strconv.Itoa(b.CountFlagged.Value()) + "/"
	str += strconv.Itoa(b.Conf.dim.mines) + ":"
	str += strconv.Itoa(b.CountMarked.Value()) + ")"
	str += b.State.Value().String() + "\n"
	for y := 0; y < b.Conf.dim.column; y++ {
		for x := 0; x < b.Conf.dim.row; x++ {
			str += b.Field[b.Idx(x, y)].String()
		}
		str += "\n"
	}
	str += eui.FormatSmartDuration(b.Conf.duration.Value(), false)
	return str
}

func (b *Board) String2() string {
	str := strconv.Itoa(b.Conf.dim.row) + ","
	str += strconv.Itoa(b.Conf.dim.column) + ",("
	str += strconv.Itoa(0) + "/"
	str += strconv.Itoa(b.Conf.dim.mines) + ")"
	str += b.State.Value().String() + "\n"
	for y := 0; y < b.Conf.dim.column; y++ {
		for x := 0; x < b.Conf.dim.row; x++ {
			cell := b.Field.Cell(b.Idx(x, y))
			if cell.mined {
				str += "*"
			} else {
				str += strconv.Itoa(b.Field[b.Idx(x, y)].count)
			}
		}
		str += "\n"
	}
	str += eui.FormatSmartDuration(b.Conf.duration.Value(), false)
	return str
}

type GameState int

const (
	GameStart GameState = iota
	GamePlay
	GamePause
	GameGameOver
	GameWin
)

func (g GameState) String() string {
	return []string{"Start", "Play", "Pause", "Game Over", "Winned"}[g]
}

type Dim struct{ row, column, mines int }

func NewDim(r, c, m int) Dim { return Dim{row: r, column: c, mines: m} }
func (d Dim) Percent() int   { return int(utils.PercentOf(float64(d.mines), float64(d.row*d.column))) }
func (d *Dim) SetMines(value float64) {
	d.mines = int(utils.ValueFromPercent(value, float64(d.row*d.column)))

}
func (d Dim) Empty() bool    { return d.row == 0 || d.column == 0 || d.mines == 0 }
func (d Dim) String() string { return fmt.Sprintf("Dim:%v %v %v", d.row, d.column, d.mines) }

type GameData struct {
	dim      Dim
	start    time.Time
	duration *eui.Signal[time.Duration]
	wined    bool
}

func NewGameData(r, c, m int) *GameData {
	return &GameData{
		dim:      NewDim(r, c, m),
		start:    time.Now(),
		duration: eui.NewSignal(func(a, b time.Duration) bool { return a == b })}
}
func (g *GameData) Set(d time.Duration, win bool) { g.duration.Emit(d); g.wined = win }
func (g *GameData) String() string {
	str := ""
	if g.wined {
		str = "Game Won!"
	} else {
		str = "Game Lost!"
	}
	str += fmt.Sprint(" Time:", eui.FormatSmartDuration(g.duration.Value(), true))
	return str
}

type GamesData map[Dim][]GameData

type Topbar struct {
	*eui.Container
	quitBtn, gameBtn  *eui.Button
	lblLeft, lblRight *eui.Label
}

func NewTopbar(fn func(*eui.Button)) *Topbar {
	t := &Topbar{Container: eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{5, 25, 15, 10, 25, 20}, 1))}
	t.quitBtn = eui.NewButton(QuitDialog, fn)
	t.lblLeft = eui.NewLabel("0/0")
	t.gameBtn = eui.NewButtonIcon([]*ebiten.Image{res.SmileSprites[0], res.SmileSprites[1]}, fn)
	t.gameBtn.SetText(PlayAgain)
	t.lblRight = eui.NewLabel("000")
	t.Add(t.quitBtn)
	t.Add(t.lblLeft)
	t.Add(eui.NewDrawable())
	t.Add(t.gameBtn)
	t.Add(eui.NewDrawable())
	t.Add(t.lblRight)
	return t
}

type GameView struct {
	*eui.Container
	topbar *Topbar
	board  *eui.Container
	game   *Board
	fn     func(*eui.Button)
}

func NewGameView(fn func(*eui.Button)) *GameView {
	c := &GameView{Container: eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 95}, 1))}
	c.fn = fn
	c.topbar = NewTopbar(c.fn)
	c.Add(c.topbar)
	c.board = eui.NewContainer(eui.NewSquareGridLayout(2, 2, 5))
	c.Add(c.board)
	c.game = NewGame()
	return c
}

func (g *GameView) New(dim Dim) {
	g.game.New(dim.row, dim.column, dim.mines)
	g.board.ResetContainer()
	g.board.SetLayout(eui.NewSquareGridLayout(dim.row, dim.column, 5))
	for idx, cell := range g.game.Field {
		x, y := g.game.Pos(idx)
		g.board.Add(eui.NewButton(cell.String(), func(b *eui.Button) {
			switch g.game.State.Value() {
			case GameStart:
				g.game.Shuffle(x, y)
				g.game.Open(x, y)
			case GamePlay:
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
					if !cell.IsOpen() {
						g.game.Open(x, y)
					} else {
						g.game.AutoMark(x, y)
					}
				}
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
					g.game.Mark(x, y)
				}
			case GameGameOver:
			}

		}))
		cell.State.Connect(func(data CellState) {
			btn := g.board.Children()[idx].(*eui.Button)
			bgUp := colornames.Silver
			bgDown := colornames.Lightgray
			switch data {
			case CellClosed:
				btn.SetText("*")
				btn.SetBg(bgUp)
			case CellOpened:
				switch g.game.Cell(x, y).Count() {
				case 0:
					btn.SetText(" ")
					btn.SetBg(bgDown)
				case 1:
					btn.SetText("1")
					btn.SetBg(bgDown)
					btn.SetFg(colornames.Blue)
				case 2:
					btn.SetText("2")
					btn.SetBg(bgDown)
					btn.SetFg(colornames.Green)
				case 3:
					btn.SetText("3")
					btn.SetBg(bgDown)
					btn.SetFg(colornames.Red)
				case 4:
					btn.SetText("4")
					btn.SetBg(bgDown)
					btn.SetFg(colornames.Navy)
				case 5:
					btn.SetText("5")
					btn.SetBg(bgDown)
					btn.SetFg(colornames.Brown)
				case 6:
					btn.SetText("6")
					btn.SetBg(bgDown)
					btn.SetFg(colornames.Teal)
				case 7:
					btn.SetText("7")
					btn.SetBg(bgDown)
					btn.SetFg(colornames.Black)
				case 8:
					btn.SetText("8")
					btn.SetBg(bgDown)
					btn.SetFg(colornames.Darkgray)
				}
			case CellFlagged:
				btn.SetText("F")
			case CellQuestioned:
				btn.SetText("Q")
			case CellFirstMined:
				btn.SetText("f")
				btn.SetFg(colornames.Red)
			case CellBlown:
				btn.SetText("*")
				btn.SetFg(colornames.Red)
			case CellWrongFlagged:
				btn.SetText("w")
				btn.SetFg(colornames.Red)
			case CellSaved:
				btn.SetText("v")
				btn.SetFg(colornames.Green)
			}
		})
	}

	g.game.State.Connect(func(data GameState) {
		switch data {
		case GameStart:
			g.game.Sw.Reset()
			g.game.CountFlagged.Emit(0)
			g.game.Conf.duration.Emit(0)
			g.topbar.gameBtn.SetUpIcon(res.SmileSprites[0])
		case GamePlay:
		case GamePause:
			g.game.Sw.Stop()
		case GameWin:
			g.game.MarkGameAfterEnd()
			g.topbar.gameBtn.SetUpIcon(res.SmileSprites[3])
		case GameGameOver:
			g.game.MarkGameAfterEnd()
			g.topbar.gameBtn.SetUpIcon(res.SmileSprites[4])
		}
	})

	g.game.CountFlagged.Connect(func(data int) {
		a, b := g.game.CountFlags()
		g.topbar.lblLeft.SetText(strconv.Itoa(a) + "/" + strconv.Itoa(b))
	})

	g.game.Conf.duration.Connect(func(data time.Duration) {
		g.topbar.lblRight.SetText(eui.FormatSmartDuration(data, false))
	})

	g.board.Layout()
}

func (g *GameView) Reset() {
	g.game.Reset()
	for i, v := range g.board.Children() {
		v.(*eui.Button).SetText(g.game.Field.Cell(i).String())
	}
}

func (g *GameView) HardReset() {
	dim := g.game.Conf.dim
	g.New(dim)
	for i, v := range g.board.Children() {
		v.(*eui.Button).SetText(g.game.Field.Cell(i).String())
	}
}

func (g *GameView) Update() {
	if g.game.Sw.IsRun() {
		g.game.Conf.duration.Emit(g.game.Sw.Duration())
	}
	if g.IsHidden() && g.game.Sw.IsRun() {
		g.game.State.Emit(GamePause)
	}
	bgUp := colornames.Silver
	bgDown := colornames.Lightgray
	bgActive := colornames.Yellow
	switch g.game.State.Value() {
	case GamePlay:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			for i, v := range g.board.Children() {
				btn := v.(*eui.Button)
				if btn.IsPressed() {
					btn.SetBg(bgActive)
					x, y := g.game.Pos(i)
					g.game.GetNeighbours(x, y, func(i1, i2 int) {
						idx := g.game.Idx(i1, i2)
						if g.game.Cell(i1, i2).IsClosed() {
							g.board.Children()[idx].(*eui.Button).SetBg(bgActive)
						}
					})
				}
			}
			g.topbar.gameBtn.SetUpIcon(res.SmileSprites[2])
		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			for i, v := range g.board.Children() {
				cell := g.game.Field.Cell(i)
				btn := v.(*eui.Button)
				if btn.Bg() == bgActive {
					if cell.IsOpen() {
						btn.SetBg(bgDown)
					} else {
						btn.SetBg(bgUp)
					}
				}
			}
			g.topbar.gameBtn.SetUpIcon(res.SmileSprites[0])
		}
	case GameWin, GameGameOver, GamePause:
		for i, v := range g.board.Children() {
			cell := g.game.Field.Cell(i)
			btn := v.(*eui.Button)
			if btn.Bg() == bgActive {
				if cell.IsClosed() {
					btn.SetBg(bgUp)
				} else if cell.IsOpen() {
					btn.SetBg(bgDown)
				}
			}
		}
	}
}

type DialogGameEnd struct {
	*eui.Container
	title        *eui.Label
	lastGameIcon *eui.Icon
}

func NewDialogGameEnd(fn func(*eui.Button)) *DialogGameEnd {
	c := &DialogGameEnd{Container: eui.NewContainer(eui.NewLayoutVerticalPercent([]int{10, 80, 10}, 5))}
	contTitle := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{10, 90}, 5))
	contTitle.Add(eui.NewButton(QuitDialog, fn))
	c.title = eui.NewLabel("")
	contTitle.Add(c.title)
	c.Add(contTitle)
	c.lastGameIcon = eui.NewIcon(res.SmileSprites[0])
	c.Add(c.lastGameIcon)
	contBtns := eui.NewContainer(eui.NewHBoxLayout(5))
	contBtns.Add(eui.NewButton(PlayAgain, fn))
	contBtns.Add(eui.NewButton(PlayRepeat, fn))
	c.Add(contBtns)
	c.Add(eui.NewGridBackground(50))
	return c
}

func NewDialogCustom(diff map[string]Dim, fn func(*eui.Button)) *eui.Container {
	c := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{10, 30, 30, 30}, 5))
	contTitle := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{10, 90}, 5))
	contTitle.Add(eui.NewButton(QuitDialog, fn))
	title := eui.NewLabel("Select:")
	contTitle.Add(title)
	c.Add(contTitle)

	gn := func(a, b int) []int {
		arr := make([]int, 0)
		for i := a; i <= b; i++ {
			arr = append(arr, i)
		}
		return arr
	}

	cRow := eui.NewSpinBox(
		"Row",
		gn(5, 50),
		0,
		func(value int) string { return strconv.Itoa(value) },
		func(a, b int) int { return a - b },
		func(a, b int) bool { return a == b },
		true,
		0)
	c.Add(cRow)

	cCol := eui.NewSpinBox(
		"Column",
		gn(5, 30),
		0,
		func(value int) string { return strconv.Itoa(value) },
		func(a, b int) int { return a - b },
		func(a, b int) bool { return a == b },
		true,
		0)
	c.Add(cCol)

	cMines := eui.NewSpinBox(
		"Percent Mines",
		gn(5, 30),
		0,
		func(value int) string { return strconv.Itoa(value) },
		func(a, b int) int { return a - b },
		func(a, b int) bool { return a == b },
		true,
		0)
	c.Add(cMines)

	result := func() string {
		row := cRow.SelectedValue.Value()
		col := cCol.SelectedValue.Value()
		perc := cMines.SelectedValue.Value()
		mines := utils.ValueFromPercent(float64(perc), float64(row*col))
		diff[Custom] = NewDim(row, col, int(mines))
		return strconv.Itoa(row) + ":" + strconv.Itoa(col) + ":" + strconv.Itoa(int(mines))
	}

	cRow.SelectedValue.ConnectAndFire(func(data int) {
		title.SetText("Selected:" + result())
	}, 5)

	cCol.SelectedValue.ConnectAndFire(func(data int) {
		title.SetText("Selected:" + result())
	}, 5)

	cMines.SelectedValue.ConnectAndFire(func(data int) {
		title.SetText("Selected:" + result())
	}, 5)

	return c
}

func NewDialogSelectGame(fn func(*eui.Button)) *eui.Container {
	c := eui.NewContainer(eui.NewSquareGridLayout(2, 2, 20))
	for _, txt := range []string{Easy, Medium, Hard, Custom} {
		c.Add(eui.NewButton(txt, fn))
	}
	return c
}

func main() {
	eui.Init(eui.GetUi().SetTitle(Title).SetSize(800, 600))
	eui.Run(func() *eui.Scene {
		var (
			gameView                   *GameView
			dialogSelect, dialogCustom eui.Drawabler
			dialogGameEnd              *DialogGameEnd
		)
		s := eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 95}, 5))
		s.Add(eui.NewTopBar(Title, nil).SetUseStopwatch())

		diff := map[string]Dim{
			Easy:   NewDim(10, 10, 10),
			Medium: NewDim(15, 15, 40),
			Hard:   NewDim(30, 15, 99),
			Custom: NewDim(0, 0, 0),
		}

		runDialogSelect := func() {
			gameView.Hide()
			dialogGameEnd.Hide()
			dialogSelect.Show()
		}
		runGame := func() {
			gameView.Enable()
			gameView.Show()
			dialogSelect.Hide()
			dialogCustom.Hide()
			dialogGameEnd.Hide()
		}
		runCustomDialog := func() {
			dialogGameEnd.Hide()
			dialogSelect.Hide()
			dialogCustom.Show()
		}
		runDialogGameEnd := func() {
			gameView.topbar.Hide()
			gameView.board.Disable()
			w0, h0 := eui.GetUi().Size()
			r := gameView.board.Rect()
			cameraRect := image.Rect(r.X, r.Y, r.Right(), r.Bottom())
			contentImg := ebiten.NewImage(w0, h0)
			gameView.board.Draw(contentImg)
			gameView.Hide()
			dialogGameEnd.lastGameIcon.SetIcon(contentImg.SubImage(cameraRect).(*ebiten.Image))
			dialogGameEnd.Show()
		}

		gameView = NewGameView(func(b *eui.Button) {
			switch b.Text() {
			case QuitDialog:
				runDialogSelect()
			case PlayAgain:
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
					gameView.Reset()
				}
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
					gameView.HardReset()
				}
			}
		})

		gameView.game.State.Connect(func(data GameState) {
			switch data {
			case GameWin:
				dialogGameEnd.title.SetText("Game Won!")
				runDialogGameEnd()
			case GameGameOver:
				dialogGameEnd.title.SetText("Game Lost!")
				runDialogGameEnd()
			}
		})

		dialogSelect = NewDialogSelectGame(func(b *eui.Button) {
			switch b.Text() {
			case Easy:
				gameView.New(diff[Easy])
				runGame()
			case Medium:
				gameView.New(diff[Medium])
				runGame()
			case Hard:
				gameView.New(diff[Hard])
				runGame()
			case Custom:
				runCustomDialog()
			}
		})

		dialogCustom = NewDialogCustom(diff, func(b *eui.Button) {
			switch b.Text() {
			case QuitDialog:
				gameView.New(diff[Custom])
				runGame()
			}
		})

		dialogGameEnd = NewDialogGameEnd(func(b *eui.Button) {
			switch b.Text() {
			case QuitDialog:
				runDialogSelect()
			case PlayAgain:
				gameView.HardReset()
				runGame()
			case PlayRepeat:
				gameView.Reset()
				runGame()
			}
		})

		board := eui.NewContainer(eui.NewStackLayout(5))
		board.Add(dialogSelect)
		board.Add(gameView)
		board.Add(dialogCustom)
		board.Add(dialogGameEnd)
		s.Add(board)
		s.Add(eui.NewGridBackground(100))
		return s
	}())
}
