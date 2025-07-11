package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
)

type BallIcon struct {
	*eui.Drawable
	size   float32
	status BallStatusType
	bg, fg color.RGBA
}

func NewBallIcon(status BallStatusType, bg, fg color.RGBA) *BallIcon {
	i := &BallIcon{
		Drawable: eui.NewDrawable(),
		bg:       bg,
		fg:       fg,
		status:   status,
	}
	i.setup(status, bg, fg)
	return i
}

func (i *BallIcon) setup(status BallStatusType, bg, fg color.RGBA) {
	i.status = status
	switch status {
	case BallHidden:
		i.size = 0
	case BallSmall:
		i.size = 0.146
	case BallMedium:
		i.size = 0.236
	case BallNormal, BallJumpCenter, BallJumpUp, BallJumpDown:
		i.size = 0.382
	case BallBig:
		i.size = 0.5
	}
	i.bg = bg
	i.fg = fg
	i.MarkDirty()
}

func (i *BallIcon) Layout() {
	i.Drawable.Layout()
	r, g, b, _ := i.bg.RGBA()
	a := 255
	bg := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	i.Image().Fill(bg)
	if i.size > 0 {
		rad := float32(i.Rect().GetLowestSize()) * i.size
		x, y := float32(i.Rect().W/2), float32(i.Rect().H/2)
		margin := (float32(i.Rect().W) - rad*2) / 3
		switch i.status {
		case BallJumpUp:
			y = float32(i.Rect().H/2) - margin
		case BallJumpCenter:
			y = float32(i.Rect().H / 2)
		case BallJumpDown:
			y = float32(i.Rect().H/2) + margin
		}
		vector.DrawFilledCircle(i.Image(), x, y, rad, i.fg, true)
	}
	i.ClearDirty()
}

func (i *BallIcon) GetImage() *ebiten.Image {
	if i.IsDirty() {
		i.Layout()
	}
	return i.Image()
}

func (i *BallIcon) Resize(rect eui.Rect) {
	i.SetRect(rect)
	i.ImageReset()
}
