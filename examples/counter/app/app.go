package app

import ui "github.com/t0l1k/eui"

func NewGame() *ui.Ui {
	u := ui.GetUi()
	theme := ui.NewTheme()
	theme.Set("bg", ui.Navy)
	u.ApplyTheme(&theme)
	u.SetTitle("Counter")
	return u
}
