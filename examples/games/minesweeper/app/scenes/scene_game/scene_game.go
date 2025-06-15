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
	eui.SceneBase
	topBar                 *eui.TopBar
	game                   *Game
	lblMines, lblGameTimer *eui.Text
	btnStatus              *eui.ButtonIcon
	btnAF                  *eui.Button
}

func NewSceneGame(title string, r, c, m int) *SceneGame {
	s := &SceneGame{}
	s.topBar = eui.NewTopBar(title, nil)
	s.Add(s.topBar)
	s.game = newGame(r, c, m)
	s.Add(s.game)
	s.game.field.State.Connect(s.UpdateData)
	s.lblMines = eui.NewText("" + strconv.Itoa(s.game.field.GetTotalMines()))
	s.Add(s.lblMines)
	s.lblGameTimer = eui.NewText("00:00")
	s.Add(s.lblGameTimer)
	s.btnStatus = eui.NewButtonIcon([]*ebiten.Image{res.SmileSprites[0], res.SmileSprites[1]}, func(b *eui.Button) {
		if b.IsMouseDownLeft() {
			s.game.New()
		} else if b.IsMouseDownRight() {
			s.game.Reset()
		} else if b.IsMouseDownMiddle() {
			s.game.field.RestoreGame()
			s.game.timer.Start()
		}
	})
	s.Add(s.btnStatus)

	s.btnAF = eui.NewButton("Auto Mark Flags", func(b *eui.Button) {
		if s.game.field.State.Value() == game.GamePlay {
			s.game.field.AutoMarkAllFlags()
		}
	})
	s.Add(s.btnAF)
	s.Resize()
	return s
}

func (s *SceneGame) Update(dt int) {
	s.checkBtnStatus()
	s.lblGameTimer.SetText(s.game.timer.StringShort())
	str := strconv.Itoa(s.game.field.GetMarkedMines()) + "/" + strconv.Itoa(s.game.field.GetTotalMines())
	s.lblMines.SetText(str)
	s.SceneBase.Update(dt)
}

func (s *SceneGame) checkBtnStatus() {
	if s.game.field.State.Value() == game.GamePlay {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			s.btnStatus.SetReleasedIcon(res.SmileSprites[0])
		} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s.btnStatus.SetReleasedIcon(res.SmileSprites[2])
		}
	}
}

func (s *SceneGame) UpdateData(value string) {
	switch value {
	case game.GameStart:
		s.btnStatus.SetReleasedIcon(res.SmileSprites[0])
	case game.GameWin:
		s.btnStatus.SetReleasedIcon(res.SmileSprites[3])
	case game.GameOver:
		s.btnStatus.SetReleasedIcon(res.SmileSprites[4])
	}
}

func (s *SceneGame) Resize() {
	w, h := eui.GetUi().Size()
	hTop := int(float64(h) * 0.05) // topbar height
	hTopHalf := int((float64(h) * 0.05) / 2)
	rect := eui.NewRect([]int{0, hTop, w, h - hTop})

	s.topBar.Resize([]int{0, 0, w, hTop})
	s.game.Resize([]int{hTopHalf, hTop + hTopHalf, w - hTop, h - hTop*2})

	x := rect.CenterX() - hTop*3
	y := 0
	r := []int{x, y, hTop * 2, hTop}
	s.lblMines.Resize(r)
	x = rect.CenterX() + hTop
	y = 0
	r = []int{x, y, hTop * 2, hTop}
	s.lblGameTimer.Resize(r)
	x = rect.CenterX() - hTop/2
	y = 0
	r = []int{x, y, hTop, hTop}
	s.btnStatus.Resize(r)
	x = rect.CenterX() + hTop*4
	r = []int{x, y, hTop * 3, hTop}
	s.btnAF.Resize(r)
}
