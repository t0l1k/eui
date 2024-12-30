package eui

type TopBar struct {
	DrawableBase
	btnMenu                        *Button
	btnFunc                        func(b *Button)
	lblTitle, tmLbl                *Text
	tmVar                          *SubjectBase
	Stopwatch                      *Stopwatch
	useSW, showSW, showTitle, show bool
	coverTitle, coverStopwatch     float64
}

// Умею показать вверху строку с меткой текста с кнопкой выход из сцены(если nil) или переопределенной функцией(для вызова диалога, например) в параметре и секундомером нахождения на сцене
func NewTopBar(title string, f func(b *Button)) *TopBar {
	t := &TopBar{coverTitle: 0.25, coverStopwatch: 0.1}
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
	t.Dirty = true
}

func (t *TopBar) SetStopwatchCoverArea(value float64) {
	t.coverStopwatch = value
	t.Dirty = true
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

func (t *TopBar) SetShowStoppwatch(value bool) {
	t.showSW = value
	if t.showSW {
		t.tmLbl.Visible(true)
	} else {
		t.tmLbl.Visible(false)
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
		t.lblTitle.Visible(true)
	} else {
		t.lblTitle.Visible(false)
	}
}

func (t *TopBar) IsVisible() bool { return t.show }
func (t *TopBar) Visible(value bool) {
	t.show = value
	if t.showTitle {
		t.lblTitle.Visible(value)
	}
	if t.showSW {
		t.tmLbl.Visible(value)
	}
	if !value {
		t.btnMenu.Disable()
	} else {
		t.btnMenu.Enable()
	}
	t.btnMenu.Visible(value)
}

func (t *TopBar) Update(dt int) {
	t.DrawableBase.Update(dt)
	if !t.useSW {
		return
	}
	t.tmVar.SetValue(t.Stopwatch.StringShort())
}

func (t *TopBar) Resize(rect []int) {
	t.Rect(NewRect(rect))
	t.SpriteBase.Resize(rect)
	x, y, w, h := 0, 0, t.GetRect().H, t.GetRect().H
	t.btnMenu.Resize([]int{x, y, w, h})
	x += h
	w = int(float64(t.rect.W) * t.coverTitle)
	t.lblTitle.Resize([]int{x, y, w, h})
	if t.useSW {
		w = int(float64(t.rect.W) * t.coverStopwatch)
		x = t.GetRect().W - w
		t.tmLbl.Resize([]int{x, y, w, h})
	}
	t.ImageReset()
}
