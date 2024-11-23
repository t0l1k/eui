package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
	"github.com/t0l1k/eui/examples/spreadsheet/sheet"
)

type SpreadsheetView struct {
	eui.DrawableBase
	curCellLbl *eui.Text
	input      *eui.InputBox2
	layFooter  *eui.BoxLayout
	laySheet   *eui.ListView
	sheet      *sheet.Sheet
	activeCell *sheet.Cell
}

func NewSpreadSheetView(row, column int) *SpreadsheetView {
	sv := SpreadsheetView{}
	sv.layFooter = eui.NewVLayout()
	sv.curCellLbl = eui.NewText("empty")
	sv.layFooter.Add(sv.curCellLbl)
	sv.input = eui.NewInputBox2(func(ib *eui.InputBox2) {
		value := ib.GetText()
		if ok, res := sv.sheet.IsFormula(sv.activeCell, value); ok {
			sv.activeCell.SetValue(res)
		} else {
			sv.activeCell.SetValue(value)
		}
		sv.laySheet.ImageReset()
	})
	sv.input.SetFocus()
	sv.input.SetAlwaysFocus()
	sv.layFooter.Add(sv.input)

	sv.laySheet = eui.NewListView()
	sv.laySheet.Rows(row + 1)
	sv.Add(sv.laySheet)
	sv.sheet = sheet.NewSheet()
	sv.laySheet.Add(eui.NewText(" "))
	for i := 0; i < row; i++ {
		grid := sheet.NewGrid(i, 0)
		sv.laySheet.Add(eui.NewText(grid.GetRow()))
	}
	for y := 0; y < column; y++ {
		grid := sheet.NewGrid(0, y+1)
		lblCol := eui.NewText(grid.GetColumn())
		sv.laySheet.Add(lblCol)
		for x := 0; x < row; x++ {
			grid := sheet.NewGrid(x, y+1)
			cell := sheet.NewCell(grid)
			btn := eui.NewButton(cell.String(), func(b *eui.Button) {
				for _, v := range sv.laySheet.GetContainer() {
					switch vv := v.(type) {
					case *eui.Button:
						if vv.GetBg() == colors.Aqua {
							vv.Bg(colors.Gray)
							vv.Fg(colors.Yellow)
							if sv.activeCell.IsActive() {
								sv.activeCell.SetInActive()
								sv.activeCell = nil
							}
							sv.input.Reset()
						}
					}
				}
				b.Bg(colors.Aqua)
				b.Fg(colors.Black)
				if cell.IsContainFormula() {
					sv.input.SetText(cell.GetFormula().String())
				} else {
					sv.input.SetText(b.GetText())
				}
				sv.curCellLbl.SetText(grid.String())
				cell.SetActive()
				sv.activeCell = cell
			})
			cell.Attach(btn)
			sv.sheet.InitCell(grid, cell)
			sv.laySheet.AddBgFg(btn, colors.Gray, colors.Yellow)
		}
	}
	return &sv
}

func (s *SpreadsheetView) Update(dt int) {
	for _, v := range s.layFooter.GetContainer() {
		v.Update(dt)
	}
	for _, v := range s.GetContainer() {
		v.Update(dt)
	}
}

func (s *SpreadsheetView) Draw(surface *ebiten.Image) {
	for _, v := range s.layFooter.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range s.GetContainer() {
		v.Draw(surface)
	}
}

func (s *SpreadsheetView) Resize(rect []int) {
	s.Rect(eui.NewRect(rect))
	cellSize := float64(s.GetRect().GetLowestSize()) * 0.05
	x0, y0, w0, h0 := s.GetRect().X, s.GetRect().Y, s.GetRect().W, s.GetRect().H
	x, y, w, h := x0, y0, int(cellSize*10), int(cellSize*2)
	s.layFooter.Resize([]int{x, y, w, h})
	y += int(cellSize * 2)
	w = w0
	h = h0 - int(cellSize*2)
	s.laySheet.Resize([]int{x, y, w, h})
}
