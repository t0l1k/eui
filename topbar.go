package eui

type TopBar struct {
	View
	btnQuit         *Button
	lblTitle, tmLbl *Text
	tmVar           *StringVar
	Stopwatch       *Stopwatch
	showStopwatch   bool
}

func NewTopBar(title string) *TopBar {
	t := &TopBar{}
	t.showStopwatch = false
	t.SetupView()
	sq := "<"
	if GetUi().IsMainScene() {
		sq = "x"
	}
	t.btnQuit = NewButton(sq, func(b *Button) {
		GetUi().Pop()
	})
	t.Add(t.btnQuit)
	t.lblTitle = NewText(title)
	t.Add(t.lblTitle)
	return t
}

func (t *TopBar) initStopwatch() {
	t.Stopwatch = NewStopwatch()
	t.tmVar = NewStringVar(t.Stopwatch.StringShort())
	t.tmLbl = NewText("0:00")
	t.tmVar.Attach(t.tmLbl)
	t.Add(t.tmLbl)
	t.Stopwatch.Start()
}

func (t *TopBar) SetShowStopwatch() {
	t.showStopwatch = !t.showStopwatch
	if t.showStopwatch {
		t.initStopwatch()
	} else {
		t.Stopwatch.Stop()
	}
}

func (t *TopBar) Update(dt int) {
	t.View.Update(dt)
	if !t.showStopwatch {
		return
	}
	t.tmVar.Set(t.Stopwatch.StringShort())
}

func (t *TopBar) Resize(arr []int) {
	t.View.Resize(arr)
	x, y, w, h := 0, 0, t.GetRect().H, t.GetRect().H
	t.btnQuit.Resize([]int{x, y, w, h})
	x += h
	w = int(float64(t.rect.W) * 0.25)
	t.lblTitle.Resize([]int{x, y, w, h})
	w = int(float64(t.rect.W) * 0.1)
	if t.showStopwatch {
		x = t.GetRect().W - w
		t.tmLbl.Resize([]int{x, y, w, h})
	}
	t.Dirty(true)
}
