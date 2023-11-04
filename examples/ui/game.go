package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	button1    *Button
	button2    *Button
	checkBox   *CheckBox
	textBoxLog *TextBox
}

func NewGame() *Game {
	g := &Game{}
	g.button1 = &Button{
		Rect: image.Rect(16, 16, 144, 48),
		Text: "Button 1",
	}
	g.button2 = &Button{
		Rect: image.Rect(160, 16, 288, 48),
		Text: "Button 2",
	}
	g.checkBox = &CheckBox{
		X:    16,
		Y:    64,
		Text: "Check Box!",
	}
	g.textBoxLog = &TextBox{
		Rect: image.Rect(16, 96, 624, 464),
	}

	g.button1.SetOnPressed(func(b *Button) {
		g.textBoxLog.AppendLine("Button 1 Pressed")
	})
	g.button2.SetOnPressed(func(b *Button) {
		g.textBoxLog.AppendLine("Button 2 Pressed")
	})
	g.checkBox.SetOnCheckChanged(func(c *CheckBox) {
		msg := "Check box check changed"
		if c.Checked() {
			msg += " (Checked)"
		} else {
			msg += " (Unchecked)"
		}
		g.textBoxLog.AppendLine(msg)
	})
	return g
}

func (g *Game) Update() error {
	g.button1.Update()
	g.button2.Update()
	g.checkBox.Update()
	g.textBoxLog.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xeb, 0xeb, 0xeb, 0xff})
	g.button1.Draw(screen)
	g.button2.Draw(screen)
	g.checkBox.Draw(screen)
	g.textBoxLog.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
