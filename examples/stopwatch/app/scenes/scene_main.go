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
	eui.SceneDefault
	topBar                   *TopBar
	frame0, frame1           *eui.BoxLayout
	list                     *eui.ListView
	stopperMain, stopperRing *eui.Stopwatch
	var0                     *eui.StringVar
	sBtns                    []string
	state                    watchState
	count                    int
	dirty                    bool
}

func NewSceneStopwatch() *SceneStopwatch {
	s := &SceneStopwatch{}
	s.stopperMain = eui.NewStopwatch()
	s.stopperRing = eui.NewStopwatch()

	s.var0 = eui.NewStringVar(s.stopperMain.String())

	s.topBar = NewTopBar("Стоппер")

	bg := eui.Green
	fg := eui.Yellow

	s.frame0 = eui.NewVLayout()
	lblTime := eui.NewText("Нажми старт", bg, fg)
	lblTime.OnlyOneFontSize(true)
	s.frame0.Add(lblTime)
	s.var0.Attach(lblTime)

	s.list = eui.NewListView(nil, 30)

	s.frame1 = eui.NewHLayout()
	s.sBtns = []string{"Обнулить", "Старт", "Круг"}
	for _, value := range s.sBtns {
		bg := eui.Gray
		fg := eui.Maroon
		button := eui.NewButton(value, bg, fg, s.stopwatchAppLogic)
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
			s.stopperMain.Reset()
			s.state = watchStart
			s.count = 0
			log.Println("set state start from reset")
			s.dirty = true
		}
	case s.sBtns[1], "Пауза":
		switch s.state {
		case watchStart:
			s.stopperMain.Start()
			s.stopperRing.Start()
			b.SetText("Пауза")
			s.state = watchPlay
			log.Println("set state play from start")
			s.dirty = true
		case watchPlay:
			s.stopperMain.Stop()
			b.SetText(s.sBtns[1])
			s.state = watchPause
			log.Println("set state pause from play")
			s.dirty = true
		case watchPause:
			s.stopperMain.Start()
			b.SetText("Пауза")
			s.state = watchPlay
			log.Println("set state from pause play")
			s.dirty = true
		}
	case s.sBtns[2]:
		switch s.state {
		case watchPlay:
			s.count++
			str := strconv.Itoa(s.count) + " R:" + s.stopperRing.String() + " M:" + s.stopperMain.String()
			s.stopperRing.Reset()
			txt := eui.NewText(str, eui.Blue, eui.Yellow)
			s.list.Add(txt)
			s.Resize()
			s.stopperRing.Start()
			log.Println("set new ring")
			s.dirty = true
		}
	}
}

func (s *SceneStopwatch) Update(dt int) {
	s.var0.Set(s.stopperMain.String())

	s.SceneDefault.Update(dt)
	if s.dirty {
		switch s.state {
		case watchStart:
			s.frame1.Container[1].(*eui.Button).SetText(s.sBtns[1])
			s.frame1.Container[0].(*eui.Button).Visible(false)
			s.frame1.Container[2].(*eui.Button).Visible(false)
			s.list.Container = nil
		case watchPlay:
			s.frame1.Container[0].(*eui.Button).Visible(true)
			s.frame1.Container[2].(*eui.Button).Visible(true)
		case watchPause:
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
	y += h1 * 2
	s.frame1.Resize([]int{x, y, w0, h1})
}
