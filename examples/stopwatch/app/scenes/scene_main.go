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
	*eui.Scene
	btnsCont       *eui.Container
	list           *eui.ListView
	swMain, swRing *eui.Stopwatch
	var0, var1     *eui.Signal[string]
	sBtns          []string
	state          watchState
	count          int
	_dirty         bool
}

func NewSceneStopwatch() *SceneStopwatch {
	s := &SceneStopwatch{Scene: eui.NewScene(eui.NewVBoxLayout(1))}
	s.swMain = eui.NewStopwatch()
	s.swRing = eui.NewStopwatch()

	s.var0 = eui.NewSignal(func(a, b string) bool { return a == b })
	s.var1 = eui.NewSignal(func(a, b string) bool { return a == b })

	lblTimeMain := eui.NewLabel("Нажми старт")
	s.var0.Connect(func(data string) { lblTimeMain.SetText(data) })

	lblTimeSecond := eui.NewLabel("0.0")
	s.var1.Connect(func(data string) { lblTimeSecond.SetText(data) })
	s.list = eui.NewListView()

	s.btnsCont = eui.NewContainer(eui.NewHBoxLayout(1))
	s.sBtns = []string{"Обнулить", "Старт", "Круг"}
	for _, value := range s.sBtns {
		button := eui.NewButton(value, s.stopwatchAppLogic)
		s.btnsCont.Add(button)
	}
	s.state = watchStart
	s._dirty = true

	s.Add(lblTimeMain)
	s.Add(lblTimeSecond)
	s.Add(s.list)
	s.Add(s.btnsCont)
	return s
}

func (s *SceneStopwatch) stopwatchAppLogic(b *eui.Button) {
	switch b.Text() {
	case s.sBtns[0]:
		switch s.state {
		case watchPause, watchPlay:
			s.swMain.Reset()
			s.swRing.Reset()
			s.state = watchStart
			s.count = 0
			log.Println("set state start from reset")
			s._dirty = true
		}
	case s.sBtns[1], "Пауза":
		switch s.state {
		case watchStart:
			s.swMain.Start()
			s.swRing.Start()
			b.SetText("Пауза")
			s.state = watchPlay
			log.Println("set state play from start")
			s.MarkDirty()
		case watchPlay:
			s.swMain.Stop()
			s.swRing.Stop()
			b.SetText(s.sBtns[1])
			s.state = watchPause
			log.Println("set state pause from play")
			s.MarkDirty()
		case watchPause:
			s.swMain.Start()
			s.swRing.Start()
			b.SetText("Пауза")
			s.state = watchPlay
			log.Println("set state from pause play")
			s._dirty = true
		}
	case s.sBtns[2]:
		switch s.state {
		case watchPlay:
			s.count++
			str := strconv.Itoa(s.count) + " R:" + s.swRing.String() + " M:" + s.swMain.String()
			s.swRing.Reset()
			txt := eui.NewLabel(str)
			s.list.AddItem(txt)
			s.swRing.Start()
			s._dirty = true
		}
	}
}

func (s *SceneStopwatch) Update(dt int) {
	s.var0.Emit(s.swMain.String())
	s.var1.Emit(s.swRing.String())

	if s._dirty {
		switch s.state {
		case watchStart:
			s.btnsCont.Children()[1].(*eui.Button).SetText(s.sBtns[1])
			s.btnsCont.Children()[0].(*eui.Button).Hide()
			s.btnsCont.Children()[2].(*eui.Button).Hide()
			s.list.Reset()
		case watchPlay:
			s.btnsCont.Children()[0].(*eui.Button).Hide()
			s.btnsCont.Children()[2].(*eui.Button).Show()
		case watchPause:
			s.btnsCont.Children()[0].(*eui.Button).Show()
			s.btnsCont.Children()[2].(*eui.Button).Hide()
		}
		s._dirty = false
	}

}
