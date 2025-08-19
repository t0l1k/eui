package eui

import (
	"image/color"

	"golang.org/x/image/colornames"
)

// ViewState перечисление состояний
type ViewState int

const (
	StateNormal ViewState = iota
	StateFocused
	StateHover
	StateSelected
	StateDisabled
	StateHidden
)

func (s ViewState) IsHovered() bool  { return s == StateHover }
func (s ViewState) IsBlurred() bool  { return s == StateNormal }
func (s ViewState) IsFocused() bool  { return s == StateFocused }
func (s ViewState) IsSelected() bool { return s == StateSelected }
func (s ViewState) IsDisabled() bool { return s == StateDisabled }
func (s ViewState) IsHidden() bool   { return s == StateHidden }

func (e ViewState) Color() color.Color {
	return [...]color.Color{
		colornames.Gray,
		colornames.Red,
		colornames.White,
		colornames.Green,
		colornames.Darkgray,
		color.Transparent,
	}[e]
}

func (s ViewState) String() string {
	return [...]string{
		"Normal",
		"Focused",
		"Hover",
		"Selected",
		"Disabled",
		"Hidden",
	}[s]
}

type ViewType int

const (
	ViewNormal ViewType = iota
	ViewModal
	ViewSystem
	ViewBackground
)

func (s ViewType) IsNormal() bool     { return s == ViewNormal }
func (s ViewType) IsModal() bool      { return s == ViewModal }
func (s ViewType) IsSystem() bool     { return s == ViewSystem }
func (s ViewType) IsBackground() bool { return s == ViewBackground }
