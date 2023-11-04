package clock

import (
	"math"

	"github.com/t0l1k/eui"
)

func GetTip(center eui.Point, percent, lenght, width, height float64) (tip eui.Point) {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	sine := math.Sin(radians)
	cosine := math.Cos(radians)
	tip.X = center.X + lenght*sine - width
	tip.Y = center.Y + lenght*cosine - height
	return tip
}

func GetAngle(percent float64) float64 {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	angle := (radians * -180 / math.Pi)
	return angle
}
