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
	return []string{
		"None",
		"Brown",
		"Yellow",
		"Green",
		"Red",
		"Aqua",
		"Blue",
		"Magenta",
		"Purple",
		"Orange",
	}[b]
}

func (b BallColor) Color() color.Color {
	return []color.Color{
		colornames.Silver,
		colornames.Maroon,
		colornames.Yellow,
		colornames.Green,
		colornames.Red,
		colornames.Aqua,
		colornames.Blue,
		colornames.Fuchsia,
		colornames.Purple,
		colornames.Orange,
	}[b]
}
