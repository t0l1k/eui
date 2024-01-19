package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
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
)

type CellData struct {
	state string
	pos   *eui.PointInt
}

func NewCellData(state string, pos *eui.PointInt) *CellData {
	return &CellData{state: state, pos: pos}
}

type CellState struct {
	eui.SubjectBase
}

func NewCellState(state string, pos *eui.PointInt) *CellState {
	c := &CellState{}
	c.SetValue(NewCellData(state, pos))
	return c
}

type Cell struct {
	state       *CellState
	pos         *eui.PointInt
	sym         int
	open, match bool
}

func NewCell(x, y int) *Cell {
	pos := eui.NewPointInt(x, y)
	return &Cell{
		pos:   pos,
		state: NewCellState(CellStateClosed, pos)}
}

func (c *Cell) Pos() (int, int)    { return c.pos.X, c.pos.Y }
func (c *Cell) SetPos(x, y int)    { c.pos = eui.NewPointInt(x, y) }
func (c *Cell) SetValue(value int) { c.sym = value }
func (c *Cell) IsOpen() bool       { return c.open }
func (c *Cell) IsMatch() bool      { return c.match }

func (c *Cell) SetMatch() {
	c.match = true
	c.state.SetValue(NewCellState(CellStateMatch, c.pos))
}

func (c *Cell) Reset() {
	c.match = false
	c.open = false
	c.state.SetValue(NewCellState(CellStateClosed, c.pos))
}

func (c *Cell) Open() {
	if c.match {
		return
	}
	c.open = !c.open
	if c.open {
		c.state.SetValue(NewCellState(CellStateOpen, c.pos))
	} else {
		c.state.SetValue(NewCellState(CellStateClosed, c.pos))
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

type Field struct {
	field, moves []*Cell
	dim          *eui.PointInt
	ClickCount   int
}

func NewField() *Field {
	f := &Field{dim: eui.NewPointInt(3, 2)}
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
}

func (f *Field) ResetGame() {
	f.moves = nil
	f.ClickCount = 0
	for _, v := range f.field {
		v.Reset()
	}
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
	eui.View
	btn   *eui.Button
	field *Field
}

func NewCellIcon(field *Field, f func(b *eui.Button)) *CellIcon {
	c := &CellIcon{}
	c.SetupView()
	c.field = field
	c.btn = eui.NewButton(CellClosed, f)
	c.Add(c.btn)
	c.Setup(f)
	return c
}

func (c *CellIcon) Setup(f func(b *eui.Button)) {
	c.btn.SetupButton(CellClosed, f)
	c.btn.Bg(eui.Teal)
	c.btn.Fg(eui.Yellow)
}

func (c *CellIcon) UpdateData(value interface{}) {
	switch v := value.(type) {
	case *CellState:
		switch v.Value().(*CellData).state {
		case CellStateClosed:
			c.btn.SetText(CellClosed)
		case CellStateOpen:
			cell := c.field.GetCell(v.Value().(*CellData).pos.X, v.Value().(*CellData).pos.Y)
			c.btn.SetText(cell.String())
		case CellStateMatch:
			cell := c.field.GetCell(v.Value().(*CellData).pos.X, v.Value().(*CellData).pos.Y)
			c.btn.SetText(cell.String())
			c.btn.Bg(eui.GreenYellow)
			c.btn.Fg(eui.Blue)
			c.btn.Disable()
		}
	}
}

type Board struct {
	eui.View
	dialog    *Dialog
	field     *Field
	varArea   *eui.StringVar
	layout    *eui.GridLayoutRightDown
	stopwatch *eui.Stopwatch
	inGame    bool
}

func NewBoard() *Board {
	b := &Board{}
	b.SetupView()
	b.dialog = NewDialog("Выбор игры", func(btn *eui.Button) {
		if btn.GetText() == bNew {
			b.NewGame()
		}
		if btn.GetText() == bReset {
			b.Reset()
		}
		if btn.GetText() == bNext {
			b.NextLevel()
		}
		b.dialog.Visible(false)
	})
	b.dialog.Visible(false)
	b.Add(b.dialog)
	b.varArea = eui.NewStringVar("")
	b.field = NewField()
	b.layout = eui.NewGridLayoutRightDown(b.field.Dim())
	b.stopwatch = eui.NewStopwatch()
	b.NewGame()
	return b
}

func (b *Board) NewGame() {
	b.inGame = false
	b.stopwatch.Reset()
	b.layout.Container = nil
	b.field.NewGame()
	for i := 0; i < len(b.field.field); i++ {
		btn := NewCellIcon(b.field, b.gameLogic)
		x, y := b.field.pos(i)
		cell := b.field.GetCell(x, y)
		cell.state.Attach(btn)
		b.layout.Add(btn)
	}
	b.layout.SetDim(b.field.Dim())
}

func (b *Board) Reset() {
	b.inGame = false
	b.field.ResetGame()
	b.stopwatch.Reset()
	for _, v := range b.layout.Container {
		v.(*CellIcon).btn.Enable()
		v.(*CellIcon).btn.Bg(eui.Teal)
		v.(*CellIcon).btn.Fg(eui.Yellow)
	}
}

func (b *Board) gameLogic(c *eui.Button) {
	for i, v := range b.layout.Container {
		if v.(*CellIcon).btn == c {
			x, y := b.field.pos(i)
			if c.IsMouseDownLeft() {
				if !b.inGame {
					b.stopwatch.Start()
					b.inGame = true
				}
				if b.inGame {
					b.field.Open(x, y)
				}
			}
		}
	}
}

func (b *Board) NextLevel() {
	b.field.NextLevel()
	b.NewGame()
}

func (b *Board) Update(dt int) {
	if b.inGame && b.field.IsWin() {
		b.stopwatch.Stop()
		b.inGame = false
		b.dialog.Visible(true)
		b.dialog.lblTitle.SetText("Победа!")
	}
	str := fmt.Sprintf("Время:[%v] Нажатий: %v Размер поля: %v", b.stopwatch, b.field.ClickCount, b.field.dim.String())
	b.varArea.SetValue(str)
	for _, v := range b.Container {
		v.Update(dt)
	}
	if b.dialog.IsVisible() {
		return
	}
	for _, v := range b.layout.Container {
		v.Update(dt)
	}
}

func (b *Board) Draw(surface *ebiten.Image) {
	for _, v := range b.Container {
		v.Draw(surface)
	}
	if b.dialog.IsVisible() {
		return
	}
	for _, v := range b.layout.Container {
		v.Draw(surface)
	}
}

func (b *Board) Resize(rect []int) {
	b.View.Resize(rect)
	b.layout.SetCellMargin(int(float64(b.GetRect().GetLowestSize()) * 0.008))
	b.layout.Resize(rect)
	w0, h0 := b.GetRect().Size()
	x, y := b.GetRect().Pos()
	w := w0 / 2
	h := h0 / 2
	x += (w0 - w) / 2
	y += (h0 - h) / 2
	b.dialog.Resize([]int{x, y, w, h})
	b.Dirty(true)
}

type Dialog struct {
	eui.View
	btnHide, btnNew, btnReset, btnNext *eui.Button
	lblTitle, lblArea                  *eui.Text
	dialFunc                           func(d *eui.Button)
}

func NewDialog(title string, f func(d *eui.Button)) *Dialog {
	t := &Dialog{}
	t.SetupView()
	t.dialFunc = f
	t.lblTitle = eui.NewText(title)
	t.Add(t.lblTitle)
	t.btnHide = eui.NewButton("x", func(b *eui.Button) {
		t.Visible(false)
	})
	t.Add(t.btnHide)
	t.btnNew = eui.NewButton(bNew, f)
	t.Add(t.btnNew)
	t.btnReset = eui.NewButton(bReset, f)
	t.Add(t.btnReset)
	t.btnNext = eui.NewButton(bNext, f)
	t.Add(t.btnNext)
	t.lblArea = eui.NewText("")
	t.Add(t.lblArea)
	return t
}

func (t *Dialog) SetTitle(title string) {
	t.lblTitle.SetText(title)
}

func (t *Dialog) Resize(rect []int) {
	t.View.Resize(rect)
	x, y := t.GetRect().Pos()
	w, h := t.GetRect().W/3, t.GetRect().H/3
	t.lblTitle.Resize([]int{x, y, w*3 - h, h})
	t.btnHide.Resize([]int{x + w*3 - h, y, h, h})
	y += h
	t.lblArea.Resize([]int{x, y, w * 3, h})
	y += h
	t.btnNew.Resize([]int{x, y, w, h})
	x += w
	t.btnReset.Resize([]int{x, y, w, h})
	x += w
	t.btnNext.Resize([]int{x, y, w, h})
	t.Dirty(true)
}

type SceneGame struct {
	eui.SceneBase
	topBar *eui.TopBar
	board  *Board
}

func NewSceneGame() *SceneGame {
	s := &SceneGame{}
	s.topBar = eui.NewTopBar("Найди пару", func(b *eui.Button) {
		s.board.dialog.SetTitle("Выбор игры")
		s.board.dialog.Visible(true)
	})
	s.topBar.SetShowStopwatch()
	s.topBar.SetTitleCoverArea(0.5)
	s.Add(s.topBar)
	s.board = NewBoard()
	s.Add(s.board)
	s.board.varArea.Attach(s.board.dialog.lblArea)
	s.Resize()
	return s
}

func (s *SceneGame) Resize() {
	w, h := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w, h})
	hT := int(float64(rect.GetLowestSize()) * 0.0382)
	s.topBar.Resize([]int{0, 0, w, hT})
	s.board.Resize([]int{hT / 2, hT + hT/2, w - hT, h - hT*3})
}

func NewGameMatch() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Найди пару")
	k := 60
	w, h := 9*k, 6*k
	u.SetSize(w, h)
	u.GetTheme().Set(eui.ViewBg, eui.Black)
	return u
}

func main() {
	eui.Init(NewGameMatch())
	eui.Run(NewSceneGame())
	eui.Quit()
}
