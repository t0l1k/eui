package main

import (
	"fmt"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/ttt/pos"
)

const (
	title = "Крестики-Нолики с minimax ai"
)

type BoardTTT struct{ *eui.Container }

func NewBoardTTT(fn func(*eui.Button)) *BoardTTT {
	d := &BoardTTT{Container: eui.NewContainer(eui.NewGridLayout(3, 3, 1))}
	d.Visible(true)
	for i := 0; i < 9; i++ {
		btn := eui.NewButton(string(pos.TurnEmpty), fn)
		d.Add(btn)
	}
	return d
}

func (d *BoardTTT) Reset() {
	for _, v := range d.Childrens() {
		v.(*eui.Button).SetText(string(pos.TurnEmpty))
	}
}

// func (d *BoardTTT) Update(dt int) {
// 	if !d.IsVisible() {
// 		return
// 	}
// 	d.Drawable.Update(dt)
// 	for _, v := range d.Childrens() {
// 		v.Update(dt)
// 	}
// }

// func (d *BoardTTT) Draw(surface *ebiten.Image) {
// 	if !d.IsVisible() {
// 		return
// 	}
// 	d.Drawable.Draw(surface)
// 	for _, v := range d.Childrens() {
// 		v.Draw(surface)
// 	}
// }

// func (d *BoardTTT) Resize(rect eui.Rect[int]) {
// 	d.SetRect(rect)
// 	d.ImageReset()
// }

type SceneMain struct {
	*eui.Scene
	topBar    *eui.TopBar
	game      *pos.Posititon
	board     *BoardTTT
	lblStatus *eui.Text
	btnReset  *eui.Button
}

func NewSceneMain() *SceneMain {
	s := &SceneMain{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.game = pos.NewPosititon(3)
	s.topBar = eui.NewTopBar(title, nil)
	s.Add(s.topBar)
	s.board = NewBoardTTT(func(btn *eui.Button) {
		for i, v := range s.board.Childrens() {
			if v.(*eui.Button) == btn {
				if s.game.IsGameEnd() || s.game.GetBoard()[i] != pos.TurnEmpty {
					return
				}
				s.game.Move(i)
				s.board.Childrens()[i].(*eui.Button).SetText(s.game.GetNextTurn())
				if !s.game.IsGameEnd() {
					best := s.game.BestMove()
					s.game.Move(best)
					s.board.Childrens()[best].(*eui.Button).SetText(s.game.GetNextTurn())
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
	s.topBar.Resize(eui.NewRect([]int{0, 0, w0, hTop}))
	s.board.Resize(eui.NewRect([]int{hTop, hTop * 2, w0 - hTop*2, h0 - hTop*4}))
	s.btnReset.Resize(eui.NewRect([]int{0, h0 - hTop, hTop * 3, hTop}))
	s.lblStatus.Resize(eui.NewRect([]int{hTop * 3, h0 - hTop, w0 - hTop*3, hTop}))
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
