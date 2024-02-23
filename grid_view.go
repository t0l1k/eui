package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GridView struct {
	DrawableBase
	r, c, strokeWidth int
	DrawRect          bool
	bg, fg            color.Color
}

func NewGridView(row, column int) *GridView {
	gr := &GridView{r: row, c: column}
	gr.DrawRect = false
	gr.strokeWidth = 1
	gr.Visible = true
	return gr
}

func (gr *GridView) Bg(clr color.Color) {
	r, g, b, _ := clr.RGBA()
	a := 0 // invisible
	gr.bg = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	gr.Dirty = true
}

func (gr *GridView) Fg(clr color.Color) {
	gr.fg = clr
	gr.Dirty = true
}

func (g *GridView) Set(r, c int) {
	g.r = r
	g.c = c
	g.Dirty = true
}

func (d *GridView) SetRow(r int)         { d.r = r }
func (d *GridView) SetColumn(c int)      { d.c = c }
func (g *GridView) SetStrokewidth(w int) { g.strokeWidth = w }

func (g *GridView) Layout() {
	g.SpriteBase.Layout()
	g.Image().Fill(g.bg)
	cellSize := func() (size int) {
		r := g.r
		c := g.c
		for r*size < g.rect.W && c*size < g.rect.H {
			size += 1
		}
		return size
	}()

	r := g.GetRect()
	w0, h0 := r.Size()
	marginX := (w0 - cellSize*g.r) / 2
	marginY := (h0 - cellSize*g.c) / 2
	x0, y0 := marginX, marginY

	if g.DrawRect {
		x, y := x0, y0
		w, h := cellSize*g.r, cellSize*g.c
		if y < 0 {
			y = 0
		}
		if h > h0 {
			h = h0
		}
		vector.StrokeRect(g.Image(), float32(x), float32(y), float32(w), float32(h), float32(g.strokeWidth), g.fg, true)
	}

	for y := y0 + cellSize; y < y0+cellSize*g.c; y += cellSize {
		vector.StrokeLine(
			g.Image(),
			float32(x0),
			float32(y),
			float32(x0)+float32(cellSize)*float32(g.r),
			float32(y),
			float32(g.strokeWidth), g.fg, true)
	}

	for x := x0 + cellSize; x < x0+cellSize*g.r; x += cellSize {
		vector.StrokeLine(
			g.Image(),
			float32(x),
			float32(y0),
			float32(x),
			float32(y0)+float32(cellSize)*float32(g.c),
			float32(g.strokeWidth), g.fg, true)
	}
	g.Dirty = true
}

func (g *GridView) Draw(surface *ebiten.Image) {
	if !g.Visible {
		return
	}
	if g.Dirty {
		g.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := g.GetRect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(g.Image(), op)
}

func (g *GridView) Resize(rect []int) {
	g.Rect(NewRect(rect))
	g.SpriteBase.Rect(NewRect(rect))
	g.Dirty = true
	g.ImageReset()
}
