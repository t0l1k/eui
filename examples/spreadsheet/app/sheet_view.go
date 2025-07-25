package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/spreadsheet/sheet"
	"golang.org/x/image/colornames"
)

type SpreadsheetView struct {
	*eui.Container
	curCellLbl *eui.Text
	input      *eui.InputBox2
	layFooter  *eui.Container
	laySheet   *eui.ListView
	sheet      *sheet.Sheet
	activeCell *sheet.Cell
}

func NewSpreadSheetView(row, column int) *SpreadsheetView {
	sv := SpreadsheetView{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	sv.layFooter = eui.NewContainer(eui.NewVBoxLayout(1))
	sv.curCellLbl = eui.NewText("empty")
	sv.layFooter.Add(sv.curCellLbl)
	sv.input = eui.NewInputBox2(func(ib *eui.InputBox2) {
		value := ib.GetText()
		if ok, res := sv.sheet.IsFormula(sv.activeCell, value); ok {
			sv.activeCell.Emit(res)
		} else {
			sv.activeCell.Emit(value)
		}
		sv.laySheet.ImageReset()
	})
	sv.input.SetFocus()
	sv.input.SetAlwaysFocus()
	sv.layFooter.Add(sv.input)
	sv.Add(sv.layFooter)

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
				for _, v := range sv.laySheet.Childrens() {
					switch vv := v.(type) {
					case *eui.Button:
						if vv.GetBg() == colornames.Aqua {
							vv.Bg(colornames.Gray)
							vv.Fg(colornames.Yellow)
							if sv.activeCell.IsActive() {
								sv.activeCell.SetInActive()
								sv.activeCell = nil
							}
							sv.input.Reset()
						}
					}
				}
				b.Bg(colornames.Aqua)
				b.Fg(colornames.Black)
				if cell.IsContainFormula() {
					sv.input.SetText(cell.GetFormula().String())
				} else {
					sv.input.SetText(b.GetText())
				}
				sv.curCellLbl.SetText(grid.String())
				cell.SetActive()
				sv.activeCell = cell
			})
			cell.Connect(func(data string) {
				btn.SetText(data)
			})
			sv.sheet.InitCell(grid, cell)
			sv.laySheet.AddBgFg(btn, colornames.Gray, colornames.Yellow)
		}
	}
	return &sv
}

func (s *SpreadsheetView) Update(dt int) {
	s.Container.Update(dt)
}

func (s *SpreadsheetView) Resize(rect eui.Rect[int]) {
	s.SetRect(rect)
	cellSize := float64(s.Rect().GetLowestSize()) * 0.05
	x0, y0, w0, h0 := s.Rect().X, s.Rect().Y, s.Rect().W, s.Rect().H
	x, y, w, h := x0, y0, int(cellSize*10), int(cellSize*2)
	s.layFooter.Resize(eui.NewRect([]int{x, y, w, h}))
	y += int(cellSize * 2)
	w = w0
	h = h0 - int(cellSize*2)
	s.laySheet.Resize(eui.NewRect([]int{x, y, w, h}))
}
