package clock

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/analog_clock/app"
)

type AnalogClock struct {
	eui.View
	MsHand, secHand, minuteHand, hourHand *Hand
	FaceBg, FaceFg                        color.Color
}

func NewAnalogClock() *AnalogClock {
	theme := eui.GetUi().GetTheme()
	a := &AnalogClock{
		MsHand:     NewHand(theme.Get(app.AppMsSecondHandFg)),
		secHand:    NewHand(theme.Get(app.AppSecondHandFg)),
		minuteHand: NewHand(theme.Get(app.AppMinuteHandFg)),
		hourHand:   NewHand(theme.Get(app.AppHourHandFg)),
	}
	a.SetupAnalogClock()
	return a
}

func (a *AnalogClock) SetupAnalogClock() {
	a.SetupView()
	a.Add(a.hourHand)
	a.Add(a.minuteHand)
	a.Add(a.secHand)
	a.Add(a.MsHand)
}

func (a *AnalogClock) Layout() {
	a.View.Layout()
	a.drawClockFace()
	log.Println("update analog clock layout done")
	a.Dirty(false)
}

func (a *AnalogClock) drawClockFace() {
	x, y := a.GetRect().Center()
	m := float64(a.GetRect().GetLowestSize()) * 0.01
	vector.DrawFilledCircle(a.GetImage(), float32(x), float32(y), float32(m)*3, a.GetBg(), true)
	center := eui.Point{X: float64(x), Y: float64(y)}
	vector.DrawFilledCircle(a.GetImage(), float32(center.X), float32(center.Y), float32(m)*2, a.FaceBg, true)
	vector.DrawFilledCircle(a.GetImage(), float32(center.X), float32(center.Y), float32(m), a.FaceFg, true)
	for i := 0; i < 60; i++ {
		var (
			tip eui.Point
			rad float64
		)
		if i%5 == 0 {
			rad = m * 2.0
		} else {
			rad = m
		}
		sz := center.Y
		if center.Y > center.X {
			sz = center.X
		}
		tip = GetTip(center, float64(i)/60, sz-m*4, 0, 0)
		vector.DrawFilledCircle(a.GetImage(), float32(tip.X), float32(tip.Y), float32(rad), a.FaceBg, true)
		vector.DrawFilledCircle(a.GetImage(), float32(tip.X), float32(tip.Y), float32(rad)/2, a.FaceFg, true)
	}
}

func (g *AnalogClock) getTime() (msec, sec, min, hour int) {
	dt := time.Now()
	msec = dt.Nanosecond() / 1e6
	sec = dt.Second()
	min = dt.Minute()
	hour = dt.Hour()
	return
}

func (g *AnalogClock) Update(dt int) {
	g.View.Update(dt)
	msec, sec, minute, hour := g.getTime()
	g.MsHand.Set(float64(msec) / 1000.0)
	g.secHand.Set((float64(sec) + g.MsHand.Get()) / 60.0)
	g.minuteHand.Set((float64(minute) + g.secHand.Get()) / 60.0)
	g.hourHand.Set((float64(hour) + g.minuteHand.Get()) / 12.0)
}

func (t *AnalogClock) Draw(surface *ebiten.Image) {
	if !t.IsVisible() {
		return
	}
	if t.IsDirty() {
		t.Layout()
		for _, c := range t.Container {
			c.Layout()
		}
	}
	op := &ebiten.DrawImageOptions{}
	x, y := t.GetRect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(t.GetImage(), op)
	for _, v := range t.Container {
		v.Draw(surface)
	}
}

func (a *AnalogClock) Resize(r []int) {
	a.Rect(r)
	a.MsHand.Resize(r)
	a.secHand.Resize(r)
	a.minuteHand.Resize(r)
	a.hourHand.Resize(r)
	a.setupHands()
	a.Dirty(true)
	a.Image(nil)
}

func (a *AnalogClock) setupHands() {
	sz := a.GetRect().GetLowestSize()
	m := (float64(sz) * 0.01)
	x, y := a.GetRect().Center()
	center := *eui.NewPoint(float64(x), float64(y))
	lenght := float64(sz/2) - m*4

	conf := eui.GetUi().GetSettings()
	a.MsHand.Setup(center, lenght, 1, conf.Get(app.ShowMSecondHand).(bool))
	a.secHand.Setup(center, lenght, 3, true)
	lenght = float64(sz/2) - m*8
	a.minuteHand.Setup(center, lenght, 5, true)
	lenght = float64(sz/2) - m*12
	a.hourHand.Setup(center, lenght, 8, true)
}
