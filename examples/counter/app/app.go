package app

import ui "github.com/t0l1k/eui"

func NewGame() *ui.Ui {
	u := ui.GetUi()
	theme := ui.NewTheme()
	theme.Set("bg", ui.Navy)
	u.ApplyTheme(&theme)
	loc := ui.NewLocale()
	loc.Set("lblUpTm", "Up")
	u.ApplyLocale(&loc)
	u.SetTitle("Counter")
	return u
}
