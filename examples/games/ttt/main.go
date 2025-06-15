package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/ttt/pos"
)

const (
	title = "Крестики-Нолики с minimax ai"
)

type BoardTTT struct {
	eui.DrawableBase
	layout *eui.GridLayoutRightDown
}

func NewBoardTTT(fn func(*eui.Button)) *BoardTTT {
	d := &BoardTTT{}
	d.Visible(true)
	d.layout = eui.NewGridLayoutRightDown(3, 3)
	for i := 0; i < 9; i++ {
		btn := eui.NewButton(string(pos.TurnEmpty), fn)
		d.layout.Add(btn)
	}
	return d
}

func (d *BoardTTT) Reset() {
	for _, v := range d.layout.GetContainer() {
		v.(*eui.Button).SetText(string(pos.TurnEmpty))
	}
}

func (d *BoardTTT) Update(dt int) {
	if !d.IsVisible() {
		return
	}
	d.DrawableBase.Update(dt)
	for _, v := range d.layout.GetContainer() {
		v.Update(dt)
	}
}

func (d *BoardTTT) Draw(surface *ebiten.Image) {
	if !d.IsVisible() {
		return
	}
	d.DrawableBase.Draw(surface)
	for _, v := range d.layout.GetContainer() {
		v.Draw(surface)
	}
}

func (d *BoardTTT) Resize(rect []int) {
	d.Rect(eui.NewRect(rect))
	d.layout.Resize(rect)
	d.ImageReset()
}

type SceneMain struct {
	eui.SceneBase
	topBar    *eui.TopBar
	game      *pos.Posititon
	board     *BoardTTT
	lblStatus *eui.Text
	btnReset  *eui.Button
}

func NewSceneMain() *SceneMain {
	s := &SceneMain{}
	s.game = pos.NewPosititon(3)
	s.topBar = eui.NewTopBar(title, nil)
	s.Add(s.topBar)
	s.board = NewBoardTTT(func(btn *eui.Button) {
		for i, v := range s.board.layout.GetContainer() {
			if v.(*eui.Button) == btn {
				if s.game.IsGameEnd() || s.game.GetBoard()[i] != pos.TurnEmpty {
					return
				}
				s.game.Move(i)
				s.board.layout.GetContainer()[i].(*eui.Button).SetText(s.game.GetNextTurn())
				if !s.game.IsGameEnd() {
					best := s.game.BestMove()
					s.game.Move(best)
					s.board.layout.GetContainer()[best].(*eui.Button).SetText(s.game.GetNextTurn())
					s.lblStatus.SetText("Turn:" + s.game.GetTurn())
				}
				if s.game.IsGameEnd() {
					str := ""
					if s.game.IsGameEnd() {
						if s.game.IsWinFor(pos.TurnX) {
							str += "Won " + string(pos.TurnX)
						} else if s.game.IsWinFor(pos.TurnO) {
							str += "Won " + string(pos.TurnO)
						} else if s.game.Blanks() == 0 {
							str = "Draw!"
						}
					}
					s.lblStatus.SetText(str)
				}
				fmt.Println(s.game)
			}
		}
	})
	s.Add(s.board)
	s.lblStatus = eui.NewText("")
	s.lblStatus.SetText("Turn:" + s.game.GetTurn())
	s.Add(s.lblStatus)
	s.btnReset = eui.NewButton("Reset Game", func(b *eui.Button) {
		s.game.Reset()
		s.board.Reset()
		s.lblStatus.SetText("Turn:" + s.game.GetTurn())
	})
	s.Add(s.btnReset)
	s.Resize()
	return s
}

func (s *SceneMain) Resize() {
	w0, h0 := eui.GetUi().Size()
	hTop := int(float64(h0) * 0.05) // topbar height
	s.topBar.Resize([]int{0, 0, w0, hTop})
	s.board.Resize([]int{hTop, hTop * 2, w0 - hTop*2, h0 - hTop*4})
	s.btnReset.Resize([]int{0, h0 - hTop, hTop * 3, hTop})
	s.lblStatus.Resize([]int{hTop * 3, h0 - hTop, w0 - hTop*3, hTop})
}

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle(title)
	k := 2
	w, h := 320*k, 200*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneMain())
	eui.Quit(func() {})
}
