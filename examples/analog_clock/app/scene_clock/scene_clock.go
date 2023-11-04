package scene_clock

import (
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/analog_clock/app/clock"
)

type SceneAnalogClock struct {
	eui.SceneDefault
	topBar *TopBar
	clock  *clock.AnalogClock
	lblTm  *eui.Text
	tmVar  *eui.StringVar
}

func NewSceneAnalogClock() *SceneAnalogClock {
	s := &SceneAnalogClock{}
	s.topBar = NewTopBar("Analog Clock Example")
	s.Add(s.topBar)
	s.clock = clock.NewAnalogClock()
	s.Add(s.clock)
	bg := eui.Blue
	fg := eui.Yellow
	s.lblTm = eui.NewText("", bg, fg)
	s.tmVar = eui.NewStringVar("")
	s.tmVar.Attach(s.lblTm)
	s.Add(s.lblTm)
	s.Resize()
	return s
}

func (s *SceneAnalogClock) Update(dt int) {
	dtFormat := "2006-01-02 15:04:05"
	tm := time.Now().Format(dtFormat)
	s.tmVar.Set(tm)
	s.SceneDefault.Update(dt)
}

func (s *SceneAnalogClock) Resize() {
	w0, h0 := eui.GetUi().Size()
	h := int(float64(h0) * 0.05)
	s.topBar.Resize([]int{0, 0, w0, h})
	s.clock.Resize([]int{0, h, w0, h0 - h})
	s.lblTm.Resize([]int{0, h0 - h, h * 4, h})
}
