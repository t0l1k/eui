package game

import (
	"image/color"

	"github.com/t0l1k/eui"
)

type BallColor int

const (
	BallNoColor BallColor = iota
	BallBrown
	BallYellow
	BallGreen
	BallRed
	BallAqua
	BallBlue
	BallMagenta
	BallPurple
	BallOrange
)

func (b BallColor) String() string {
	switch b {
	case BallBrown:
		return "Brown"
	case BallYellow:
		return "Yellow"
	case BallGreen:
		return "Green"
	case BallRed:
		return "Red"
	case BallAqua:
		return "Aqua"
	case BallBlue:
		return "Blue"
	case BallMagenta:
		return "Magenta"
	case BallPurple:
		return "Purple"
	case BallOrange:
		return "Orange"
	default:
		return "None"
	}
}

func (b BallColor) Color() color.RGBA {
	switch b {
	case BallBrown:
		return eui.Maroon
	case BallYellow:
		return eui.Yellow
	case BallGreen:
		return eui.Green
	case BallRed:
		return eui.Red
	case BallAqua:
		return eui.Aqua
	case BallBlue:
		return eui.Blue
	case BallMagenta:
		return eui.Fuchsia
	case BallPurple:
		return eui.Purple
	case BallOrange:
		return eui.Orange
	default:
		return eui.Silver
	}
}
