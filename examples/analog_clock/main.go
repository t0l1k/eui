// Пример аналоговые часы с плавным ходом секундной и миллисекундной стрелки, ещё есть строка состояния вверху, там выход со сцены и приложения, ещё внизу-влево есть метка с датой и временем, и кнопкой включить и выключить показ миллискундной стрелки внизу-справа.
// Ещё в примере показано как настроить внешний вид сцены и приложения

package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

// Стрелка часов с рисованием прозрачного фона и самой стрелки, можно напрямую рисовать через draw, только стрелку, но так элегантнее, все через макет и уже в контейнере обновить и перерисовать
type Hand struct {
	*eui.Drawable
	faceCenter, tip          eui.Point
	value, lenght, thickness float64
}

func NewHand(bg, fg color.Color) *Hand {
	h := &Hand{Drawable: eui.NewDrawable()}
	r, g, b, _ := bg.RGBA()
	a := 0
	col := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	h.Bg(col)
	h.Fg(fg)
	return h
}

func (h *Hand) Setup(center eui.Point, lenght float64, thickness float64, visible bool) {
	h.faceCenter = center
	h.lenght = lenght
	h.thickness = thickness
	h.Visible(visible)
}

func (h *Hand) ToggleVisible() {
	h.Visible(!h.IsVisible())
	h.MarkDirty()
}

func (h *Hand) Get() float64 { return h.value }

func (h *Hand) Set(value float64) {
	if h.value == value {
		return
	}
	h.value = value
	h.tip = GetTip(h.faceCenter, h.value, h.lenght, 0, 0)
	h.MarkDirty()
}

func (h *Hand) Layout() {
	h.Drawable.Layout() // подготовить холст
	x1 := h.Rect().CenterX()
	y1 := h.Rect().CenterY()
	x2 := int(h.tip.X)
	y2 := int(h.tip.Y)
	vector.StrokeLine(h.Image(), float32(x1), float32(y1), float32(x2), float32(y2), float32(h.thickness), h.GetFg(), true)
	h.ClearDirty()
}

func (h *Hand) Draw(surface *ebiten.Image) {
	if !h.IsVisible() || h.IsDisabled() {
		return
	}
	if h.IsDirty() {
		h.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := h.Rect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(h.Image(), op)
}

func (t *Hand) Resize(rect eui.Rect) { t.SetRect(rect); t.MarkDirty() }

// Сами часы, рисуется "лицо часов", а поверх него стрелки в порядке добавления друг поверх друга.
type AnalogClock struct {
	*eui.Container
	MsHand, secHand, minuteHand, hourHand *Hand
	FaceBg, FaceFg                        color.Color
}

func NewAnalogClock() *AnalogClock {
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(AppfaceBg)
	a := &AnalogClock{
		Container:  eui.NewContainer(eui.NewAbsoluteLayout()),
		MsHand:     NewHand(bg, theme.Get(AppMsSecondHandFg)),
		secHand:    NewHand(bg, theme.Get(AppSecondHandFg)),
		minuteHand: NewHand(bg, theme.Get(AppMinuteHandFg)),
		hourHand:   NewHand(bg, theme.Get(AppHourHandFg)),
	}
	a.Add(a.hourHand)
	a.Add(a.minuteHand)
	a.Add(a.secHand)
	a.Add(a.MsHand)
	a.Visible(true)
	return a
}

func (a *AnalogClock) Layout() {
	a.Drawable.Layout()
	a.Image().Fill(a.GetBg())
	a.drawClockFace()
	log.Println("update analog clock layout done")
	a.ClearDirty()
}

func (a *AnalogClock) drawClockFace() {
	x, y := a.Rect().Center()
	m := float64(a.Rect().GetLowestSize()) * 0.01
	vector.DrawFilledCircle(a.Image(), float32(x), float32(y), float32(m)*3, a.GetBg(), true)
	center := eui.Point{X: float64(x), Y: float64(y)}
	vector.DrawFilledCircle(a.Image(), float32(center.X), float32(center.Y), float32(m)*2, a.FaceBg, true)
	vector.DrawFilledCircle(a.Image(), float32(center.X), float32(center.Y), float32(m), a.FaceFg, true)
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
		vector.DrawFilledCircle(a.Image(), float32(tip.X), float32(tip.Y), float32(rad), a.FaceBg, true)
		vector.DrawFilledCircle(a.Image(), float32(tip.X), float32(tip.Y), float32(rad)/2, a.FaceFg, true)
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
		for _, c := range t.Childrens() {
			c.Layout()
		}
	}
	op := &ebiten.DrawImageOptions{}
	x, y := t.Rect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(t.Image(), op)
	for _, v := range t.Childrens() {
		v.Draw(surface)
	}
}

func (a *AnalogClock) Resize(r eui.Rect) {
	a.SetRect(r)
	a.MsHand.Resize(r)
	a.secHand.Resize(r)
	a.minuteHand.Resize(r)
	a.hourHand.Resize(r)
	a.setupHands()
	a.ImageReset()
}

func (a *AnalogClock) setupHands() {
	sz := a.Rect().GetLowestSize()
	m := (float64(sz) * 0.01)
	x, y := a.Rect().Center()
	center := eui.NewPoint(float64(x), float64(y))
	lenght := float64(sz/2) - m*4

	conf := eui.GetUi().GetSettings()
	a.MsHand.Setup(center, lenght, 1, conf.Get(ShowMSecondHand).(bool))
	a.secHand.Setup(center, lenght, 3, true)
	lenght = float64(sz/2) - m*8
	a.minuteHand.Setup(center, lenght, 5, true)
	lenght = float64(sz/2) - m*12
	a.hourHand.Setup(center, lenght, 8, true)
}

func GetTip(center eui.Point, percent, lenght, width, height float64) (tip eui.Point) {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	sine := math.Sin(radians)
	cosine := math.Cos(radians)
	tip.X = center.X + lenght*sine - width
	tip.Y = center.Y + lenght*cosine - height
	return tip
}

func GetAngle(percent float64) float64 {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	angle := (radians * -180 / math.Pi)
	return angle
}

type SceneAnalogClock struct {
	*eui.Scene
	topBar   *eui.TopBar
	clock    *AnalogClock
	lblTm    *eui.Text
	tmVar    *eui.Signal[string]
	checkBox *eui.Checkbox
}

func NewSceneAnalogClock() *SceneAnalogClock {
	s := &SceneAnalogClock{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar("Analog Clock Example", nil)
	s.Add(s.topBar)
	s.clock = NewAnalogClock()
	s.Add(s.clock)
	s.lblTm = eui.NewText("")
	s.Add(s.lblTm)
	s.tmVar = eui.NewSignal(func(a, b string) bool { return a == b })
	s.tmVar.Connect(func(data string) {
		s.lblTm.SetText(data)
	})
	conf := eui.GetUi().GetSettings()
	s.checkBox = eui.NewCheckbox("MSecond View?", func(c *eui.Checkbox) {
		s.clock.MsHand.ToggleVisible()
		conf.Set(ShowMSecondHand, c.IsChecked())
	})
	s.Add(s.checkBox)
	s.checkBox.SetChecked(conf.Get(ShowMSecondHand).(bool))
	s.setupTheme()
	s.Resize()
	return s
}

func (s *SceneAnalogClock) setupTheme() {
	theme := eui.GetUi().GetTheme()
	s.topBar.Bg(theme.Get(AppBg))
	s.topBar.Fg(theme.Get(AppFg))
	s.clock.Bg(theme.Get(AppBg))
	s.clock.Fg(theme.Get(AppFg))
	s.clock.FaceBg = theme.Get(AppfaceBg)
	s.clock.FaceFg = theme.Get(AppfaceFg)
	s.lblTm.Bg(theme.Get(ApplblBg))
	s.lblTm.Fg(theme.Get(ApplblFg))
}

func (s *SceneAnalogClock) Update(dt int) {
	dtFormat := "2006-01-02 15:04:05"
	tm := time.Now().Format(dtFormat)
	s.tmVar.Emit(tm)
	s.Scene.Update(dt)
}

func (s *SceneAnalogClock) Resize() {
	w0, h0 := eui.GetUi().Size()
	h := int(float64(h0) * 0.05)
	s.topBar.Resize(eui.NewRect([]int{0, 0, w0, h}))
	s.clock.Resize(eui.NewRect([]int{0, h, w0, h0 - h}))
	s.lblTm.Resize(eui.NewRect([]int{0, h0 - h, h * 4, h}))
	s.checkBox.Resize(eui.NewRect([]int{w0 - h*4, h0 - h, h * 4, h}))
}

func NewGame() *eui.Ui {
	u := eui.GetUi().SetTitle("Analog Clock").SetSize(800, 600)
	u.GetSettings().Set(eui.UiFullscreen, false)
	u.GetSettings().Set(ShowMSecondHand, false)
	setAppTheme()
	return u
}

const (
	AppBg eui.ThemeValue = iota + 100
	AppFg
	ApplblBg
	ApplblFg
	AppfaceBg
	AppfaceFg
	AppMsSecondHandFg
	AppSecondHandFg
	AppMinuteHandFg
	AppHourHandFg

	ShowMSecondHand eui.SettingName = iota + 100
)

func setAppTheme() {
	theme := eui.GetUi().GetTheme()
	theme.Set(AppBg, colornames.Silver)
	theme.Set(AppFg, colornames.Black)
	theme.Set(ApplblBg, colornames.Greenyellow)
	theme.Set(ApplblFg, colornames.Black)
	theme.Set(AppfaceBg, colornames.Navy)
	theme.Set(AppfaceFg, colornames.Greenyellow)
	theme.Set(AppMsSecondHandFg, colornames.Black)
	theme.Set(AppSecondHandFg, colornames.Red)
	theme.Set(AppMinuteHandFg, colornames.Blue)
	theme.Set(AppHourHandFg, colornames.Navy)
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneAnalogClock())
	eui.Quit(func() {})
}
