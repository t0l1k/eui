package clock

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
)

type Hand struct {
	eui.DrawableBase
	faceCenter, tip          eui.Point
	value, lenght, thickness float64
	fg                       color.Color
}

func NewHand(clr color.Color) *Hand {
	h := &Hand{fg: clr}
	return h
}

func (h *Hand) Setup(center eui.Point, lenght float64, thickness float64, visible bool) {
	h.faceCenter = center
	h.lenght = lenght
	h.thickness = thickness
	h.Visible = visible
}

func (h *Hand) ToggleVisible() {
	h.Visible = !h.Visible
	h.Dirty = true
}

func (h *Hand) Get() float64 {
	return h.value
}

func (h *Hand) Set(value float64) {
	if h.value == value {
		return
	}
	h.value = value
	h.tip = GetTip(h.faceCenter, h.value, h.lenght, 0, 0)
	h.Dirty = true
}

func (h *Hand) Draw(surface *ebiten.Image) {
	if !h.Dirty || !h.Visible {
		return
	}
	x := h.Rect.X
	y := h.Rect.Y
	x1 := x + h.Rect.CenterX()
	y1 := y + h.Rect.CenterY()
	x2 := x + int(h.tip.X)
	y2 := y + int(h.tip.Y)
	vector.StrokeLine(surface, float32(x1), float32(y1), float32(x2), float32(y2), float32(h.thickness), h.fg, true)
	h.Dirty = false
}

func (h *Hand) Resize(r []int) {
	h.Rect = eui.NewRect(r)
	h.Dirty = true
}
