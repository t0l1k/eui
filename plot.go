package eui

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type Plot struct {
	*Drawable
	xArr, yArr, values             []int
	bg, fg, fgValues               color.Color
	pTitle, xAxisTitle, yAxisTitle string
}

func NewPlot(xArr, yArr, values []int, title, xTitle, yTitle string) *Plot {
	p := &Plot{Drawable: NewDrawable(), xArr: xArr, yArr: yArr, values: values, pTitle: title, xAxisTitle: xTitle, yAxisTitle: yTitle}
	p.bg = colornames.Gray
	p.fg = colornames.Yellowgreen
	p.fgValues = colornames.Blue
	return p
}

func (p *Plot) Layout() {
	p.Drawable.Layout()
	p.Image().Fill(p.bg)
	axisXMax := len(p.xArr)
	axisYMax := len(p.yArr)
	w0, h0 := p.Rect().Size()
	margin := int(float64(p.Rect().GetLowestSize()) * 0.05)
	x, y, w, h := margin, margin, w0-margin*2, h0-margin*2
	axisRect := NewRect([]int{x, y, w, h})
	stroke := float32(axisRect.GetLowestSize()) * 0.003
	if stroke < 0.5 {
		stroke = 0.5
	}

	lerp := func(t, inStart, inEnd, outStart, outEnd float64) float64 {
		return outStart + (t-inStart)/(inEnd-inStart)*(outEnd-outStart)
	}
	xPos := func(x float64) float64 {
		return math.Round(lerp(x, 0, float64(axisXMax), float64(axisRect.Left()), float64(axisRect.Right())))
	}
	yPos := func(y float64) float64 {
		return math.Round(lerp(y, 0, float64(axisYMax), float64(axisRect.Bottom()), float64(axisRect.Top())))
	}
	zip := func(a, b []int) [][]int {
		if len(a) != len(b) {
			panic("len(a) != len(b)")
		}
		r := make([][]int, 0)
		for i := 0; i < len(a); i++ {
			arr := make([]int, 0)
			arr = append(arr, a[i], b[i])
			r = append(r, arr)
		}
		return r
	}

	{ // X Axis
		x1, y1 := axisRect.BottomLeft()
		x2, y2 := axisRect.BottomRight()
		vector.StrokeLine(p.Image(), float32(x1), float32(y1), float32(x2), float32(y2), stroke, p.fg, true)

		xTicks := len(p.xArr)
		for i := 1; i < xTicks+1; i++ {
			x := axisXMax * i / xTicks
			x1, y1 := xPos(float64(x)), float64(axisRect.Bottom())
			x2, y2 := xPos(float64(x)), float64(axisRect.Bottom()+margin/4)
			vector.StrokeLine(p.Image(), float32(x1), float32(y1), float32(x2), float32(y2), stroke, p.fg, true)

			x1, y1 = xPos(float64(x)), float64(axisRect.Bottom())
			x2, y2 = xPos(float64(x)), float64(axisRect.Top())
			vector.StrokeLine(p.Image(), float32(x1), float32(y1), float32(x2), float32(y2), stroke, p.fg, true)
			boxSize := margin / 2
			xL, yL, w, h := int(xPos(float64(x))-float64(boxSize)/2), axisRect.Bottom()+boxSize/2, boxSize, boxSize
			lbl := NewLabel(strconv.Itoa(p.xArr[i-1]))
			defer lbl.Close()
			lbl.SetBg(p.bg)
			lbl.SetFg(p.fg)
			lbl.SetRect(NewRect([]int{xL, yL, w, h}))
			lbl.Draw(p.Image())
		}
		boxSize := margin
		xL, yL, w, h := axisRect.Right()-boxSize*3, axisRect.Bottom()-boxSize, boxSize*3, boxSize
		lbl := NewLabel(p.xAxisTitle)
		defer lbl.Close()
		lbl.SetBg(p.bg)
		lbl.SetFg(p.fg)
		lbl.SetRect(NewRect([]int{xL, yL, w, h}))
		lbl.Draw(p.Image())
	}
	{ // Y Axis
		x1, y1 := axisRect.BottomLeft()
		x2, y2 := axisRect.TopLeft()
		vector.StrokeLine(p.Image(), float32(x1), float32(y1), float32(x2), float32(y2), stroke, p.fg, true)
		yTicks := len(p.yArr)
		for i := 1; i < yTicks+1; i++ {
			y := axisYMax * i / yTicks
			x1, y1 := axisRect.Left(), yPos(float64(y))
			x2, y2 := axisRect.Left()-margin/4, yPos(float64(y))
			vector.StrokeLine(p.Image(), float32(x1), float32(y1), float32(x2), float32(y2), stroke, p.fg, true)
			x1, y1 = axisRect.Left(), yPos(float64(y))
			x2, y2 = axisRect.Right(), yPos(float64(y))
			vector.StrokeLine(p.Image(), float32(x1), float32(y1), float32(x2), float32(y2), stroke, p.fg, true)
			boxSize := margin / 2
			xL, yL, w, h := axisRect.Left()-int(float64(boxSize)*1.5), int(yPos(float64(y))-float64(boxSize)/2), boxSize, boxSize
			lbl := NewLabel(strconv.Itoa(p.yArr[i-1]))
			defer lbl.Close()
			lbl.SetBg(p.bg)
			lbl.SetFg(p.fg)
			lbl.SetRect(NewRect([]int{xL, yL, w, h}))
			lbl.Draw(p.Image())
		}
		boxSize := margin
		xL, yL, w, h := axisRect.Left(), axisRect.Top()-boxSize, boxSize*3, boxSize
		lbl := NewLabel(p.yAxisTitle)
		defer lbl.Close()
		lbl.SetBg(p.bg)
		lbl.SetFg(p.fg)
		lbl.SetRect(NewRect([]int{xL, yL, w, h}))
		lbl.Draw(p.Image())
	}
	{ // values line
		points := zip(p.xArr, p.values)
		var results []float64
		xx := xPos(float64(axisXMax) * float64(0) / float64(len(p.xArr)))
		yy := yPos(float64(0))
		results = append(results, xx, yy)
		for i := 0; i < len(points); i++ {
			x, y := points[i][0], points[i][1]
			xx := xPos(float64(axisXMax) * float64(x) / float64(len(p.xArr)))
			yy := yPos(float64(y))
			results = append(results, xx, yy)
		}
		for i, j := 0, 1; i < len(results)-2; i, j = i+2, j+2 {
			x1, y1, x2, y2 := results[i], results[j], results[i+2], results[j+2]
			vector.StrokeLine(p.Image(), float32(x1), float32(y1), float32(x2), float32(y2), stroke, p.fgValues, true)
		}
	}
	{
		boxSize := margin * 8
		xL, yL, w, h := axisRect.Right()/2-boxSize/2, axisRect.Top()-boxSize/10, boxSize, boxSize/4
		lbl := NewLabel(p.pTitle)
		defer lbl.Close()
		lbl.SetBg(p.bg)
		lbl.SetFg(p.fg)
		lbl.SetRect(NewRect([]int{xL, yL, w, h}))
		lbl.Draw(p.Image())
	}
	p.ClearDirty()
}

func (p *Plot) Draw(surface *ebiten.Image) {
	if p.IsHidden() {
		return
	}
	if p.IsDirty() {
		p.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := p.Rect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(p.Image(), op)
}

func (p *Plot) SetRect(rect Rect[int]) {
	p.Drawable.SetRect(rect)
	p.ImageReset()
}
