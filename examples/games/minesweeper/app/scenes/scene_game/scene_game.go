package scene_game

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/minesweeper/app/game"
	"github.com/t0l1k/eui/examples/games/minesweeper/app/res"
)

type SceneGame struct {
	*eui.Scene
	game                   *Game
	lblMines, lblGameTimer *eui.Label
	btnStatus              *eui.Button
	btnAF                  *eui.Button
}

func NewSceneGame(title string, r, c, m int) *SceneGame {
	s := &SceneGame{Scene: eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 90, 5}, 5))}
	s.game = newGame(r, c, m)
	s.game.field.State.Connect(s.UpdateData)
	s.lblMines = eui.NewLabel("" + strconv.Itoa(s.game.field.GetTotalMines()))
	s.lblGameTimer = eui.NewLabel("00:00")
	s.btnStatus = eui.NewButtonIcon([]*ebiten.Image{res.SmileSprites[1], res.SmileSprites[0]}, func(b *eui.Button) {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			s.game.New()
		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			s.game.Reset()
		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
			s.game.field.RestoreGame()
			s.game.timer.Start()
		}
	})
	s.btnStatus.SetState(eui.StateNormal)

	s.btnAF = eui.NewButton("Auto Mark Flags", func(b *eui.Button) {
		if s.game.field.State.Value() == game.GamePlay {
			s.game.field.AutoMarkAllFlags()
		}
	})

	contStatus := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{40, 5, 5, 10, 10, 40}, 1))
	contStatus.Add(eui.NewDrawable())
	contStatus.Add(s.lblMines)
	contStatus.Add(s.btnStatus)
	contStatus.Add(s.btnAF)
	contStatus.Add(s.lblGameTimer)
	contStatus.Add(eui.NewDrawable())

	s.Add(eui.NewTopBar(title, nil))
	s.Add(s.game)
	s.Add(contStatus)
	return s
}

func (s *SceneGame) Update(dt int) {
	s.checkBtnStatus()
	s.lblGameTimer.SetText(eui.FormatSmartDuration(s.game.timer.Duration(), false))
	str := strconv.Itoa(s.game.field.GetMarkedMines()) + "/" + strconv.Itoa(s.game.field.GetTotalMines())
	s.lblMines.SetText(str)
	s.Scene.Update(dt)
}

func (s *SceneGame) checkBtnStatus() {
	if s.game.field.State.Value() == game.GamePlay {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			s.btnStatus.SetUpIcon(res.SmileSprites[0])
		} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s.btnStatus.SetUpIcon(res.SmileSprites[2])
		}
	}
}

func (s *SceneGame) UpdateData(value string) {
	switch value {
	case game.GameStart:
		s.btnStatus.SetUpIcon(res.SmileSprites[0])
	case game.GameWin:
		s.btnStatus.SetUpIcon(res.SmileSprites[3])
	case game.GameOver:
		s.btnStatus.SetUpIcon(res.SmileSprites[4])
	}
}
