package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/spreadsheet/sheet"
	"golang.org/x/image/colornames"
)

const Title = "Spreadsheet example"

func main() {
	eui.Init(eui.GetUi().SetTitle(Title).SetSize(800, 600))
	eui.Run(func() *eui.Scene {
		sc := eui.NewScene(eui.NewVBoxLayout(10))
		sc.Add(func(row, column int) *eui.Container {
			sv := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{30, 70}, 5))
			layFooter := eui.NewContainer(eui.NewVBoxLayout(5))
			curCellLbl := eui.NewLabel("empty")
			layFooter.Add(curCellLbl)
			laySheet := eui.NewListView()
			sh := sheet.NewSheet()
			var activeCell *sheet.Cell
			input := eui.NewTextInputLine(func(ib *eui.TextInputLine) {
				value := ib.Text()
				if ok, res := sh.IsFormula(activeCell, value); ok {
					activeCell.Emit(res)
				} else {
					activeCell.Emit(value)
				}
				laySheet.ImageReset()
			})
			layFooter.Add(input)
			sv.Add(layFooter)

			laySheet.Rows(row + 1)
			sv.Add(laySheet)
			laySheet.AddItem(eui.NewLabel(" "))
			for i := 0; i < row; i++ {
				grid := sheet.NewGrid(i, 0)
				laySheet.AddItem(eui.NewLabel(grid.GetRow()))
			}
			for y := 0; y < column; y++ {
				grid := sheet.NewGrid(0, y+1)
				lblCol := eui.NewLabel(grid.GetColumn())
				laySheet.AddItem(lblCol)
				for x := 0; x < row; x++ {
					grid := sheet.NewGrid(x, y+1)
					cell := sheet.NewCell(grid)
					btn := eui.NewButton(cell.String(), func(b *eui.Button) {
						for _, v := range laySheet.Items() {
							switch vv := v.(type) {
							case *eui.Button:
								if vv.Bg() == colornames.Aqua {
									vv.SetBg(colornames.Gray)
									vv.SetFg(colornames.Yellow)
									if activeCell.IsActive() {
										activeCell.SetInActive()
										activeCell = nil
									}
									input.SetText("")
								}
							}
						}
						b.SetBg(colornames.Aqua)
						b.SetFg(colornames.Black)
						if cell.IsContainFormula() {
							input.SetText(cell.GetFormula().String())
						} else {
							input.SetText(b.Text())
						}
						curCellLbl.SetText(grid.String())
						cell.SetActive()
						activeCell = cell
					})
					cell.Connect(func(data string) {
						btn.SetText(data)
					})
					sh.InitCell(grid, cell)
					laySheet.AddBgFg(btn, colornames.Gray, colornames.Yellow)
				}
			}
			return sv
		}(5, 25))
		return sc
	}())
	eui.Quit(func() {})
}
