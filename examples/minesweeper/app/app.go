package app

import ui "github.com/t0l1k/eui"

func NewGame() *ui.Ui {
	u := ui.GetUi()
	theme := ui.NewTheme()
	theme.Set("fg", ui.Yellow)
	theme.Set("bg", ui.Navy)
	loc := ui.NewLocale()
	loc.Set("lblUpTm", "Up")
	u.ApplyLocale(&loc)
	u.ApplyTheme(&theme)
	u.SetTitle("Minesweeper")
	return u
}
