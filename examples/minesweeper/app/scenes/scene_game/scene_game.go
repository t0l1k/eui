package scene_game

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/app/game"
	"github.com/t0l1k/eui/examples/minesweeper/app/res"
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
	s.topBar = eui.NewTopBar(title)
	s.Add(s.topBar)
	s.game = newGame(r, c, m)
	s.Add(s.game)
	s.lblMines = eui.NewText("" + strconv.Itoa(s.game.field.GetTotalMines()))
	s.Add(s.lblMines)
	s.lblGameTimer = eui.NewText("00:00")
	s.Add(s.lblGameTimer)
	s.btnStatus = eui.NewButtonIcon([]*ebiten.Image{res.SmileSprites[0], res.SmileSprites[1]}, func(b *eui.ButtonIcon) {
		if b.IsMouseDownLeft() {
			s.game.New()
		} else if b.IsMouseDownRight() {
			s.game.Reset()
		} else if b.IsMouseDownMiddle() {
			s.game.field.RestoreGame()
			s.game.redraw()
			s.game.timer.Start()
		}
	})
	s.Add(s.btnStatus)

	s.btnAF = eui.NewButton("Auto Mark Flags", func(b *eui.Button) {
		switch s.game.field.GetState() {
		case game.GamePlay:
			s.game.field.AutoMarkAllFlags()
			s.game.redraw()
		}
	})
	s.Add(s.btnAF)
	s.Resize()
	log.Println("Game init done")
	return s
}

func (s *SceneGame) Update(dt int) {
	switch s.game.field.GetState() {
	case game.GameStart, game.GamePlay:
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			s.btnStatus.SetReleasedIcon(res.SmileSprites[0])
			log.Println("smile0")
		} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s.btnStatus.SetReleasedIcon(res.SmileSprites[2])
			log.Println("smile2")
		}
	case game.GameWin:
		s.btnStatus.SetReleasedIcon(res.SmileSprites[3])
	case game.GameOver:
		s.btnStatus.SetReleasedIcon(res.SmileSprites[4])
	}
	s.lblGameTimer.SetText(s.game.timer.StringShort())
	str := strconv.Itoa(s.game.field.GetLeftMines()) + "/" + strconv.Itoa(s.game.field.GetTotalMines())
	s.lblMines.SetText(str)
	s.SceneBase.Update(dt)
}

func (s *SceneGame) Resize() {
	w, h := eui.GetUi().Size()
	hTop := int(float64(h) * 0.05) // topbar height
	rect := eui.NewRect([]int{0, hTop, w, h - hTop})

	s.topBar.Resize([]int{0, 0, w, hTop})
	s.game.Resize([]int{0, hTop, w, h - hTop})

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