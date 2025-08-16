package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

const title = "Найди пару"

const (
	bNew   = "Новая"
	bReset = "Повторить"
	bNext  = "Следующий уровень"
	bCont  = "Продолжить"
	bQuit  = "X"
)

var dialogMenu = []string{bNew, bReset, bNext, bCont}

type GameState int

const (
	GameStart GameState = iota
	GamePlay
	GamePause
	GameWin
)

func (s GameState) String() string {
	return [...]string{
		"Start",
		"Play",
		"Pause",
		"Win",
	}[s]
}

type CellState int

const (
	CellClose CellState = iota
	CellOpen
	CellMatch
)

func (a CellState) Eq(b CellState) bool { return a == b }

type Cell struct {
	State *eui.Signal[CellState]
	Value int
}

func NewCell(value int) *Cell {
	return &Cell{State: eui.NewSignal(func(a, b CellState) bool { return a.Eq(b) }), Value: value}
}
func (c *Cell) IsOpen() bool    { return c.State.Value().Eq(CellOpen) }
func (c *Cell) IsClosed() bool  { return c.State.Value().Eq(CellClose) }
func (c *Cell) IsMatch() bool   { return c.State.Value().Eq(CellMatch) }
func (a *Cell) Eq(b *Cell) bool { return a.Value == b.Value }
func (c *Cell) Open() *Cell     { c.State.Emit(CellOpen); return c }
func (c *Cell) Close() *Cell    { c.State.Emit(CellClose); return c }
func (c *Cell) Match() *Cell    { c.State.Emit(CellMatch); return c }
func (c *Cell) String() string {
	switch c.State.Value() {
	case CellClose:
		return fmt.Sprintf("%v", "*")
	default:
		return fmt.Sprintf("%v", c.Value)
	}
}

type Game struct {
	state        *eui.Signal[GameState]
	status       *eui.Signal[string]
	dim          eui.Point[int]
	field, moves []*Cell
	clickCount   int
	sw           *eui.Stopwatch
	completed    bool
}

func NewGame() *Game {
	return &Game{dim: eui.NewPoint(3, 2), state: eui.NewSignal(func(a, b GameState) bool { return a == b }), status: eui.NewSignal(func(a, b string) bool { return a == b }), sw: eui.NewStopwatch()}
}
func (g *Game) New() {
	g.state.Emit(GameStart)
	g.field = nil
	value, i := 1, 0
	for range g.dim.X * g.dim.Y {
		g.field = append(g.field, NewCell(value))
		i++
		if i > 1 {
			value++
			i = 0
		}
	}
	g.shuffle()
	g.clickCount = 0
	g.completed = false
	g.moves = nil
	g.sw.Reset()
	g.status.Emit(g.gameStatus())
	log.Println("Game:New:", g.state.Value(), g.field)
}

func (g *Game) Reset() {
	g.state.Emit(GameStart)
	for y := range g.dim.Y {
		for x := range g.dim.X {
			g.field[g.idx(x, y)].Close()
		}
	}
	g.clickCount = 0
	g.moves = nil
	g.sw.Reset()
	g.status.Emit(g.gameStatus())
	log.Println("Game:Reset:", g.field)
}

func (g *Game) NextLevel() {
	g.dim.X++
	if g.dim.X > g.dim.Y*3 {
		g.dim.Y++
		g.dim.X = g.dim.X / 2
	}
	if (g.dim.X*g.dim.Y)%2 != 0 {
		g.dim.X++
	}
	log.Println("Game:NextLevel:", g.dim)
}

func (g *Game) shuffle() {
	for i := 0; i < 10; i++ {
		for i := range g.field {
			x, y := rand.Intn(g.dim.X), rand.Intn(g.dim.Y)
			tmpX, tmpY := g.pos(i)
			v1 := g.field[g.idx(x, y)].Value
			v2 := g.field[g.idx(tmpX, tmpY)].Value
			g.field[g.idx(x, y)].Value = v2
			g.field[g.idx(tmpX, tmpY)].Value = v1
		}
	}
}

func (g *Game) open(x, y int) {
	cell := g.field[g.idx(x, y)]
	if cell.IsOpen() {
		return
	}
	cell.Open()
	g.moves = append(g.moves, cell)
	g.check()
	g.clickCount++
	g.status.Emit(g.gameStatus())
}

func (g *Game) check() {
	var a1, a2, a3 *Cell

	if len(g.moves) == 2 {
		a1, a2 = g.moves[0], g.moves[1]
		if a1.Eq(a2) {
			a1.Match()
			a2.Match()
		}
	} else if len(g.moves) == 3 {
		a1, a2, a3 = g.moves[0], g.moves[1], g.moves[2]
		if a1.Eq(a2) {
			a1.Match()
			a2.Match()
		} else {
			a1.Close()
			a2.Close()
		}
		g.moves = nil
		g.moves = append(g.moves, a3)
	}
}

func (g *Game) checkWin() {
	for _, v := range g.field {
		if !v.IsMatch() {
			return
		}
	}
	g.completed = true
	g.state.Emit(GameWin)
	g.sw.Stop()
	g.status.Emit(g.gameStatus())
}

func (g *Game) gameStatus() string {
	return fmt.Sprintf("%v Поле:%v Нажатий:%v Время:%v", g.state.Value(), g.dim.String(), g.clickCount, eui.FormatSmartDuration(g.sw.Duration(), true))
}

func (g *Game) Dim() (int, int)        { return g.dim.X, g.dim.Y }
func (g *Game) pos(idx int) (int, int) { return idx % g.dim.X, idx / g.dim.X }
func (g *Game) idx(x, y int) int       { return y*g.dim.X + x }
func (g *Game) Cell(x, y int) *Cell    { return g.field[g.idx(x, y)] }

func (g *Game) String() (result string) {
	result = fmt.Sprintf("\nРазмер поля [%vx%v]\nНажатий:%v\n", g.dim.X, g.dim.Y, g.clickCount)
	for y := 0; y < g.dim.Y; y++ {
		for x := 0; x < g.dim.X; x++ {
			result += fmt.Sprintf("[%.2v(%v)]", g.field[g.idx(x, y)].State.Value(), g.field[g.idx(x, y)].Value)
		}
		result += "\n"
	}
	return result
}

type Board struct {
	*eui.Container
	game *Game
	fn   func(*eui.Button)
}

func NewBoard(fn func(*eui.Button)) *Board {
	b := &Board{Container: eui.NewContainer(eui.NewGridLayout(3, 2, 10)), fn: fn}
	b.game = NewGame()
	log.Println("Board:Init", b.game)
	return b
}

func (b *Board) New() {
	b.game.New()
	b.ResetContainer()

	w, h := b.game.dim.X, b.game.dim.Y
	b.SetLayout(eui.NewGridLayout(w, h, 10))

	for _, cell := range b.game.field {
		icon := eui.NewButton(cell.String(), b.fn)
		b.Add(icon)

		cell.State.Connect(func(state CellState) {
			switch state {
			case CellClose:
				icon.SetText(cell.String())
				icon.Bg(colornames.Gray)
				icon.Fg(colornames.Yellow)

			case CellOpen:
				icon.SetText(cell.String())
				icon.Bg(colornames.Teal)
				icon.Fg(colornames.Yellow)
			case CellMatch:
				icon.SetText(cell.String())
				icon.Bg(colornames.Greenyellow)
				icon.Fg(colornames.Black)
			}
		})
	}

	b.MarkDirty()
	b.Layout()
	log.Println("Board:New", b.game, b.Children())
	for _, v := range b.Children() {
		log.Println("B:", b.Rect(), v.(*eui.Button).Text(), v.Rect())
	}
}

func (b *Board) Reset() {
	b.game.Reset()
	w, h := b.game.dim.X, b.game.dim.Y
	for i := range w * h {
		x, y := b.game.pos(i)
		value := b.game.Cell(x, y).String()
		b.Children()[i].(*eui.Button).SetText(value)
	}
	log.Println("Board:Reset", b.game, b.Children())
	for _, v := range b.Children() {
		log.Println("B:", v.(*eui.Button).Text(), v.Rect())
	}
}

func (b *Board) Update(int) {
	b.game.status.Emit(b.game.gameStatus())
	if b.IsHidden() && b.game.state.Value() == GamePlay {
		b.game.state.Emit(GamePause)
	}
}

type Dialog struct {
	*eui.Container
	msg *eui.Label
}

func NewDialog(title, message string, fn func(*eui.Button)) *Dialog {
	d := &Dialog{Container: eui.NewContainer(eui.NewLayoutVerticalPercent([]int{10, 70, 20}, 1))}
	titleContainer := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{10, 90}, 1))
	titleContainer.Add(eui.NewButton("X", func(b *eui.Button) { eui.GetUi().Pop() }))
	titleLbl := eui.NewLabel(title)
	titleLbl.SetAlign(eui.LabelAlignLeft)
	titleContainer.Add(titleLbl)
	d.Add(titleContainer)
	d.msg = eui.NewLabel(message)
	d.Add(d.msg)
	btnsContainer := eui.NewContainer(eui.NewHBoxLayout(1))
	for _, v := range dialogMenu {
		btnsContainer.Add(eui.NewButton(v, fn))
	}
	d.Add(btnsContainer)
	return d
}

type TopBar struct {
	*eui.Container
	timerVar *eui.Signal[time.Duration]
}

func NewTopbar(title string, fn func(*eui.Button)) *TopBar {
	t := &TopBar{Container: eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{5, 25, 60, 10}, 1))}

	tmLbl := eui.NewLabel("-:--")
	t.timerVar = eui.NewSignal(func(a, b time.Duration) bool { return a == b })
	t.timerVar.ConnectAndFire(func(data time.Duration) {
		tmLbl.SetText(eui.FormatSmartDuration(data, false))
	}, 0)

	btnLbl := "X"
	if fn != nil {
		btnLbl = "Menu"
	} else {
		fn = func(b *eui.Button) {
			eui.GetUi().Pop()
		}
	}
	t.Add(eui.NewButton(btnLbl, fn))
	t.Add(eui.NewLabel(title))
	t.Add(eui.NewDrawable())
	t.Add(tmLbl)
	return t
}
func (t *TopBar) Tick(td eui.TickData) { t.timerVar.Emit(t.timerVar.Value() + td.Duration()) }

func main() {
	eui.Init(func() *eui.Ui {
		return eui.GetUi().SetTitle(title).SetSize(800, 600)
	}())
	eui.Run(func() *eui.Scene {
		var (
			dialog *Dialog
			board  *Board
		)
		s := &eui.Scene{Container: eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 90, 5}, 5))}

		statusLbl := eui.NewLabel("StatusLine")

		board = NewBoard(func(b *eui.Button) {
			for i, v := range board.Children() {
				x, y := board.game.pos(i)
				cell := board.game.Cell(board.game.pos(i))
				icon := v.(*eui.Button)
				if icon == b {
					switch board.game.state.Value() {
					case GameStart:
						board.game.state.Emit(GamePlay)
						board.game.sw.Start()
						board.game.open(x, y)
					case GamePlay:
						board.game.open(x, y)
						board.game.checkWin()
					case GamePause:
						board.game.state.Emit(GamePlay)
						board.game.sw.Start()
					case GameWin:
					}
					log.Println("Pressed:", b.Text(), b.Rect(), icon.Text(), icon.Rect(), i, x, y, cell)
				}
			}
		})
		board.game.state.Connect(func(data GameState) {
			switch data {
			case GameStart:
			case GamePlay:
			case GamePause:
				board.game.sw.Stop()
				board.Hide()
			case GameWin:
				board.Hide()
				dialog.Show()
			}
		})

		board.game.status.Connect(func(data string) {
			statusLbl.SetText(data)
			dialog.msg.SetText(board.game.gameStatus())
		})

		dialog = NewDialog("Select", "Press New Game", func(b *eui.Button) {
			switch b.Text() {
			case bNew:
				board.New()
				dialog.Hide()
				board.Show()
			case bReset:
				board.Reset()
			case bNext:
				if board.game.completed {
					board.game.NextLevel()
					board.New()
				}
			case bCont:
			}
			dialog.Hide()
			board.Show()
		})

		boardContainer := eui.NewContainer(eui.NewStackLayout(5))
		boardContainer.Add(dialog)
		boardContainer.Add(board)
		s.Add(NewTopbar(title, func(b *eui.Button) {
			board.Hide()
			dialog.Show()
		}))
		s.Add(boardContainer)
		s.Add(statusLbl)
		return s
	}())
}
