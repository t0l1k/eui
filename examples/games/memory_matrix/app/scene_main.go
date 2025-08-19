package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/memory_matrix/mem"
	"golang.org/x/image/colornames"
)

type SceneMain struct{ *eui.Scene }

func NewSceneMain() *SceneMain {
	s := &SceneMain{Scene: eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 90, 5}, 5))}
	topBar := eui.NewTopBar(title, nil)
	topBar.SetUseStopwatch()
	topBar.SetShowStoppwatch(true)
	lblStatus := eui.NewLabel("")
	var board *BoardMem
	board = NewBoardMem(func(btn *eui.Button) {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
			switch board.Game().Stage() {
			case mem.Preparation:
				board.game.Move(0)
				board.showTimer.On()
				board.SetupShow()
			case mem.Restart:
				board.game.NextLevel()
				board.SetupPreparation()
			}
		}
		for i, v := range board.Children() {
			switch vv := v.(type) {
			case *eui.Button:
				if vv == btn {
					switch board.Game().Stage() {
					case mem.Recollection:
						if board.game.Move(i) {
							v.(*eui.Button).SetBg(colornames.Aqua)
						} else {
							v.(*eui.Button).SetBg(colornames.Orange)
						}
					}
				}
			}
		}
	})
	board.varMsg.Connect(func(data string) { lblStatus.SetText(data) })
	board.varColor.Connect(func(arr []color.Color) {
		lblStatus.SetBg(arr[0])
		lblStatus.SetFg(arr[1])
	})
	lblStatus.SetText(board.Game().String())
	s.Add(topBar)
	s.Add(board)
	s.Add(lblStatus)
	return s
}
