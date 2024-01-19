package scenes

import (
	"log"
	"strconv"

	"github.com/t0l1k/eui"
)

type watchState int

const (
	watchStart watchState = iota
	watchPlay
	watchPause
)

type SceneStopwatch struct {
	eui.SceneBase
	topBar         *eui.TopBar
	frame0, frame1 *eui.BoxLayout
	list           *eui.ListView
	swMain, swRing *eui.Stopwatch
	var0, var1     *eui.StringVar
	sBtns          []string
	state          watchState
	count          int
	dirty          bool
}

func NewSceneStopwatch() *SceneStopwatch {
	s := &SceneStopwatch{}
	s.swMain = eui.NewStopwatch()
	s.swRing = eui.NewStopwatch()

	s.var0 = eui.NewStringVar(s.swMain.String())
	s.var1 = eui.NewStringVar(s.swRing.String())

	s.topBar = eui.NewTopBar("Секундомер", nil)

	s.frame0 = eui.NewVLayout()
	lblTimeMain := eui.NewText("Нажми старт")
	lblTimeMain.OnlyOneFontSize(true)
	s.frame0.Add(lblTimeMain)
	s.var0.Attach(lblTimeMain)

	lblTimeSecond := eui.NewText("0.0")
	lblTimeSecond.OnlyOneFontSize(true)
	s.frame0.Add(lblTimeSecond)
	s.var1.Attach(lblTimeSecond)

	s.list = eui.NewListView()

	s.frame1 = eui.NewHLayout()
	s.sBtns = []string{"Обнулить", "Старт", "Круг"}
	for _, value := range s.sBtns {
		button := eui.NewButton(value, s.stopwatchAppLogic)
		s.frame1.Add(button)
	}
	s.state = watchStart
	s.Add(s.topBar)
	s.Add(s.frame0)
	s.Add(s.list)
	s.Add(s.frame1)
	s.dirty = true
	s.Resize()
	return s
}

func (s *SceneStopwatch) stopwatchAppLogic(b *eui.Button) {
	switch b.GetText() {
	case s.sBtns[0]:
		switch s.state {
		case watchPause, watchPlay:
			s.swMain.Reset()
			s.swRing.Reset()
			s.state = watchStart
			s.count = 0
			log.Println("set state start from reset")
			s.dirty = true
		}
	case s.sBtns[1], "Пауза":
		switch s.state {
		case watchStart:
			s.swMain.Start()
			s.swRing.Start()
			b.SetText("Пауза")
			s.state = watchPlay
			log.Println("set state play from start")
			s.dirty = true
		case watchPlay:
			s.swMain.Stop()
			s.swRing.Stop()
			b.SetText(s.sBtns[1])
			s.state = watchPause
			log.Println("set state pause from play")
			s.dirty = true
		case watchPause:
			s.swMain.Start()
			s.swRing.Start()
			b.SetText("Пауза")
			s.state = watchPlay
			log.Println("set state from pause play")
			s.dirty = true
		}
	case s.sBtns[2]:
		switch s.state {
		case watchPlay:
			s.count++
			str := strconv.Itoa(s.count) + " R:" + s.swRing.String() + " M:" + s.swMain.String()
			s.swRing.Reset()
			txt := eui.NewText(str)
			s.list.Add(txt)
			s.swRing.Start()
			s.dirty = true
		}
	}
}

func (s *SceneStopwatch) Update(dt int) {
	s.var0.SetValue(s.swMain.String())
	s.var1.SetValue(s.swRing.String())

	s.SceneBase.Update(dt)
	if s.dirty {
		switch s.state {
		case watchStart:
			s.frame1.Container[1].(*eui.Button).SetText(s.sBtns[1])
			s.frame1.Container[0].(*eui.Button).Visible(false)
			s.frame1.Container[2].(*eui.Button).Visible(false)
			s.list.Reset()
		case watchPlay:
			s.frame1.Container[0].(*eui.Button).Visible(false)
			s.frame1.Container[2].(*eui.Button).Visible(true)
		case watchPause:
			s.frame1.Container[0].(*eui.Button).Visible(true)
			s.frame1.Container[2].(*eui.Button).Visible(false)
		}
		s.dirty = false
	}
}

func (s *SceneStopwatch) Resize() {
	w0, h0 := eui.GetUi().Size()
	x, y := 0, 0
	rect := eui.NewRect([]int{x, y, w0, h0})
	hTop := int(float64(rect.GetLowestSize()) * 0.05)
	s.topBar.Resize([]int{x, y, w0, hTop})

	y += hTop
	h1 := (h0 - hTop) / 5
	s.frame0.Resize([]int{x, y, w0, h1 * 2})
	y += h1 * 2
	s.list.Resize([]int{x, y, w0, h1 * 2})
	s.list.Itemsize(hTop)
	y += h1 * 2
	s.frame1.Resize([]int{x, y, w0, h1})
}
