package app

import "github.com/t0l1k/eui"

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Analog Clock")
	k := 2
	w, h := 320*k, 200*k
	u.SetSize(w, h)
	u.GetSettings().Set(eui.UiFullscreen, false)
	u.GetSettings().Set(ShowMSecondHand, false)
	setAppTheme()
	return u
}

const (
	AppBg eui.ThemeValue = iota + 100
	AppFg
	ApplblBg
	ApplblFg
	AppfaceBg
	AppfaceFg
	AppMsSecondHandFg
	AppSecondHandFg
	AppMinuteHandFg
	AppHourHandFg

	ShowMSecondHand eui.SettingName = iota + 100
)

func setAppTheme() {
	theme := eui.GetUi().GetTheme()
	theme.Set(AppBg, eui.Silver)
	theme.Set(AppFg, eui.Black)
	theme.Set(ApplblBg, eui.GreenYellow)
	theme.Set(ApplblFg, eui.Black)
	theme.Set(AppfaceBg, eui.Navy)
	theme.Set(AppfaceFg, eui.GreenYellow)
	theme.Set(AppMsSecondHandFg, eui.Black)
	theme.Set(AppSecondHandFg, eui.Red)
	theme.Set(AppMinuteHandFg, eui.Blue)
	theme.Set(AppHourHandFg, eui.Navy)
}
