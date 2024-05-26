package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GridView struct {
	DrawableBase
	r, c, strokeWidth float64
	DrawRect          bool
	bg, fg            color.Color
}

func NewGridView(row, column float64) *GridView {
	gr := &GridView{r: row, c: column}
	gr.DrawRect = false
	gr.strokeWidth = 1
	gr.Visible(true)
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

func (g *GridView) Set(r, c float64) {
	g.r = r
	g.c = c
	g.Dirty = true
}

func (g *GridView) SetRow(r float64)         { g.r = r; g.Dirty = true }
func (g *GridView) SetColumn(c float64)      { g.c = c; g.Dirty = true }
func (g *GridView) SetStrokewidth(w float64) { g.strokeWidth = w; g.Dirty = true }

func (g *GridView) Layout() {
	g.SpriteBase.Layout()
	g.Image().Fill(g.bg)
	cellSizeW := func() (size float64) {
		r := g.r
		for r*size < float64(g.rect.W) {
			size += 0.01
		}
		return size
	}()

	cellSizeH := func() (size float64) {
		c := g.c
		for c*size < float64(g.rect.H) {
			size += 0.01
		}
		return size
	}()

	r := g.GetRect()
	w0, h0 := r.Size()
	marginX := (float64(w0) - cellSizeW*g.r) / 2
	marginY := (float64(h0) - cellSizeH*g.c) / 2
	x0, y0 := marginX, marginY

	if g.DrawRect {
		x, y := x0, y0
		w, h := cellSizeW*g.r, cellSizeH*g.c
		if y < 0 {
			y = 0
		}
		if h > float64(h0) {
			h = float64(h0)
		}
		vector.StrokeRect(g.Image(), float32(x), float32(y), float32(w), float32(h), float32(g.strokeWidth), g.fg, true)
	}

	if g.c > 1 {
		for x := x0 + cellSizeW; x < x0+cellSizeW*g.r; x += cellSizeW {
			vector.StrokeLine(
				g.Image(),
				float32(x),
				float32(y0),
				float32(x),
				float32(y0)+float32(cellSizeH)*float32(g.c),
				float32(g.strokeWidth), g.fg, true)
		}
	}

	if g.r > 1 {
		for y := y0 + cellSizeH; y < y0+cellSizeH*g.c; y += cellSizeH {
			vector.StrokeLine(
				g.Image(),
				float32(x0),
				float32(y),
				float32(x0)+float32(cellSizeW)*float32(g.r),
				float32(y),
				float32(g.strokeWidth), g.fg, true)
		}
	}
	g.Dirty = true
}

func (g *GridView) Draw(surface *ebiten.Image) {
	if !g.IsVisible() {
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
