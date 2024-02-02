package scene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/lines/game"
)

type BallAnimType int

const (
	BallAnimNo BallAnimType = iota
	BallAnimFilledNext
	BallAnimFilled
	BallAnimJump
	BallAnimFilledAfterMove
	BallAnimDelete
)

type CellIcon struct {
	eui.DrawableBase
	btn        *eui.ButtonIcon
	cell       *game.Cell
	anim       BallAnimType
	animStatus BallStatusType
	jumpDt     int
	icon       *BallIcon
	fg         color.RGBA
}

func NewCellIcon(cell *game.Cell, f func(b *eui.ButtonIcon)) *CellIcon {
	c := &CellIcon{}
	c.cell = cell
	c.icon = NewBallIcon(BallHidden, game.BallNoColor.Color(), game.BallNoColor.Color())
	c.icon.Resize([]int{0, 0, 1, 1})
	c.icon.Visible = true
	c.Visible = true
	c.btn = eui.NewButtonIcon([]*ebiten.Image{c.icon.GetImage(), c.icon.GetImage()}, f)
	c.btn.Bg(game.BallNoColor.Color())
	c.btn.Fg(game.BallNoColor.Color())
	c.anim = BallAnimNo
	c.animStatus = BallHidden
	return c
}

func (c *CellIcon) Update(dt int) {
	c.btn.Update(dt)
	if c.anim == BallAnimFilledNext && !c.animStatus.IsMedium() {
		c.jumpDt += dt
		if c.jumpDt > 250 {
			c.animStatus.FilledNext()
			c.jumpDt = 0
			c.updateIcon(c.animStatus)
			// log.Println("cell icon anim filled next", c.animStatus.String())
		}
	}
	if c.anim == BallAnimJump {
		c.jumpDt += dt
		if c.jumpDt > 100 {
			c.animStatus.Jump()
			c.jumpDt = 0
			c.updateIcon(c.animStatus)
			// log.Println("cell icon anim jump", c.animStatus.String())
		}
	}

	if c.anim == BallAnimFilledAfterMove && !c.animStatus.IsHidden() {
		c.jumpDt += dt
		if c.jumpDt > 250 {
			c.jumpDt = 0
			c.anim = BallAnimFilled
			c.updateIcon(c.animStatus)
			c.cell.SetFilled()
			log.Println("cell icon anim show move way end", c.animStatus.String(), c.cell.IsFilledAfterMove())
		}
	}
	if c.anim == BallAnimDelete {
		if !c.animStatus.IsHidden() {
			c.jumpDt += dt
			if c.jumpDt > 50 {
				c.animStatus.Delete()
				c.jumpDt = 0
				c.updateIcon(c.animStatus)
				log.Println("cell icon anim delete", c.animStatus.String())
			}
		} else {
			c.animStatus = BallHidden
			c.updateIcon(c.animStatus)
			c.anim = BallAnimNo
			c.cell.Reset()
			log.Println("cell icon anim delete done", c.animStatus.String())
		}
	}
}

func (c *CellIcon) UpdateData(value interface{}) {
	switch v := value.(type) {
	case *game.CellData:
		switch v.State {
		case game.CellEmpty:
			c.anim = BallAnimNo
			c.animStatus = BallHidden
			c.fg = v.Color.Color()
			c.updateIcon(BallHidden)
		case game.CellFilledNext:
			c.anim = BallAnimFilledNext
			c.animStatus = BallHidden
			c.fg = v.Color.Color()
			c.updateIcon(BallSmall)
		case game.CellFilled:
			c.anim = BallAnimFilled
			c.animStatus = BallMedium
			c.fg = v.Color.Color()
			c.updateIcon(BallNormal)
		case game.CellMarkedForMove:
			c.anim = BallAnimJump
			c.animStatus = BallJumpCenter
			c.fg = v.Color.Color()
			c.updateIcon(BallNormal)
		case game.CellFilledAfterMove:
			c.anim = BallAnimFilledAfterMove
			c.animStatus = BallJumpCenter
			c.fg = v.Color.Color()
			c.updateIcon(BallNormal)
		case game.CellMarkedForDelete:
			c.anim = BallAnimDelete
			c.animStatus = BallBig
			c.fg = v.Color.Color()
			c.updateIcon(BallBig)
		}
	}
}

func (c *CellIcon) updateIcon(ballStatus BallStatusType) {
	bg := game.BallNoColor.Color()
	rect := c.Rect.GetArr()
	c.icon = NewBallIcon(ballStatus, bg, c.fg)
	c.icon.setup(ballStatus)
	c.icon.Resize(rect)
	c.btn.SetIcons([]*ebiten.Image{c.icon.GetImage(), c.icon.GetImage()})
}

func (c *CellIcon) Draw(surface *ebiten.Image) {
	if !c.Visible {
		return
	}
	c.btn.Draw(surface)
	if c.Dirty {
		c.btn.Layout()
		c.icon.Layout()
	}
	if c.anim == BallAnimJump || c.anim == BallAnimFilled {
		op := &ebiten.DrawImageOptions{}
		x, y := c.Rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(c.icon.Image, op)
	}
}

func (c *CellIcon) Resize(rect []int) {
	c.Rect = eui.NewRect(rect)
	c.btn.Resize(rect)
	c.icon.Resize(rect)
	c.Dirty = true
}
