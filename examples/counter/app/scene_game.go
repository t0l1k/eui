package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type SceneGame struct {
	ui.ContainerDefault
	topBar *TopBar
}

func NewSceneGame() *SceneGame {
	sc := &SceneGame{}
	sc.topBar = NewTopBar()
	sc.Add(sc.topBar)
	return sc
}

func (sc *SceneGame) Entered() {
	sc.Resize()

	sc.topBar.btnQuit.SetBg(ui.Aqua)
	sc.topBar.btnQuit.SetFg(ui.Black)
}
func (sc *SceneGame) Update(dt int) {
	for _, v := range sc.Container {
		v.Update(dt)
	}
}
func (sc *SceneGame) Draw(surface *ebiten.Image) {
	surface.Fill(ui.Red)
	for _, v := range sc.Container {
		v.Draw(surface)
	}
}
func (sc *SceneGame) Resize() {
	sc.topBar.Resize()
}

func (sc *SceneGame) Close() {
	for _, v := range sc.Container {
		v.Close()
	}
}
