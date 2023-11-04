package clock

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
)

type Hand struct {
	eui.DrawableBase
	faceCenter, tip                 eui.Point
	value, lenght, thickness, angle float64
	fg                              color.Color
}

func NewHand() *Hand {
	h := &Hand{}
	return h
}

func (h *Hand) Setup(center eui.Point, lenght float64, fg color.Color, thickness float64) {
	h.faceCenter = center
	h.lenght = lenght
	h.fg = fg
	h.thickness = thickness
}

func (h *Hand) GetTip() eui.Point {
	return h.tip
}

func (h *Hand) SetTip() {
	h.tip = GetTip(h.faceCenter, h.value, h.lenght, 0, 0)
}

func (h *Hand) Get() float64 {
	return h.value
}

func (h *Hand) Set(value float64) {
	if h.value == value {
		return
	}
	h.value = value
}

func (h *Hand) GetAngle() float64 {
	return h.angle
}

func (h *Hand) SetAngle() {
	h.angle = GetAngle(h.value)
}

func (h *Hand) Draw(surface *ebiten.Image) {
	x := h.Rect.X
	y := h.Rect.Y
	x1 := x + h.Rect.CenterX()
	y1 := y + h.Rect.CenterY()
	x2 := x + int(h.tip.X)
	y2 := y + int(h.tip.Y)
	vector.StrokeLine(surface, float32(x1), float32(y1), float32(x2), float32(y2), float32(h.thickness), h.fg, true)
}

func (h *Hand) Resize(r []int) {
	h.Rect = eui.NewRect(r)
}
