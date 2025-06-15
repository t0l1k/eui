package game

import (
	"image/color"

	"golang.org/x/image/colornames"
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
		return colornames.Maroon
	case BallYellow:
		return colornames.Yellow
	case BallGreen:
		return colornames.Green
	case BallRed:
		return colornames.Red
	case BallAqua:
		return colornames.Aqua
	case BallBlue:
		return colornames.Blue
	case BallMagenta:
		return colornames.Fuchsia
	case BallPurple:
		return colornames.Purple
	case BallOrange:
		return colornames.Orange
	default:
		return colornames.Silver
	}
}
