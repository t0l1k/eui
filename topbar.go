package eui

type TopBar struct {
	*Container
	btnMenu                        *Button
	btnFunc                        func(b *Button)
	lblTitle, tmLbl                *Text
	tmVar                          *Signal[string]
	Stopwatch                      *Stopwatch
	useSW, showSW, showTitle, show bool
	coverTitle, coverStopwatch     float64
}

// Умею показать вверху строку с меткой текста с кнопкой выход из сцены(если nil) или переопределенной функцией(для вызова диалога, например) в параметре и секундомером нахождения на сцене
func NewTopBar(title string, f func(b *Button)) *TopBar {
	t := &TopBar{Container: NewContainer(NewAbsoluteLayout()), coverTitle: 0.25, coverStopwatch: 0.1}
	t.show = true
	t.showTitle = true
	t.useSW = false
	t.btnFunc = f
	btnText := "Menu"
	if f == nil {
		btnText = "<"
		t.btnFunc = func(b *Button) {
			GetUi().Pop()
		}
		if GetUi().IsMainScene() {
			btnText = "x"
		}
	}
	t.btnMenu = NewButton(btnText, t.btnFunc)
	t.Add(t.btnMenu)
	t.lblTitle = NewText(title)
	t.Add(t.lblTitle)
	t.setTheme()
	return t
}

func (t *TopBar) SetTitle(text string) {
	t.lblTitle.SetText(text)
}

func (t *TopBar) SetButtonText(text string) {
	t.btnMenu.SetText(text)
}

func (t *TopBar) SetButtonFunc(f func(b *Button)) {
	t.btnFunc = f
}

func (t *TopBar) SetTitleCoverArea(value float64) {
	t.coverTitle = value
	t.MarkDirty()
}

func (t *TopBar) SetStopwatchCoverArea(value float64) {
	t.coverStopwatch = value
	t.MarkDirty()
}

func (t *TopBar) setTheme() {
	theme := GetUi().GetTheme()
	t.bg = theme.Get(TopBarBg)
	t.btnMenu.bg = theme.Get(TopBarQuitBg)
	t.btnMenu.fg = theme.Get(TopBarQuitFg)
	t.lblTitle.bg = theme.Get(TopBarTitleBg)
	t.lblTitle.fg = theme.Get(TopBarTitleFg)
}

func (t *TopBar) initStopwatch() {
	t.Stopwatch = NewStopwatch()
	t.tmVar = NewSignal(func(a, b string) bool { return a == b })
	t.tmLbl = NewText("0:00")
	theme := GetUi().GetTheme()
	t.tmLbl.bg = theme.Get(TopBarStopwatchBg)
	t.tmLbl.fg = theme.Get(TopBarStopwatchFg)
	t.tmVar.Connect(func(data string) { t.tmLbl.SetText(data) })
	t.Add(t.tmLbl)
	t.Stopwatch.Start()
}

func (t *TopBar) SetShowStoppwatch(value bool) {
	t.showSW = value
	if t.showSW {
		t.tmLbl.SetHidden(false)
	} else {
		t.tmLbl.SetHidden(true)
	}
}

func (t *TopBar) SetUseStopwatch() {
	t.useSW = !t.useSW
	if t.useSW {
		t.initStopwatch()
	} else {
		t.Stopwatch.Stop()
	}
}

func (t *TopBar) SetShowTitle(value bool) {
	t.showTitle = value
	if t.showTitle {
		t.lblTitle.SetHidden(false)
	} else {
		t.lblTitle.SetHidden(true)
	}
}

func (t *TopBar) IsVisible() bool { return t.show }
func (t *TopBar) Visible(value bool) {
	t.show = value
	if t.showTitle {
		t.lblTitle.SetHidden(value)
	}
	if t.showSW {
		t.tmLbl.SetHidden(value)
	}
	if !value {
		t.btnMenu.Disable()
	} else {
		t.btnMenu.Enable()
	}
	t.btnMenu.SetHidden(value)
}

func (t *TopBar) Update(dt int) {
	t.Container.Update(dt)
	if !t.useSW {
		return
	}
	t.tmVar.Emit(t.Stopwatch.StringShort())
}

func (t *TopBar) Resize(rect Rect[int]) {
	t.SetRect(rect)
	x, y, w, h := 0, 0, t.Rect().H, t.Rect().H
	t.btnMenu.Resize(NewRect([]int{x, y, w, h}))
	x += h
	w = int(float64(rect.W) * t.coverTitle)
	t.lblTitle.Resize(NewRect([]int{x, y, w, h}))
	if t.useSW {
		w = int(float64(rect.W) * t.coverStopwatch)
		x = t.Rect().W - w
		t.tmLbl.Resize(NewRect([]int{x, y, w, h}))
	}
	t.ImageReset()
}
