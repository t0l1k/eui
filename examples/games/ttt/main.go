package main

import (
	"fmt"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/ttt/pos"
)

const title = "Крестики-Нолики с minimax ai"

func main() {
	eui.Init(func() *eui.Ui {
		u := eui.GetUi()
		u.SetTitle(title)
		k := 2
		w, h := 320*k, 200*k
		u.SetSize(w, h)
		return u
	}())
	eui.Run(
		func() *eui.Scene {
			s := eui.NewScene(eui.NewLayoutVerticalPercent([]int{10, 80, 10}, 5))

			game := pos.NewPosititon(3)

			lblStatus := eui.NewText("")
			lblStatus.SetText("Turn:" + game.GetTurn())

			board := eui.NewContainer(eui.NewGridLayout(3, 3, 5))
			for i := range 9 {
				board.Add(eui.NewButton("", func(b *eui.Button) {
					if game.IsGameEnd() || game.GetBoard()[i] != pos.TurnEmpty {
						return
					}
					game.Move(i)
					board.Children()[i].(*eui.Button).SetText(game.GetNextTurn())
					if !game.IsGameEnd() {
						best := game.BestMove()
						game.Move(best)
						board.Children()[best].(*eui.Button).SetText(game.GetNextTurn())
						lblStatus.SetText("Turn:" + game.GetTurn())
					}
					if game.IsGameEnd() {
						str := ""
						if game.IsGameEnd() {
							if game.IsWinFor(pos.TurnX) {
								str += "Won " + string(pos.TurnX)
							} else if game.IsWinFor(pos.TurnO) {
								str += "Won " + string(pos.TurnO)
							} else if game.Blanks() == 0 {
								str = "Draw!"
							}
						}
						lblStatus.SetText(str)
					}
					fmt.Println(game)
				}))
			}

			statusLine := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{30, 70}, 5))
			statusLine.Add(eui.NewButton("Reset game", func(b *eui.Button) {
				game.Reset()
				lblStatus.SetText("Turn:" + game.GetTurn())
				for _, v := range board.Children() {
					v.(*eui.Button).SetText(string(pos.TurnEmpty))
				}
			}))
			statusLine.Add(lblStatus)

			s.Add(eui.NewTopBar(title, nil))
			s.Add(board)
			s.Add(statusLine)
			return s
		}())
	eui.Quit(func() {})
}
