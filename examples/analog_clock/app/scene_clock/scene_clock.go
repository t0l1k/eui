package scene_clock

import (
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/analog_clock/app"
	"github.com/t0l1k/eui/examples/analog_clock/app/clock"
)

type SceneAnalogClock struct {
	eui.SceneBase
	topBar   *eui.TopBar
	clock    *clock.AnalogClock
	lblTm    *eui.Text
	tmVar    *eui.StringVar
	checkBox *eui.Checkbox
}

func NewSceneAnalogClock() *SceneAnalogClock {
	s := &SceneAnalogClock{}
	s.topBar = eui.NewTopBar("Analog Clock Example")
	s.Add(s.topBar)
	s.clock = clock.NewAnalogClock()
	s.Add(s.clock)
	s.lblTm = eui.NewText("")
	s.tmVar = eui.NewStringVar("")
	s.tmVar.Attach(s.lblTm)
	s.Add(s.lblTm)
	conf := eui.GetUi().GetSettings()
	s.checkBox = eui.NewCheckbox("MSecond View?", func(c *eui.Checkbox) {
		s.clock.MsHand.ToggleVisible()
		conf.Set(app.ShowMSecondHand, c.IsChecked())
	})
	s.checkBox.SetChecked(conf.Get(app.ShowMSecondHand).(bool))
	s.Add(s.checkBox)
	s.setupTheme()
	s.Resize()
	return s
}

func (s *SceneAnalogClock) setupTheme() {
	theme := eui.GetUi().GetTheme()
	s.topBar.Bg(theme.Get(app.AppBg))
	s.topBar.Fg(theme.Get(app.AppFg))
	s.clock.Bg(theme.Get(app.AppBg))
	s.clock.Fg(theme.Get(app.AppFg))
	s.clock.FaceBg = theme.Get(app.AppfaceBg)
	s.clock.FaceFg = theme.Get(app.AppfaceFg)
	s.lblTm.Bg(theme.Get(app.ApplblBg))
	s.lblTm.Fg(theme.Get(app.ApplblFg))
}

func (s *SceneAnalogClock) Update(dt int) {
	dtFormat := "2006-01-02 15:04:05"
	tm := time.Now().Format(dtFormat)
	s.tmVar.Set(tm)
	s.SceneBase.Update(dt)
}

func (s *SceneAnalogClock) Resize() {
	w0, h0 := eui.GetUi().Size()
	h := int(float64(h0) * 0.05)
	s.topBar.Resize([]int{0, 0, w0, h})
	s.clock.Resize([]int{0, h, w0, h0 - h})
	s.lblTm.Resize([]int{0, h0 - h, h * 4, h})
	s.checkBox.Resize([]int{w0 - h*4, h0 - h, h * 4, h})
}
