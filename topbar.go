package eui

type TopBar struct {
	View
	btnMenu                    *Button
	btnFunc                    func(b *Button)
	lblTitle, tmLbl            *Text
	tmVar                      *SubjectBase
	Stopwatch                  *Stopwatch
	showStopwatch              bool
	coverTitle, coverStopwatch float64
}

// Умею показать вверху строку с меткой текста с кнопкой выход из сцены(если nil) или переопределенной функцией(для вызова диалога, например) в параметре и секундомером нахождения на сцене
func NewTopBar(title string, f func(b *Button)) *TopBar {
	t := &TopBar{coverTitle: 0.25, coverStopwatch: 0.1}
	t.showStopwatch = false
	t.SetupView()
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

func (t *TopBar) SetButtonText(text string) {
	t.btnMenu.SetText(text)
}

func (t *TopBar) SetButtonFunc(f func(b *Button)) {
	t.btnFunc = f
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
	t.btnMenu.bg = theme.Get(TopBarQuitBg)
	t.btnMenu.fg = theme.Get(TopBarQuitFg)
	t.lblTitle.bg = theme.Get(TopBarTitleBg)
	t.lblTitle.fg = theme.Get(TopBarTitleFg)
}

func (t *TopBar) initStopwatch() {
	t.Stopwatch = NewStopwatch()
	t.tmVar = NewSubject()
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
	t.btnMenu.Resize([]int{x, y, w, h})
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
