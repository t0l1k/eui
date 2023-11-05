package clock

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
)

type AnalogClock struct {
	eui.View
	msHand, secHand, minuteHand, hourHand *Hand
	msHandVisible                         bool
}

func NewAnalogClock() *AnalogClock {
	a := &AnalogClock{
		msHand:        NewHand(),
		secHand:       NewHand(),
		minuteHand:    NewHand(),
		hourHand:      NewHand(),
		msHandVisible: true,
	}
	a.SetupAnalogClock()
	return a
}

func (a *AnalogClock) SetupAnalogClock() {
	a.SetupView()
	a.Add(a.hourHand)
	a.Add(a.minuteHand)
	a.Add(a.secHand)
	a.Add(a.msHand)
}

func (a *AnalogClock) Layout() {
	a.View.Layout()
	// a.drawBack()
	a.drawClockFace()
	// a.drawDebugLines()
	log.Println("update analog clock layout done")
}

func (a *AnalogClock) drawClockFace() {
	x, y := a.GetRect().Center()
	m := float64(a.GetRect().GetLowestSize()) * 0.01
	bg := eui.Aqua
	fg := eui.Black
	vector.DrawFilledCircle(a.GetImage(), float32(x), float32(y), float32(m)*3, bg, true)
	center := eui.Point{X: float64(x), Y: float64(y)}
	vector.DrawFilledCircle(a.GetImage(), float32(center.X), float32(center.Y), float32(m)*2, bg, true)
	vector.DrawFilledCircle(a.GetImage(), float32(center.X), float32(center.Y), float32(m), fg, true)
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
		tip = GetTip(center, float64(i)/60, center.Y-m*4, 0, 0)
		vector.DrawFilledCircle(a.GetImage(), float32(tip.X), float32(tip.Y), float32(rad), bg, true)
		vector.DrawFilledCircle(a.GetImage(), float32(tip.X), float32(tip.Y), float32(rad)/2, fg, true)
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
	g.msHand.Set(float64(msec) / 1000.0)
	g.secHand.Set((float64(sec) + g.msHand.Get()) / 60.0)
	g.minuteHand.Set((float64(minute) + g.secHand.Get()) / 60.0)
	g.hourHand.Set((float64(hour) + g.minuteHand.Get()) / 12.0)
	if g.msHandVisible {
		g.msHand.SetTip()
	}
	g.secHand.SetTip()
	g.minuteHand.SetTip()
	g.hourHand.SetTip()
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
	a.msHand.Resize(r)
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
	a.msHand.Setup(center, lenght, eui.Yellow, 1)
	a.secHand.Setup(center, lenght, eui.Red, 3)
	lenght = float64(sz/2) - m*8
	a.minuteHand.Setup(center, lenght, eui.Green, 5)
	lenght = float64(sz/2) - m*12
	a.hourHand.Setup(center, lenght, eui.Blue, 8)
}

// func (a *AnalogClock) drawDebugLines() {
// 	_, _, w, h := a.GetRect().GetRectFloat()
// 	vector.StrokeLine(a.GetImage(), 0, 0, w, h, 2, eui.Red, true)
// 	vector.StrokeLine(a.GetImage(), 0, h, w, 0, 2, eui.Red, true)
// 	vector.StrokeLine(a.GetImage(), 0, h/2, w, h/2, 2, eui.Red, true)
// 	vector.StrokeLine(a.GetImage(), w/2, h, w/2, 0, 2, eui.Red, true)
// }

// func (a *AnalogClock) drawBack() {
// 	m := float64(a.GetRect().GetLowestSize()) * 0.01
// 	_, _, w0, h0 := a.GetRect().GetRectFloat()
// 	bg := eui.Aqua
// 	fg := eui.Black
// 	a.GetImage().Fill(bg)
// 	vector.StrokeRect(a.GetImage(), 0, 0, w0, h0, 3, fg, true)
// 	for y := 0; y < a.GetRect().H; y += int(m) {
// 		for x := 0; x < a.GetRect().W; x += int(m) {
// 			vector.StrokeRect(a.GetImage(), float32(x), 0, float32(x), h0, 1, fg, true)
// 			vector.StrokeRect(a.GetImage(), 0, float32(y), w0, float32(y), 1, fg, true)
// 		}
// 	}
// }
