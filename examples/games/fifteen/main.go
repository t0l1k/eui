package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

var Title = "Fifteen game"

type GameState int

const (
	GameStart GameState = iota
	GamePlay
	GamePause
	GameWin
)

func (s GameState) String() string { return []string{"Start", "Play", "Pause", "Win"}[s] }

type Dim struct{ row, column int }

func NewDim(r, c int) Dim            { return Dim{row: r, column: c} }
func (d Dim) Idx(x, y int) int       { return y*d.row + x }
func (d Dim) Pos(idx int) (int, int) { return idx % d.row, idx / d.row }
func (d Dim) IsEdge(x, y int) bool   { return x < 0 || x > d.row-1 || y < 0 || y > d.column-1 }
func (d Dim) Area() int              { return d.row * d.column }
func (d Dim) String() string         { return fmt.Sprintf("Dim:[%v,%v]", d.row, d.column) }

type Cell struct{ *eui.Signal[int] }

func NewCell(value int) *Cell {
	c := &Cell{Signal: eui.NewSignal(func(a, b int) bool { return a == b })}
	c.Emit(value)
	return c
}
func (c *Cell) Move(value int) { c.Signal.Emit(value) }
func (c Cell) Value() int      { return c.Signal.Value() }
func (c Cell) IsBlank() bool   { return c.Value() == 0 }
func (c Cell) String() string {
	switch c.Value() {
	case 0:
		return "*"
	default:
		return strconv.Itoa(c.Value())
	}
}

type Field []*Cell

func NewField(dim Dim) Field {
	f := Field{}
	for y := 0; y < dim.column; y++ {
		for x := 0; x < dim.row; x++ {
			f = append(f, NewCell(dim.Idx(x, y)+1))
		}
	}
	f[len(f)-1] = NewCell(0)
	return f
}
func (f Field) Cell(idx int) *Cell { return f[idx] }
func (f Field) Blank() (int, *Cell) {
	for i, cell := range f {
		if cell.IsBlank() {
			return i, cell
		}
	}
	return -1, nil
}
func (f *Field) Swap(a, b int) {
	aValue := f.Cell(a).Value()
	bValue := f.Cell(b).Value()
	f.Cell(a).Move(bValue)
	f.Cell(b).Move(aValue)
}

type Dir int

const (
	Up Dir = iota
	Down
	Left
	Right
)

func (d Dir) String() string { return []string{"Up", "Dowm", "Left", "Right"}[d] }

type Game struct {
	state  *eui.Signal[GameState]
	field  Field
	dim    Dim
	clicks int
	sw     *eui.Stopwatch
}

func NewGame(dim Dim) *Game {
	g := &Game{dim: dim, field: NewField(dim), state: eui.NewSignal(func(a, b GameState) bool { return a == b }), sw: eui.NewStopwatch()}
	return g
}

func (g *Game) Move(dir Dir) {
	idx, _ := g.field.Blank()
	x, y := g.dim.Pos(idx)
	switch dir {
	case Up:
		if !g.dim.IsEdge(x, y-1) {
			y--
		}
	case Down:
		if !g.dim.IsEdge(x, y+1) {
			y++
		}
	case Left:
		if !g.dim.IsEdge(x-1, y) {
			x--
		}
	case Right:
		if !g.dim.IsEdge(x+1, y) {
			x++
		}
	}
	newIdx := g.dim.Idx(x, y)
	if idx == newIdx {
		return
	}
	g.field.Swap(idx, newIdx)
	g.clicks++
}

func (g *Game) Shuffle() *Game {
	for i := 0; i < g.dim.Area()*3; i++ {
		d := rand.Intn(4)
		g.Move(Dir(d))
	}
	g.clicks = 0
	return g
}

func (g *Game) IsWin() bool {
	for i := 0; i < len(g.field)-1; i++ {
		if g.field.Cell(i).Value() != i+1 {
			return false
		}
	}
	return true
}

func (g *Game) WinLabel() string {
	s := strings.Builder{}
	s.WriteString("Победа за ")
	s.WriteString(eui.FormatSmartDuration(g.sw.Duration(), true))
	s.WriteString(" секунд, кликов ")
	s.WriteString(strconv.Itoa(g.clicks))
	return s.String()
}

func (f *Game) GetNeighthors(x, y int, fn func(int, int)) {
	for dir := 0; dir < 4; dir++ {
		var dx, dy int
		switch dir {
		case 0:
			dx--
		case 1:
			dx++
		case 2:
			dy--
		case 3:
			dy++
		}
		nx := x + dx
		ny := y + dy
		if !f.dim.IsEdge(nx, ny) {
			fn(nx, ny)
		}
	}
}

func (g *Game) String() string {
	s := strings.Builder{}
	s.WriteString(g.dim.String())
	s.WriteString("\n")
	for y := 0; y < g.dim.column; y++ {
		for x := 0; x < g.dim.row; x++ {
			cell := g.field.Cell(g.dim.Idx(x, y))
			s.WriteString(fmt.Sprintf("%3v", cell.String()))
		}
		s.WriteString("\n")
	}
	return s.String()
}

func main() {
	eui.Init(eui.GetUi().SetTitle(Title).SetSize(640, 400))

	eui.Run(func() *eui.Scene {
		var (
			game    *Game
			dim     Dim
			dimIdx  int
			dims    []Dim = []Dim{NewDim(2, 2), NewDim(3, 3), NewDim(4, 4), NewDim(5, 5)}
			btnDims *eui.Button
		)
		scene := eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 90, 5}, 5))

		board := eui.NewContainer(eui.NewSquareGridLayout(2, 2, 5))

		boardReset := func() {
			dim = dims[dimIdx]
			game = NewGame(dim)
			btnDims.SetText(dim.String())

			board.ResetContainer()
			board.SetLayout(eui.NewSquareGridLayout(dim.row, dim.column, 5))

			for idx0, cell := range game.field {
				board.Add(eui.NewButton(cell.String(), func(b *eui.Button) {
					switch game.state.Value() {
					case GameStart:
						game.state.Emit(GamePlay)
					case GamePlay:
						x, y := game.dim.Pos(idx0)
						game.GetNeighthors(x, y, func(i1, i2 int) {
							newIdx := game.dim.Idx(i1, i2)
							newCell := game.field.Cell(newIdx)
							if newCell.IsBlank() {
								game.field.Swap(newIdx, game.dim.Idx(x, y))
								game.clicks++
							}
						})
						if game.IsWin() {
							game.state.Emit(GameWin)
						}
					case GameWin:
						game.state.Emit(GameStart)
					}
				}))

				cell.Connect(func(value int) {
					btn := board.Children()[idx0].(*eui.Button)
					switch value {
					case 0:
						btn.SetText("")
						btn.Disable()
						btn.SetBg(colornames.Navy)
					default:
						btn.SetText(cell.String())
						btn.Enable()
						btn.SetBg(colornames.Brown)
						btn.SetFg(colornames.Yellow)
					}
				})

				game.state.Connect(func(data GameState) {
					switch data {
					case GameStart:
						game.state.Emit(GamePlay)
					case GamePlay:
						game.clicks = 0
						game.sw.Reset()
						game.Shuffle()
						game.sw.Start()
					case GameWin:
						game.sw.Stop()
						eui.NewSnackBar(game.WinLabel()).ShowTime(3 * time.Second)
					}
				})
			}
			if !board.Rect().IsEmpty() { // Иначе паника при первом запуске, не инициализированы размеры компонентов ещё
				board.Layout()
			}
		}

		nextDim := func() {
			dimIdx++
			if len(dims) == dimIdx {
				dimIdx = 0
			}
			boardReset()
		}

		bottomBar := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{30, 30, 30, 10}, 1))
		btnDims = eui.NewButton(dim.String(), func(b *eui.Button) {
			nextDim()
		})
		lblState := eui.NewLabel("")
		lblClicks := eui.NewLabel("")
		lblSw := eui.NewLabel("")
		bottomBar.Add(btnDims)
		bottomBar.Add(lblState)
		bottomBar.Add(lblClicks)
		bottomBar.Add(lblSw)

		eui.GetUi().TickListener().Connect(func(data eui.Event) {
			lblState.SetText(game.state.Value().String())
			lblClicks.SetText(strconv.Itoa(game.clicks))
			lblSw.SetText(eui.FormatSmartDuration(game.sw.Duration(), true))
		})

		eui.GetUi().KeyboardListener().Connect(func(data eui.Event) {
			kd := data.Value.(eui.KeyboardData)

			if kd.IsPressed(ebiten.KeySpace) {
				game.state.Emit(GamePlay)
			}

			if kd.IsPressed(ebiten.KeyEnter) {
				nextDim()
			}

			switch game.state.Value() {
			case GameStart:
			case GamePlay:
				if kd.IsPressed(ebiten.KeyArrowLeft) {
					game.Move(Left)
				}
				if kd.IsPressed(ebiten.KeyArrowRight) {
					game.Move(Right)
				}
				if kd.IsPressed(ebiten.KeyArrowUp) {
					game.Move(Up)
				}
				if kd.IsPressed(ebiten.KeyArrowDown) {
					game.Move(Down)
				}

				if game.IsWin() {
					game.state.Emit(GameWin)
				}
			}
		})

		scene.Add(eui.NewTopBar(Title, nil).SetShowStoppwatch(true).SetUseStopwatch())
		scene.Add(board)
		scene.Add(bottomBar)

		boardReset()

		return scene
	}())

	eui.Quit(func() {})
}
