package game

import (
	"image/color"

	"github.com/t0l1k/eui/colors"
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
		return colors.Maroon
	case BallYellow:
		return colors.Yellow
	case BallGreen:
		return colors.Green
	case BallRed:
		return colors.Red
	case BallAqua:
		return colors.Aqua
	case BallBlue:
		return colors.Blue
	case BallMagenta:
		return colors.Fuchsia
	case BallPurple:
		return colors.Purple
	case BallOrange:
		return colors.Orange
	default:
		return colors.Silver
	}
}
