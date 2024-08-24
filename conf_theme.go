package eui

import (
	"image/color"

	"github.com/t0l1k/eui/colors"
)

type ThemeValue int

const (
	ButtonBg ThemeValue = iota
	ButtonFg
	ButtonHover
	ButtonFocus
	ButtonNormal
	ButtonActive
	ButtonSelected
	ButtonDisabled
	TextBg
	TextFg
	CheckboxBg
	CheckboxFg
	ComboBoxBg
	ComboBoxFg
	InputBoxBg
	InputBoxFg
	ListViewBg
	ListViewFg
	ListViewItemBg
	ListViewItemFg
	SceneBg
	SceneFg
	ViewBg
	TopBarBg
	TopBarQuitBg
	TopBarQuitFg
	TopBarTitleBg
	TopBarTitleFg
	TopBarStopwatchBg
	TopBarStopwatchFg
)

type Theme map[ThemeValue]color.Color

func DefaultTheme() *Theme {
	t := NewTheme()
	t.Set(ButtonFg, colors.Black)
	t.Set(ButtonBg, colors.Silver)
	t.Set(ButtonHover, colors.Yellow)
	t.Set(ButtonFocus, colors.Red)
	t.Set(ButtonNormal, colors.Green)
	t.Set(ButtonActive, colors.White)
	t.Set(ButtonSelected, colors.Maroon)
	t.Set(ButtonDisabled, colors.Gray)
	t.Set(TextBg, colors.Green)
	t.Set(TextFg, colors.White)
	t.Set(CheckboxBg, colors.Teal)
	t.Set(CheckboxFg, colors.Black)
	t.Set(ComboBoxBg, colors.Navy)
	t.Set(ComboBoxFg, colors.White)
	t.Set(InputBoxBg, colors.GreenYellow)
	t.Set(InputBoxFg, colors.Fuchsia)
	t.Set(ListViewBg, colors.Silver)
	t.Set(ListViewFg, colors.Orange)
	t.Set(ListViewItemBg, colors.Blue)
	t.Set(ListViewItemFg, colors.Yellow)
	t.Set(SceneBg, colors.Navy)
	t.Set(SceneFg, colors.Yellow)
	t.Set(ViewBg, colors.Navy)
	t.Set(TopBarBg, colors.Gray)
	t.Set(TopBarQuitBg, colors.Silver)
	t.Set(TopBarQuitFg, colors.Black)
	t.Set(TopBarTitleBg, colors.YellowGreen)
	t.Set(TopBarTitleFg, colors.Black)
	t.Set(TopBarStopwatchBg, colors.YellowGreen)
	t.Set(TopBarStopwatchFg, colors.Black)
	return &t
}

func NewTheme() Theme {
	return make(Theme)
}

func (t Theme) Get(set ThemeValue) (value color.Color) {
	return t[set]
}

func (t Theme) Set(set ThemeValue, value color.Color) {
	t[set] = value
}
