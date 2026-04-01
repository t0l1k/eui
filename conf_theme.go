package eui

import (
	"image/color"

	"golang.org/x/image/colornames"
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
	t.Set(ButtonFg, colornames.Black)
	t.Set(ButtonBg, colornames.Silver)
	t.Set(ButtonHover, colornames.Yellow)
	t.Set(ButtonFocus, colornames.Red)
	t.Set(ButtonNormal, colornames.Green)
	t.Set(ButtonActive, colornames.White)
	t.Set(ButtonSelected, colornames.Maroon)
	t.Set(ButtonDisabled, colornames.Gray)
	t.Set(TextBg, colornames.Green)
	t.Set(TextFg, colornames.White)
	t.Set(CheckboxBg, colornames.Teal)
	t.Set(CheckboxFg, colornames.Yellow)
	t.Set(ComboBoxBg, colornames.Navy)
	t.Set(ComboBoxFg, colornames.White)
	t.Set(InputBoxBg, colornames.Greenyellow)
	t.Set(InputBoxFg, colornames.Fuchsia)
	t.Set(ListViewBg, colornames.Silver)
	t.Set(ListViewFg, colornames.Orange)
	t.Set(ListViewItemBg, colornames.Blue)
	t.Set(ListViewItemFg, colornames.Yellow)
	t.Set(SceneBg, colornames.Navy)
	t.Set(SceneFg, colornames.Yellow)
	t.Set(ViewBg, colornames.Navy)
	t.Set(TopBarBg, colornames.Gray)
	t.Set(TopBarQuitBg, colornames.Silver)
	t.Set(TopBarQuitFg, colornames.Black)
	t.Set(TopBarTitleBg, colornames.Yellowgreen)
	t.Set(TopBarTitleFg, colornames.Black)
	t.Set(TopBarStopwatchBg, colornames.Yellowgreen)
	t.Set(TopBarStopwatchFg, colornames.Black)
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
