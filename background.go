package eui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type GridBackground struct {
	*Drawable
	spacing int
}

func NewGridBackground(spacing int) *GridBackground {
	g := &GridBackground{Drawable: NewDrawable(), spacing: spacing}
	g.SetViewType(ViewBackground)
	return g
}

func (g *GridBackground) Layout() {
	w, h := g.rect.Size()
	g.Drawable.Layout()
	x, y := 0, 0
	for x < w {
		vector.StrokeLine(g.Image(), float32(x), float32(y), float32(x), float32(h), 1, colornames.Greenyellow, true)
		x += g.spacing
	}
	x = 0
	for y < h {
		vector.StrokeLine(g.Image(), float32(x), float32(y), float32(w), float32(y), 1, colornames.Greenyellow, true)
		y += g.spacing
	}
	g.ClearDirty()
	log.Println("GridBackground:Layout:", g.spacing, g.Rect())
}
