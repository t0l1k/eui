package app

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/mines"
	"github.com/t0l1k/eui/examples/minesweeper/res"
)

type SceneGame struct {
	ui.ContainerDefault
	topBar            *TopBar
	game              *Game
	lblMines, lblTime *ui.Label
	btnStatus         *ui.ButtonIcon
}

func NewSceneGame(r, c, m int) *SceneGame {
	s := &SceneGame{}
	s.topBar = NewTopBar("")
	s.Add(s.topBar)
	s.game = newGame(r, c, m)
	s.Add(s.game)
	rect := []int{0, 0, 1, 1}
	bg := ui.Yellow
	fg := ui.Black
	s.lblMines = ui.NewLabel(""+strconv.Itoa(s.game.field.GetTotalMines()), rect, bg, fg)
	s.Add(s.lblMines)
	s.lblTime = ui.NewLabel("00:00", rect, bg, fg)
	s.Add(s.lblTime)
	s.btnStatus = ui.NewButtonIcon([]*ebiten.Image{res.Smiles[0], res.Smiles[1]}, rect, func(b *ui.ButtonIcon) {
		if b.IsMouseDownLeft() {
			s.game.New()
		} else if b.IsMouseDownRight() {
			s.game.Reset()
		} else if b.IsMouseDownMiddle() {
			s.game.field.RestoreGame()
			s.game.redraw()
			s.game.timer.Resume()
		}
	})
	s.Add(s.btnStatus)
	return s
}

func (s *SceneGame) Entered() {
	s.Resize()
}

func (s *SceneGame) Update(dt int) {
	s.lblTime.SetText(s.game.timer.String())
	str := strconv.Itoa(s.game.field.GetLeftMines()) + "/" + strconv.Itoa(s.game.field.GetTotalMines())
	s.lblMines.SetText(str)
	switch s.game.field.GetState() {
	case mines.GameStart:
		s.btnStatus.SetIconRelesed(res.Smiles[0])
	case mines.GamePlay:
		if s.game.fEmoji {
			s.btnStatus.SetIconRelesed(res.Smiles[2])
		} else {
			s.btnStatus.SetIconRelesed(res.Smiles[0])
		}
	case mines.GameWin:
		s.btnStatus.SetIconRelesed(res.Smiles[3])
	case mines.GameOver:
		s.btnStatus.SetIconRelesed(res.Smiles[4])
	}
	for _, c := range s.Container {
		c.Update(dt)
	}
}

func (s *SceneGame) Draw(surface *ebiten.Image) {
	surface.Fill(ui.Black)
	for _, c := range s.Container {
		c.Draw(surface)
	}
}

func (s *SceneGame) Resize() {
	s.topBar.Resize()
	s.game.Resize()
	w, h := ebiten.WindowSize()
	hTop := int(float64(h) * 0.05) // topbar height
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	x := rect.CenterX() - hTop*3
	y := 0
	r := []int{x, y, hTop * 2, hTop}
	s.lblMines.Resize(r)
	x = rect.CenterX() + hTop
	y = 0
	r = []int{x, y, hTop * 2, hTop}
	s.lblTime.Resize(r)
	x = rect.CenterX() - hTop/2
	y = 0
	r = []int{x, y, hTop, hTop}
	s.btnStatus.Resize(r)
}

func (s *SceneGame) Close() {
	for _, c := range s.Container {
		c.Close()
	}
}
