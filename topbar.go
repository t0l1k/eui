package eui

type TopBar struct {
	View
	btnQuit                    *Button
	lblTitle, tmLbl            *Text
	tmVar                      *StringVar
	Stopwatch                  *Stopwatch
	showStopwatch              bool
	coverTitle, coverStopwatch float64
}

func NewTopBar(title string) *TopBar {
	t := &TopBar{coverTitle: 0.25, coverStopwatch: 0.1}
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
	t.setTheme()
	return t
}

func (t *TopBar) SetTitleCoverArea(value float64) {
	t.coverTitle = value
	t.Dirty(true)
}

func (t *TopBar) SetStopwatchCoverArea(value float64) {
	t.coverStopwatch = value
	t.Dirty(true)
}

func (t *TopBar) setTheme() {
	theme := GetUi().GetTheme()
	t.bg = theme.Get(TopBarBg)
	t.btnQuit.bg = theme.Get(TopBarQuitBg)
	t.btnQuit.fg = theme.Get(TopBarQuitFg)
	t.lblTitle.bg = theme.Get(TopBarTitleBg)
	t.lblTitle.fg = theme.Get(TopBarTitleFg)
}

func (t *TopBar) initStopwatch() {
	t.Stopwatch = NewStopwatch()
	t.tmVar = NewStringVar(t.Stopwatch.StringShort())
	t.tmLbl = NewText("0:00")
	theme := GetUi().GetTheme()
	t.tmLbl.bg = theme.Get(TopBarStopwatchBg)
	t.tmLbl.fg = theme.Get(TopBarStopwatchFg)
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
	t.tmVar.SetValue(t.Stopwatch.StringShort())
}

func (t *TopBar) Resize(arr []int) {
	t.View.Resize(arr)
	x, y, w, h := 0, 0, t.GetRect().H, t.GetRect().H
	t.btnQuit.Resize([]int{x, y, w, h})
	x += h
	w = int(float64(t.rect.W) * t.coverTitle)
	t.lblTitle.Resize([]int{x, y, w, h})
	if t.showStopwatch {
		w = int(float64(t.rect.W) * t.coverStopwatch)
		x = t.GetRect().W - w
		t.tmLbl.Resize([]int{x, y, w, h})
	}
	t.Dirty(true)
}
