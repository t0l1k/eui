package eui

type Topbar struct {
	*Container
	btnMenu                        *Button
	btnFunc                        func(b *Button)
	lblTitle, tmLbl                *Label
	tmVar                          *Signal[string]
	Stopwatch                      *Stopwatch
	useSW, showSW, showTitle, show bool
}

// Умею показать вверху строку с меткой текста с кнопкой выход из сцены(если nil) или переопределенной функцией(для вызова диалога, например) в параметре и секундомером нахождения на сцене
func NewTopBar(title string, fn func(b *Button)) *Topbar {
	t := &Topbar{Container: NewContainer(NewLayoutHorizontalPercent([]int{5, 30, 55, 10}, 1))}
	t.show = true
	t.showTitle = true
	t.useSW = false
	t.btnFunc = fn
	btnText := "Menu"
	if fn == nil {
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
	t.lblTitle = NewLabel(title)
	t.Add(t.lblTitle)
	t.Add(NewDrawable())
	t.tmLbl = NewLabel("")
	t.Add(t.tmLbl)
	t.setTheme()
	return t
}

func (t *Topbar) SetTitle(text string)            { t.lblTitle.SetText(text) }
func (t *Topbar) SetButtonText(text string)       { t.btnMenu.SetText(text) }
func (t *Topbar) SetButtonFunc(f func(b *Button)) { t.btnFunc = f }

func (t *Topbar) setTheme() {
	theme := GetUi().Theme()
	t.bg = theme.Get(TopBarBg)
	t.btnMenu.bg = theme.Get(TopBarQuitBg)
	t.btnMenu.fg = theme.Get(TopBarQuitFg)
	t.lblTitle.bg = theme.Get(TopBarTitleBg)
	t.lblTitle.fg = theme.Get(TopBarTitleFg)
}

func (t *Topbar) initStopwatch() {
	t.Stopwatch = NewStopwatch()
	t.tmVar = NewSignal(func(a, b string) bool { return a == b })
	t.tmLbl.SetText("00:00")
	theme := GetUi().Theme()
	t.tmLbl.bg = theme.Get(TopBarStopwatchBg)
	t.tmLbl.fg = theme.Get(TopBarStopwatchFg)
	t.tmVar.Connect(func(data string) { t.tmLbl.SetText(data) })
	t.Stopwatch.Start()
}

func (t *Topbar) SetShowStoppwatch(value bool) *Topbar {
	t.showSW = value
	if t.showSW {
		t.tmLbl.Show()
	} else {
		t.tmLbl.Hide()
	}
	return t
}

func (t *Topbar) SetUseStopwatch() *Topbar {
	t.useSW = !t.useSW
	if t.useSW {
		t.initStopwatch()
	} else {
		t.Stopwatch.Stop()
	}
	return t
}

func (t *Topbar) SetShowTitle(value bool) *Topbar {
	t.showTitle = value
	if t.showTitle {
		t.lblTitle.Show()
	} else {
		t.lblTitle.Hide()
	}
	return t
}

func (t *Topbar) Tick(TickData) {
	if !t.useSW {
		return
	}
	t.tmVar.Emit(FormatSmartDuration(t.Stopwatch.Duration(), false))
}
