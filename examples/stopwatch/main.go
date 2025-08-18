package main

import (
	"fmt"
	"log"
	"time"

	"github.com/t0l1k/eui"
)

const title = "Stopwatch Example"

const (
	SwIdle eui.ViewState = iota + 100
	SwRun
	SwPause
)

type Btns int

const (
	SwStart Btns = iota
	SwPlay
	SwStop
	SwReset
	SwRing
)

func (s Btns) String() string {
	return []string{"Старт", "Пуск", "Стоп", "Обнулить", "Круг"}[s]
}

type Lap struct {
	index      int
	total, lap time.Duration
}

func (l Lap) String() string {
	return fmt.Sprintf("#%v Ring:%v Total:%v", l.index, l.lap.Truncate(time.Millisecond), l.total.Truncate(time.Millisecond))
}

type SwScene struct {
	*eui.Scene
	swMain, swLap *eui.Stopwatch
	swState       *eui.Signal[eui.ViewState]
	sTotal, sLap  *eui.Signal[time.Duration]
	interval      time.Duration
	last          time.Time
}

func NewSwScene() *SwScene {
	var (
		lblMainSw, lblRingSw *eui.Label
		listRing             *eui.ListView
		laps                 []Lap
	)
	s := &SwScene{Scene: eui.NewScene(eui.NewVBoxLayout(10)), interval: 70 * time.Millisecond}
	lblMainSw = eui.NewLabel("MainSW")
	lblRingSw = eui.NewLabel("RingSw")
	listRing = eui.NewListView()

	s.swMain = eui.NewStopwatch()
	s.swLap = eui.NewStopwatch()
	s.swState = eui.NewSignal(func(a, b eui.ViewState) bool { return a == b })
	s.sTotal = eui.NewSignal(func(a, b time.Duration) bool { return a == b })
	s.sLap = eui.NewSignal(func(a, b time.Duration) bool { return a == b })

	s.sTotal.ConnectAndFire(func(data time.Duration) {
		lblMainSw.SetText(eui.FormatSmartDuration(data, true))
	}, s.swMain.Duration())

	s.sLap.ConnectAndFire(func(data time.Duration) {
		lblRingSw.SetText(eui.FormatSmartDuration(data, true))
	}, s.swLap.Duration())

	swStart := func() {
		s.swMain.Reset().Start()
		s.swLap.Reset().Start()
		s.swState.Emit(SwRun)
		log.Println("SwStart")
	}

	swPlay := func() {
		s.swMain.Start()
		s.swLap.Start()
		s.swState.Emit(SwRun)
		s.last = time.Now()
		log.Println("SwPlay")
	}

	swStop := func() {
		if s.swState.Value() == SwRun {
			s.swMain.Stop()
			s.swLap.Stop()
			s.swState.Emit(SwPause)
			log.Println("SwStop")
		}
	}

	swReset := func() {
		s.swMain.Reset()
		s.swLap.Reset()
		s.swState.Emit(SwIdle)
		s.sTotal.Emit(s.swMain.Duration())
		s.sLap.Emit(s.swLap.Duration())
		laps = nil
		log.Println("SwReset")
	}

	swRing := func() Lap {
		if s.swState.Value() != SwRun {
			return Lap{}
		}
		s.swLap.Stop()
		lap := Lap{
			index: len(laps) + 1,
			total: s.swMain.Duration(),
			lap:   s.swLap.Duration(),
		}
		laps = append(laps, lap)
		s.swLap.Reset().Start()
		s.sLap.Emit(s.swLap.Duration())
		return lap
	}

	contBtns := eui.NewContainer(eui.NewHBoxLayout(5))
	for i := range 5 {
		contBtns.Add(eui.NewButton(Btns(i).String(), func(b *eui.Button) {
			if b.IsDisabled() {
				return
			}
			switch b.Text() {
			case SwStart.String():
				swStart()
				log.Println("Sw", SwStart.String())
			case SwPlay.String():
				swPlay()
				log.Println("Sw", SwPlay.String())
			case SwStop.String():
				swStop()
				log.Println("Sw", SwStop.String())
			case SwReset.String():
				listRing.Reset()
				swReset()
				log.Println("Sw", SwReset.String())
			case SwRing.String():
				if entry := swRing(); entry.lap > 0 {
					listRing.AddItem(eui.NewLabel(entry.String()))
				}
				log.Println("Sw", SwRing.String())
			}
		}))
	}

	visibleBtnsContainer := eui.NewContainerVisibleByFilter(eui.NewHBoxLayout(5), contBtns, func(d eui.Drawabler, ss eui.ViewState) bool {
		btn, ok := d.(*eui.Button)
		if !ok {
			return false
		}
		layout := func(s eui.ViewState) (result []string) {
			switch s {
			case SwIdle:
				result = append(result, SwStart.String())
			case SwRun:
				result = append(result, SwStop.String(), SwRing.String())
			case SwPause:
				result = append(result, SwPlay.String(), SwReset.String())
			}
			return result
		}
		for _, lbl := range layout(ss) {
			if btn.Text() == lbl {
				return true
			}
		}
		return false
	})

	s.swState.ConnectAndFire(func(data eui.ViewState) {
		visibleBtnsContainer.UpdateBy(data)
	}, SwIdle)

	s.Add(lblMainSw)
	s.Add(lblRingSw)
	s.Add(listRing)
	s.Add(visibleBtnsContainer)

	return s
}

func (s *SwScene) Tick(dt eui.TickData) {
	if !(s.swState.Value() == SwRun) {
		return
	}
	if time.Since(s.last) > s.interval {
		s.sTotal.Emit(s.swMain.Duration())
		s.sLap.Emit(s.swLap.Duration())
		s.last = time.Now()
	}
}

func main() {
	eui.Init(func() *eui.Ui { return eui.GetUi().SetTitle(title).SetSize(320, 200) }())
	eui.Run(NewSwScene())
}
