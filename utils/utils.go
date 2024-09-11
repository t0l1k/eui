package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
)

func DrawDebugLines(surface *ebiten.Image, rect *eui.Rect) {
	x0, y0 := rect.Pos()
	w0, h0 := rect.BottomRight()
	x, y := float32(x0), float32(y0)
	w, h := float32(w0), float32(h0)
	vector.StrokeLine(surface, x, y, w, h, 2, colors.Red, true)
	vector.StrokeLine(surface, x, h, w, y, 1, colors.Red, true)
	vector.StrokeLine(surface, x, h/2, w, h/2, 1, colors.Red, true)
	vector.StrokeLine(surface, w/2, h, w/2, y, 1, colors.Red, true)
}
