package eui

import "image/color"

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

var (
	Black  = color.RGBA{0, 0, 0, 255}
	Gray   = color.RGBA{128, 128, 128, 255}
	Silver = color.RGBA{192, 192, 192, 255}
	White  = color.RGBA{255, 255, 255, 255}

	Orange  = color.RGBA{255, 165, 0, 255}
	Fuchsia = color.RGBA{255, 0, 255, 255}
	Purple  = color.RGBA{128, 0, 128, 255}
	Red     = color.RGBA{255, 0, 0, 255}
	Maroon  = color.RGBA{128, 0, 0, 255}

	Yellow      = color.RGBA{255, 255, 0, 255}
	GreenYellow = color.RGBA{173, 255, 47, 255}
	YellowGreen = color.RGBA{154, 205, 50, 255}
	Olive       = color.RGBA{128, 128, 0, 255}
	Lime        = color.RGBA{0, 255, 0, 255}
	Green       = color.RGBA{0, 128, 0, 255}

	Aqua = color.RGBA{0, 255, 255, 255}
	Teal = color.RGBA{0, 128, 128, 255}
	Blue = color.RGBA{0, 0, 255, 255}
	Navy = color.RGBA{0, 0, 128, 255}
)

type Theme map[ThemeValue]color.Color

func DefaultTheme() *Theme {
	t := NewTheme()
	t.Set(ButtonFg, Black)
	t.Set(ButtonBg, Silver)
	t.Set(ButtonHover, Yellow)
	t.Set(ButtonFocus, Red)
	t.Set(ButtonNormal, Green)
	t.Set(ButtonActive, White)
	t.Set(ButtonSelected, Maroon)
	t.Set(ButtonDisabled, Gray)
	t.Set(TextBg, Green)
	t.Set(TextFg, White)
	t.Set(CheckboxBg, Teal)
	t.Set(CheckboxFg, Black)
	t.Set(ComboBoxBg, Navy)
	t.Set(ComboBoxFg, White)
	t.Set(InputBoxBg, GreenYellow)
	t.Set(InputBoxFg, Fuchsia)
	t.Set(ListViewBg, Silver)
	t.Set(ListViewFg, Orange)
	t.Set(ListViewItemBg, Blue)
	t.Set(ListViewItemFg, Yellow)
	t.Set(SceneBg, Navy)
	t.Set(SceneFg, Yellow)
	t.Set(ViewBg, Navy)
	t.Set(TopBarBg, Gray)
	t.Set(TopBarQuitBg, Silver)
	t.Set(TopBarQuitFg, Black)
	t.Set(TopBarTitleBg, YellowGreen)
	t.Set(TopBarTitleFg, Black)
	t.Set(TopBarStopwatchBg, YellowGreen)
	t.Set(TopBarStopwatchFg, Black)
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
