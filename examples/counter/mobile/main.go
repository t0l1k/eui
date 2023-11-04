package emob

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/counter/app"
	"github.com/t0l1k/eui/examples/counter/app/scene_counter"
)

func init() {
	eui.Init(app.NewGame())
	eui.Run(scene_counter.NewSceneCounter())
	mobile.SetGame(eui.GetUi())
}

func Dummy() {}
