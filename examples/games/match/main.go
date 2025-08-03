package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

const (
	GameStart = "start"
	GamePlay  = "play"
	GamePause = "pause"
	GameWin   = "win"
)

const (
	CellClosed      = "*"
	CellStateClosed = "cell closed"
	CellStateOpen   = "cell open"
	CellStateMatch  = "cell match"
)

const (
	bNew   = "Новая"
	bReset = "Повторить"
	bNext  = "Следующий"
	bCont  = "Продолжить"
	bQuit  = "X"
)

type CellData struct {
	state string
	pos   eui.Point[int]
}

func NewCellData(state string, pos eui.Point[int]) *CellData {
	return &CellData{state: state, pos: pos}
}

type Cell struct {
	state       *eui.Signal[*CellData]
	pos         eui.Point[int]
	sym         int
	open, match bool
}

func NewCell(x, y int) *Cell {
	pos := eui.NewPoint(x, y)
	return &Cell{
		pos:   pos,
		state: eui.NewSignal[*CellData]()}
}

func (c *Cell) Pos() (int, int)    { return c.pos.X, c.pos.Y }
func (c *Cell) SetPos(x, y int)    { c.pos = eui.NewPoint(x, y) }
func (c *Cell) SetValue(value int) { c.sym = value }
func (c *Cell) IsOpen() bool       { return c.open }
func (c *Cell) IsMatch() bool      { return c.match }

func (c *Cell) SetMatch() {
	c.match = true
	c.state.Emit(NewCellData(CellStateMatch, c.pos))
}

func (c *Cell) Reset() {
	c.match = false
	c.open = false
	c.state.Emit(NewCellData(CellStateClosed, c.pos))
}

func (c *Cell) Open() {
	if c.match {
		return
	}
	c.open = !c.open
	if c.open {
		c.state.Emit(NewCellData(CellStateOpen, c.pos))
	} else {
		c.state.Emit(NewCellData(CellStateClosed, c.pos))
	}
}

func (c *Cell) String() (result string) {
	if c.open || c.match {
		result = fmt.Sprintf("%v", c.sym)
	} else {
		result = CellClosed
	}
	return result
}

type FieldState struct{ *eui.Signal[string] }

func NewFieldState() *FieldState {
	return &FieldState{Signal: eui.NewSignal(func(a, b string) bool { return a == b })}
}

type Field struct {
	State        *FieldState
	field, moves []*Cell
	dim          eui.Point[int]
	ClickCount   int
}

func NewField() *Field {
	f := &Field{State: NewFieldState(), dim: eui.NewPoint(3, 2)}
	return f
}

func (f *Field) NewGame() {
	f.field = nil
	f.moves = nil
	f.ClickCount = 0
	a := 1
	i := 0
	for y := 0; y < f.dim.Y; y++ {
		for x := 0; x < f.dim.X; x++ {
			c := NewCell(x, y)
			c.SetValue(a)
			f.field = append(f.field, c)
			i++
			if i > 1 {
				i = 0
				a++
			}
		}
	}
	f.shuffle()
	f.State.Emit(GameStart)
	log.Println("Set Game start")
}

func (f *Field) ResetGame() {
	f.moves = nil
	f.ClickCount = 0
	for _, v := range f.field {
		v.Reset()
	}
	f.State.Emit(GameStart)
	log.Println("Set Game start")
}

func (f *Field) NextLevel() {
	f.dim.X++
	if f.dim.X > f.dim.Y*3 {
		f.dim.Y++
		f.dim.X = f.dim.X / 2
	}
	if (f.dim.X*f.dim.Y)%2 != 0 {
		f.dim.X++
	}
	f.NewGame()
}

func (f *Field) shuffle() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		for _, cell := range f.field {
			x, y := rand.Intn(f.dim.X), rand.Intn(f.dim.Y)
			tmpX, tmpY := cell.Pos()
			v1 := f.field[f.idx(x, y)].sym
			v2 := f.field[f.idx(tmpX, tmpY)].sym
			f.field[f.idx(x, y)].sym = v2
			f.field[f.idx(tmpX, tmpY)].sym = v1
		}
	}
}

func (f *Field) Open(x, y int) {
	if f.field[f.idx(x, y)].IsOpen() {
		return
	}
	f.field[f.idx(x, y)].Open()
	f.moves = append(f.moves, f.GetCell(x, y))
	f.checkPair()
	f.ClickCount++
}

func (f *Field) checkPair() {
	var (
		a1, a2, a3 *Cell
	)
	if len(f.moves) == 2 {
		a1, a2 = f.moves[0], f.moves[1]
		if a1.sym == a2.sym {
			a1.SetMatch()
			a2.SetMatch()
		}
	} else if len(f.moves) == 3 {
		a1, a2, a3 = f.moves[0], f.moves[1], f.moves[2]
		if a1.sym == a2.sym {
			a1.SetMatch()
			a2.SetMatch()
		} else {
			a1.Open()
			a2.Open()
		}
		f.moves = nil
		f.moves = append(f.moves, a3)
	}
}

func (f *Field) IsWin() bool {
	for _, v := range f.field {
		if !v.IsMatch() {
			return false
		}
	}
	f.State.Emit(GameWin)
	log.Println("Set game win!")
	return true
}

func (f *Field) Dim() (int, int)        { return f.dim.X, f.dim.GetY() }
func (f *Field) pos(idx int) (int, int) { return idx % f.dim.X, idx / f.dim.X }
func (f *Field) idx(x, y int) int       { return y*f.dim.X + x }
func (f *Field) GetCell(x, y int) *Cell { return f.field[f.idx(x, y)] }

func (f *Field) String() (result string) {
	result = fmt.Sprintf("Размер поля [%vx%v]\nНажатий:%v\n", f.dim.X, f.dim.Y, f.ClickCount)
	for y := 0; y < f.dim.Y; y++ {
		for x := 0; x < f.dim.X; x++ {
			result += fmt.Sprintf("[%.2v]", f.field[f.idx(x, y)].sym)
		}
		result += "\n"
	}
	return result
}

type CellIcon struct {
	*eui.Container
	btn   *eui.Button
	field *Field
}

func NewCellIcon(field *Field, f func(b *eui.Button)) *CellIcon {
	c := &CellIcon{Container: eui.NewContainer(eui.NewVBoxLayout(1))}
	c.field = field
	c.btn = eui.NewButton(CellClosed, f)
	c.Add(c.btn)
	c.Setup(f)
	return c
}

func (c *CellIcon) Setup(f func(b *eui.Button)) {
	c.btn.SetupButton(CellClosed, f)
	c.btn.Bg(colornames.Teal)
	c.btn.Fg(colornames.Yellow)
}
func (d *CellIcon) SetHidden(value bool) {
	d.Traverse(func(c eui.Drawabler) {
		c.SetHidden(value)
		if !value {
			c.Enable()
		} else {
			c.Disable()
		}
	}, false)
}

type Board struct {
	*eui.Container
	field     *Field
	varArea   *eui.Signal[string]
	layout    *eui.Container
	stopwatch *eui.Stopwatch
	bottomLbl *eui.Text
}

func NewBoard() *Board {
	b := &Board{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	b.varArea = eui.NewSignal(func(a, b string) bool { return a == b })
	b.field = NewField()
	b.field.State.Connect(func(value string) {
		switch value {
		case GameStart:
		case GamePlay:
			if !b.stopwatch.IsRun() {
				b.stopwatch.Start()
			}
		case GamePause:
			b.stopwatch.Stop()
			b.SetHidden(true)
		case GameWin:
			b.stopwatch.Stop()
			b.SetHidden(true)
			str := fmt.Sprintf("Время:[%v] Нажатий: %v Размер поля: %v", b.stopwatch, b.field.ClickCount, b.field.dim.String())
			b.varArea.Emit(str)
		}
		log.Println("board got:", value)
	})
	r, c := b.field.Dim()
	b.layout = eui.NewContainer(eui.NewGridLayout(float64(r), float64(c), 1))
	b.Add(b.layout)
	b.stopwatch = eui.NewStopwatch()
	b.bottomLbl = eui.NewText("")
	b.Add(b.bottomLbl)
	b.varArea.Connect(func(data string) {
		b.bottomLbl.SetText(data)
	})
	b.NewGame()
	return b
}

func (b *Board) NewGame() {
	b.stopwatch.Reset()
	b.layout.ResetContainer()
	b.field.NewGame()
	r, c := b.field.Dim()
	b.layout.SetLayout(eui.NewGridLayout(float64(r), float64(c), 1))
	for i := 0; i < len(b.field.field); i++ {
		btn := NewCellIcon(b.field, b.gameLogic)
		x, y := b.field.pos(i)
		cell := b.field.GetCell(x, y)
		cell.state.Connect(func(value *CellData) {
			c := btn
			switch value.state {
			case CellStateClosed:
				c.btn.SetText(CellClosed)
			case CellStateOpen:
				cell := c.field.GetCell(value.pos.X, value.pos.Y)
				c.btn.SetText(cell.String())
			case CellStateMatch:
				cell := c.field.GetCell(value.pos.X, value.pos.Y)
				c.btn.SetText(cell.String())
				c.btn.Bg(colornames.Greenyellow)
				c.btn.Fg(colornames.Blue)
				c.btn.Disable()
			}

		})
		b.layout.Add(btn)
	}
	b.bottomLbl.SetHidden(false)
}

func (b *Board) Reset() {
	b.field.ResetGame()
	b.stopwatch.Reset()
	for _, v := range b.layout.Childrens() {
		v.SetHidden(false)
		cell, ok := v.(*CellIcon)
		if ok {
			cell.btn.Bg(colornames.Teal)
			cell.btn.Fg(colornames.Yellow)
		}
	}
	b.bottomLbl.SetHidden(false)
}

func (b *Board) NextLevel() {
	b.field.NextLevel()
	b.NewGame()
}

func (b *Board) gameLogic(c *eui.Button) {
	for i, v := range b.layout.Childrens() {
		if v.(*CellIcon).btn == c {
			x, y := b.field.pos(i)
			if c.IsMouseDownLeft() {
				switch b.field.State.Value() {
				case GameStart:
					b.field.State.Emit(GamePlay)
					log.Println("Set game play")
					b.stopwatch.Start()
					b.field.Open(x, y)
				case GamePlay:
					b.field.Open(x, y)
					b.field.IsWin()
				}
			}
		}
	}
}
func (d *Board) SetHidden(value bool) {
	d.Traverse(func(c eui.Drawabler) {
		c.SetHidden(value)
		if !value {
			c.Enable()
		} else {
			c.Disable()
		}
	}, false)
	d.MarkDirty()
}

func (b *Board) Update(dt int) {
	b.Container.Update(dt)
	str := fmt.Sprintf("Время:[%v] Нажатий: %v Размер поля: %v", b.stopwatch, b.field.ClickCount, b.field.dim.String())
	b.varArea.Emit(str)
}

func (b *Board) SetRect(rect eui.Rect[int]) {
	b.Container.SetRect(rect)
	hT := int(float64(b.Rect().GetLowestSize()) * 0.05)
	x, y := b.Rect().X, b.Rect().Y
	w, h := b.Rect().W, b.Rect().H-hT
	b.layout.SetRect(eui.NewRect([]int{x, y, w, h}))
	y += h
	h = hT
	b.bottomLbl.SetRect(eui.NewRect([]int{x, y, w, h}))
	b.ImageReset()
}

type Dialog struct {
	*eui.Container
	btnQuit, btnNew, btnReset, btnNext, btnCont *eui.Button
	title, message                              *eui.Text
	dialFunc                                    func(d *eui.Button)
	board                                       *Board
}

func NewDialog(title string, board *Board, f func(d *eui.Button)) *Dialog {
	t := &Dialog{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	t.board = board
	t.dialFunc = f
	t.title = eui.NewText(title)
	t.Add(t.title)
	t.btnQuit = eui.NewButton(bQuit, func(b *eui.Button) {
		eui.GetUi().Pop()
	})
	t.Add(t.btnQuit)
	t.btnNew = eui.NewButton(bNew, f)
	t.Add(t.btnNew)
	t.btnReset = eui.NewButton(bReset, f)
	t.Add(t.btnReset)
	t.btnNext = eui.NewButton(bNext, f)
	t.Add(t.btnNext)
	t.btnCont = eui.NewButton(bCont, f)
	t.Add(t.btnCont)
	t.message = eui.NewText("")
	t.Add(t.message)
	return t
}

func (d *Dialog) SetHidden(value bool) {
	d.Traverse(func(c eui.Drawabler) {
		c.SetHidden(value)
		if !value {
			c.Enable()
		} else {
			c.Disable()
		}
	}, false)
	d.MarkDirty()
}

func (t *Dialog) SetTitle(title string) {
	t.title.SetText(title)
}

func (t *Dialog) SetRect(rect eui.Rect[int]) {
	t.Container.SetRect(rect)
	x, y := t.Rect().Pos()
	w, h := t.Rect().W/4, t.Rect().H/3
	t.title.SetRect(eui.NewRect([]int{x, y, w*4 - h, h}))
	t.btnQuit.SetRect(eui.NewRect([]int{x + w*4 - h, y, h, h}))
	y += h
	t.message.SetRect(eui.NewRect([]int{x, y, w * 4, h}))
	y += h
	t.btnNew.SetRect(eui.NewRect([]int{x, y, w, h}))
	x += w
	t.btnReset.SetRect(eui.NewRect([]int{x, y, w, h}))
	x += w
	t.btnNext.SetRect(eui.NewRect([]int{x, y, w, h}))
	x += w
	t.btnCont.SetRect(eui.NewRect([]int{x, y, w, h}))
	t.ImageReset()
}

type SceneGame struct {
	*eui.Scene
	topBar *eui.TopBar
	board  *Board
	dialog *Dialog
}

func NewSceneGame() *SceneGame {
	s := &SceneGame{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar("Найди пару", func(b *eui.Button) {
		s.dialog.SetTitle("Выбор игры")
		s.dialog.SetHidden(false)
		s.board.SetHidden(true)
		if s.board.field.State.Value() == GamePlay {
			s.board.field.State.Emit(GamePause)
		}
		log.Println("Set Game to pause")
	})
	s.topBar.SetUseStopwatch()
	s.topBar.SetTitleCoverArea(0.5)
	s.Add(s.topBar)

	s.dialog = NewDialog("Выбор игры", s.board, func(btn *eui.Button) {
		if btn.GetText() == bNew {
			s.board.NewGame()
		}
		if btn.GetText() == bReset {
			s.board.Reset()
		}
		if btn.GetText() == bNext {
			s.board.NextLevel()
		}
		if btn.GetText() == bCont {
			if s.board.field.State.Value() == GamePause {
				s.board.field.State.Emit(GamePlay)
			}
			s.board.SetHidden(false)
			log.Println("Set Game continue play")
		}
		s.dialog.SetHidden(true)
	})
	s.dialog.SetHidden(true)
	s.Add(s.dialog)

	s.board = NewBoard()
	s.board.varArea.Connect(func(data string) {
		s.dialog.message.SetText(data)
	})
	s.board.field.State.Connect(func(value string) {
		b := s.dialog
		switch value {
		case GamePlay:
			if !b.IsHidden() {
				b.board.field.State.Emit(GamePause)
				log.Println("Set Game Pause")
			}
		case GamePause:
			b.title.SetText("Пауза!")
			b.SetHidden(false)
		case GameWin:
			b.SetHidden(false)
			b.title.SetText("Победа!")
		}

		log.Println("dialog got:", value)

	})
	s.Add(s.board)
	return s
}

func (s *SceneGame) SetRect(rect eui.Rect[int]) {
	w0, h0 := rect.Size()
	x, y := 0, 0
	w, h := w0, h0
	hT := int(float64(rect.GetLowestSize()) * 0.1)
	h = hT
	s.topBar.SetRect(eui.NewRect([]int{x, y, w, h}))
	x, y = hT/2, hT+hT/2
	w, h = w0-hT, h0-hT*2
	s.board.SetRect(eui.NewRect([]int{x, y, w, h}))
	s.dialog.SetRect(eui.NewRect([]int{x, y, w, h}))
}

func NewGameMatch() *eui.Ui {
	k := 60
	w, h := 9*k, 6*k
	u := eui.GetUi().SetTitle("Найди пару").SetSize(w, h)
	u.GetTheme().Set(eui.ViewBg, colornames.Black)
	return u
}

func main() {
	eui.Init(NewGameMatch())
	eui.Run(NewSceneGame())
	eui.Quit(func() {})
}
