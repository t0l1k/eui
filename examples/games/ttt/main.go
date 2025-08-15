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

			lblStatus := eui.NewLabel("")
			statusGame := eui.NewSignal(func(a, b string) bool { return a == b })
			statusGame.ConnectAndFire(func(data string) {
				lblStatus.SetText(data)
			}, "Turn:"+game.GetTurn())

			board := eui.NewContainer(eui.NewSquareGridLayout(3, 3, 5))

			move := func(where int, who string) {
				board.Children()[where].(*eui.Button).SetText(who)
			}

			for i := range 9 {
				board.Add(eui.NewButton("", func(b *eui.Button) {
					if game.IsGameEnd() || game.GetBoard()[i] != pos.TurnEmpty {
						return
					}
					game.Move(i)
					move(i, game.GetNextTurn())
					if !game.IsGameEnd() {
						best := game.BestMove()
						game.Move(best)
						move(best, game.GetNextTurn())
					}
					if game.IsGameEnd() {
						if game.IsWinFor(pos.TurnX) {
							statusGame.Emit("Won " + string(pos.TurnX))
						} else if game.IsWinFor(pos.TurnO) {
							statusGame.Emit("Won " + string(pos.TurnO))
						} else if game.Blanks() == 0 {
							statusGame.Emit("Draw!")
						}
					} else {
						statusGame.Emit("Turn:" + game.GetTurn())
					}
					fmt.Println(game)
				}))
			}

			statusLine := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{30, 70}, 5))
			statusLine.Add(eui.NewButton("Reset game", func(b *eui.Button) {
				game.Reset()
				statusGame.Emit("Turn:" + game.GetTurn())
				for i := range board.Children() {
					move(i, string(pos.TurnEmpty))
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
