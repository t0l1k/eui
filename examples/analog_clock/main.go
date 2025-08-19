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
	faceCenter, tip          eui.Point[float64]
	value, lenght, thickness float64
	show                     bool
}

func NewHand(bg, fg color.Color) *Hand {
	h := &Hand{Drawable: eui.NewDrawable()}
	h.SetBg(color.Transparent)
	h.SetFg(fg)
	return h
}

func (h *Hand) Setup(center eui.Point[float64], lenght float64, thickness float64, visible bool) {
	h.faceCenter = center
	h.lenght = lenght
	h.thickness = thickness
	h.show = !visible
	h.ToggleVisible()
}

func (h *Hand) ToggleVisible() {
	h.show = !h.show
	if h.show {
		h.Show()
	} else {
		h.Hide()
	}
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
	vector.StrokeLine(h.Image(), float32(x1), float32(y1), float32(x2), float32(y2), float32(h.thickness), h.Fg(), true)
	h.ClearDirty()
}

func (h *Hand) Draw(surface *ebiten.Image) {
	if h.IsHidden() || h.IsDisabled() || !h.show {
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

// Сами часы, рисуется "лицо часов", а поверх него стрелки в порядке добавления друг поверх друга.
type AnalogClock struct {
	*eui.Container
	MsHand, secHand, minuteHand, hourHand *Hand
	FaceBg, FaceFg                        color.Color
}

func NewAnalogClock() *AnalogClock {
	theme := eui.GetUi().Theme()
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
	return a
}

func (a *AnalogClock) Layout() {
	a.Drawable.Layout()
	a.Image().Fill(a.Bg())
	a.drawClockFace()
	log.Println("update analog clock layout done")
	a.ClearDirty()
}

func (a *AnalogClock) drawClockFace() {
	x, y := a.Rect().Center()
	m := float64(a.Rect().GetLowestSize()) * 0.01
	vector.DrawFilledCircle(a.Image(), float32(x), float32(y), float32(m)*3, a.Bg(), true)
	center := eui.Point[float64]{X: float64(x), Y: float64(y)}
	vector.DrawFilledCircle(a.Image(), float32(center.X), float32(center.Y), float32(m)*2, a.FaceBg, true)
	vector.DrawFilledCircle(a.Image(), float32(center.X), float32(center.Y), float32(m), a.FaceFg, true)
	for i := 0; i < 60; i++ {
		var (
			tip eui.Point[float64]
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

func (g *AnalogClock) Update() {
	msec, sec, minute, hour := g.getTime()
	g.MsHand.Set(float64(msec) / 1000.0)
	g.secHand.Set((float64(sec) + g.MsHand.Get()) / 60.0)
	g.minuteHand.Set((float64(minute) + g.secHand.Get()) / 60.0)
	g.hourHand.Set((float64(hour) + g.minuteHand.Get()) / 12.0)
}

func (t *AnalogClock) Draw(surface *ebiten.Image) {
	if t.IsHidden() {
		return
	}
	if t.IsDirty() {
		t.Layout()
		for _, c := range t.Children() {
			c.Layout()
		}
	}
	op := &ebiten.DrawImageOptions{}
	x, y := t.Rect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(t.Image(), op)
	for _, v := range t.Children() {
		v.Draw(surface)
	}
}

func (a *AnalogClock) SetRect(r eui.Rect[int]) {
	a.Drawable.SetRect(r)
	a.MsHand.SetRect(r)
	a.secHand.SetRect(r)
	a.minuteHand.SetRect(r)
	a.hourHand.SetRect(r)
	a.setupHands()
	a.ImageReset()
}

func (a *AnalogClock) setupHands() {
	sz := a.Rect().GetLowestSize()
	m := (float64(sz) * 0.01)
	x, y := a.Rect().Center()
	center := eui.NewPoint(float64(x), float64(y))
	lenght := float64(sz/2) - m*4

	conf := eui.GetUi().Settings()
	a.MsHand.Setup(center, lenght, 1, conf.Get(ShowMSecondHand).(bool))
	a.secHand.Setup(center, lenght, 3, true)
	lenght = float64(sz/2) - m*8
	a.minuteHand.Setup(center, lenght, 5, true)
	lenght = float64(sz/2) - m*12
	a.hourHand.Setup(center, lenght, 8, true)
}

func GetTip(center eui.Point[float64], percent, lenght, width, height float64) (tip eui.Point[float64]) {
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
	topBar   *eui.Topbar
	clock    *AnalogClock
	lblTm    *eui.Label
	tmVar    *eui.Signal[string]
	checkBox *eui.Checkbox
}

func NewSceneAnalogClock() *SceneAnalogClock {
	s := &SceneAnalogClock{Scene: eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 90, 5}, 5))}
	s.topBar = eui.NewTopBar("Analog Clock Example", nil).SetUseStopwatch()
	s.Add(s.topBar)
	s.clock = NewAnalogClock()
	s.Add(s.clock)
	contStatus := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{30, 40, 30}, 5))
	s.Add(contStatus)
	s.lblTm = eui.NewLabel("")
	contStatus.Add(s.lblTm)
	contStatus.Add(eui.NewLabel(""))
	s.tmVar = eui.NewSignal(func(a, b string) bool { return a == b })
	s.tmVar.Connect(func(data string) {
		s.lblTm.SetText(data)
	})
	conf := eui.GetUi().Settings()
	s.checkBox = eui.NewCheckbox("MSecond View?", func(c *eui.Checkbox) {
		s.clock.MsHand.ToggleVisible()
		conf.Set(ShowMSecondHand, c.IsChecked())
	})
	contStatus.Add(s.checkBox)
	s.checkBox.SetChecked(conf.Get(ShowMSecondHand).(bool))
	s.Add(eui.NewGridBackground(50))
	s.setupTheme()
	return s
}

func (s *SceneAnalogClock) setupTheme() {
	theme := eui.GetUi().Theme()
	s.topBar.SetBg(theme.Get(AppBg))
	s.topBar.SetFg(theme.Get(AppFg))
	s.clock.SetBg(theme.Get(AppBg))
	s.clock.SetFg(theme.Get(AppFg))
	s.clock.FaceBg = theme.Get(AppfaceBg)
	s.clock.FaceFg = theme.Get(AppfaceFg)
	s.lblTm.SetBg(theme.Get(ApplblBg))
	s.lblTm.SetFg(theme.Get(ApplblFg))
}

func (s *SceneAnalogClock) Update() {
	dtFormat := "2006-01-02 15:04:05"
	tm := time.Now().Format(dtFormat)
	s.tmVar.Emit(tm)
	s.Scene.Update()
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
	theme := eui.GetUi().Theme()
	theme.Set(AppBg, color.Transparent)
	theme.Set(AppFg, colornames.Black)
	theme.Set(ApplblBg, colornames.Greenyellow)
	theme.Set(ApplblFg, colornames.Black)
	theme.Set(AppfaceBg, colornames.Navy)
	theme.Set(AppfaceFg, colornames.Greenyellow)
	theme.Set(AppMsSecondHandFg, colornames.Yellow)
	theme.Set(AppSecondHandFg, colornames.Red)
	theme.Set(AppMinuteHandFg, colornames.Magenta)
	theme.Set(AppHourHandFg, colornames.Aqua)
}

func main() {
	eui.Init(func() *eui.Ui {
		u := eui.GetUi().SetTitle("Analog Clock").SetSize(800, 600)
		u.Settings().Set(eui.UiFullscreen, false)
		u.Settings().Set(ShowMSecondHand, false)
		setAppTheme()
		return u
	}())
	eui.Run(NewSceneAnalogClock())
	eui.Quit(func() {})
}
