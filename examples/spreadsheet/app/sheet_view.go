package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/spreadsheet/sheet"
	"golang.org/x/image/colornames"
)

type SpreadsheetView struct {
	*eui.Container
	curCellLbl *eui.Label
	input      *eui.TextInputLine
	layFooter  *eui.Container
	laySheet   *eui.ListView
	sheet      *sheet.Sheet
	activeCell *sheet.Cell
}

func NewSpreadSheetView(row, column int) *SpreadsheetView {
	sv := SpreadsheetView{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	sv.layFooter = eui.NewContainer(eui.NewVBoxLayout(1))
	sv.curCellLbl = eui.NewLabel("empty")
	sv.layFooter.Add(sv.curCellLbl)
	sv.input = eui.NewTextInputLine(func(ib *eui.TextInputLine) {
		value := ib.Text()
		if ok, res := sv.sheet.IsFormula(sv.activeCell, value); ok {
			sv.activeCell.Emit(res)
		} else {
			sv.activeCell.Emit(value)
		}
		sv.laySheet.ImageReset()
	})
	sv.layFooter.Add(sv.input)
	sv.Add(sv.layFooter)

	sv.laySheet = eui.NewListView()
	sv.laySheet.Rows(row + 1)
	sv.Add(sv.laySheet)
	sv.sheet = sheet.NewSheet()
	sv.laySheet.AddItem(eui.NewLabel(" "))
	for i := 0; i < row; i++ {
		grid := sheet.NewGrid(i, 0)
		sv.laySheet.AddItem(eui.NewLabel(grid.GetRow()))
	}
	for y := 0; y < column; y++ {
		grid := sheet.NewGrid(0, y+1)
		lblCol := eui.NewLabel(grid.GetColumn())
		sv.laySheet.AddItem(lblCol)
		for x := 0; x < row; x++ {
			grid := sheet.NewGrid(x, y+1)
			cell := sheet.NewCell(grid)
			btn := eui.NewButton(cell.String(), func(b *eui.Button) {
				for _, v := range sv.laySheet.Items() {
					switch vv := v.(type) {
					case *eui.Button:
						if vv.Bg() == colornames.Aqua {
							vv.SetBg(colornames.Gray)
							vv.SetFg(colornames.Yellow)
							if sv.activeCell.IsActive() {
								sv.activeCell.SetInActive()
								sv.activeCell = nil
							}
							sv.input.SetText("")
						}
					}
				}
				b.SetBg(colornames.Aqua)
				b.SetFg(colornames.Black)
				if cell.IsContainFormula() {
					sv.input.SetText(cell.GetFormula().String())
				} else {
					sv.input.SetText(b.Text())
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

func (s *SpreadsheetView) SetRect(rect eui.Rect[int]) {
	s.Container.SetRect(rect)
	cellSize := float64(s.Rect().GetLowestSize()) * 0.05
	x0, y0, w0, h0 := s.Rect().X, s.Rect().Y, s.Rect().W, s.Rect().H
	x, y, w, h := x0, y0, int(cellSize*10), int(cellSize*2)
	s.layFooter.SetRect(eui.NewRect([]int{x, y, w, h}))
	y += int(cellSize * 2)
	w = w0
	h = h0 - int(cellSize*2)
	s.laySheet.SetRect(eui.NewRect([]int{x, y, w, h}))
}
