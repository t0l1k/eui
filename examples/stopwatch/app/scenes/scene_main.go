package scenes

import (
	"log"

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
	btnsCont *eui.Container
	// list           *eui.ListView
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

	lblTimeMain := eui.NewText("Нажми старт")
	lblTimeMain.OnlyOneFontSize(true)
	s.var0.Connect(func(data string) { lblTimeMain.SetText(data) })

	lblTimeSecond := eui.NewText("0.0")
	lblTimeSecond.OnlyOneFontSize(true)
	s.var1.Connect(func(data string) { lblTimeSecond.SetText(data) })
	// s.list = eui.NewListView()

	s.btnsCont = eui.NewContainer(eui.NewHBoxLayout(1))
	s.sBtns = []string{"Обнулить", "Старт", "Круг"}
	for _, value := range s.sBtns {
		button := eui.NewButton(value, func(b *eui.Button) {
			log.Println(b.GetText())
		})
		s.btnsCont.Add(button)
	}
	s.state = watchStart
	s._dirty = true

	s.Add(lblTimeMain)
	s.Add(lblTimeSecond)
	// s.Add(s.list)
	s.Add(s.btnsCont)
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
			// str := strconv.Itoa(s.count) + " R:" + s.swRing.String() + " M:" + s.swMain.String()
			s.swRing.Reset()
			// txt := eui.NewText(str)
			// s.list.Add(txt)
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
			s.btnsCont.Childrens()[1].(*eui.Button).SetText(s.sBtns[1])
			s.btnsCont.Childrens()[0].(*eui.Button).Visible(false)
			s.btnsCont.Childrens()[2].(*eui.Button).Visible(false)
			// s.list.Reset()
		case watchPlay:
			s.btnsCont.Childrens()[0].(*eui.Button).Visible(false)
			s.btnsCont.Childrens()[2].(*eui.Button).Visible(true)
		case watchPause:
			s.btnsCont.Childrens()[0].(*eui.Button).Visible(true)
			s.btnsCont.Childrens()[2].(*eui.Button).Visible(false)
		}
		s._dirty = false
	}
	// s.SceneBase.Update(dt)

	// for _, v := range s.frame0.GetContainer() {
	// 	v.Update(dt)
	// }
	// for _, v := range s.frame1.GetContainer() {
	// 	v.Update(dt)
	// }
}

// func (s *SceneStopwatch) Draw(surface *ebiten.Image) {
// 	s.SceneBase.Draw(surface)
// 	for _, v := range s.frame0.GetContainer() {
// 		v.Draw(surface)
// 	}
// 	for _, v := range s.frame1.GetContainer() {
// 		v.Draw(surface)
// 	}
// }

// func (s *SceneStopwatch) Resize() {
// 	w0, h0 := eui.GetUi().Size()
// 	x, y := 0, 0
// 	rect := eui.NewRect([]int{x, y, w0, h0})
// 	hTop := int(float64(rect.GetLowestSize()) * 0.05)
// 	s.topBar.Resize([]int{x, y, w0, hTop})

// 	y += hTop
// 	h1 := (h0 - hTop) / 5
// 	s.frame0.Resize([]int{x, y, w0, h1 * 2})
// 	y += h1 * 2
// 	s.list.Resize([]int{x, y, w0, h1 * 2})
// 	s.list.Itemsize(hTop)
// 	y += h1 * 2
// 	s.frame1.Resize([]int{x, y, w0, h1})
// }
