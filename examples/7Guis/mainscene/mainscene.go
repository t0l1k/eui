package mainscene

import (
	"log"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/cells"
	"github.com/t0l1k/eui/examples/7Guis/circle"
	"github.com/t0l1k/eui/examples/7Guis/counter"
	"github.com/t0l1k/eui/examples/7Guis/crud"
	"github.com/t0l1k/eui/examples/7Guis/data"
	"github.com/t0l1k/eui/examples/7Guis/flightbooker"
	"github.com/t0l1k/eui/examples/7Guis/temper"
	"github.com/t0l1k/eui/examples/7Guis/timer"
)

func NewMainScene() *eui.Scene {
	var (
		dialogSelect, counterDemo, temperatureDemo, flightBookerDemo, timerDemo, crudDemo, circleDemo, cellsDemo eui.Drawabler
	)
	counterDemo = counter.NewCounterDemo(func(b *eui.Button) {
		switch b.Text() {
		case data.QuitDemo:
			counterDemo.Hide()
			dialogSelect.Show()
		}
		log.Println("CounterDemo", b.Text())
	})

	temperatureDemo = temper.NewTemperatureDemo(func(b *eui.Button) {
		switch b.Text() {
		case data.QuitDemo:
			temperatureDemo.Hide()
			dialogSelect.Show()
		}
		log.Println("TemperatureDemo", b.Text())
	})

	flightBookerDemo = flightbooker.NewFlightBookerDemo(func(b *eui.Button) {
		switch b.Text() {
		case data.QuitDemo:
			flightBookerDemo.Hide()
			dialogSelect.Show()
		}
		log.Println("FlightBookerDemo", b.Text())
	})

	timerDemo = timer.NewTimerDemo(func(b *eui.Button) {
		switch b.Text() {
		case data.QuitDemo:
			timerDemo.Hide()
			dialogSelect.Show()
		}
		log.Println("TimerDemo", b.Text())
	})

	crudDemo = crud.NewCrudDemo(func(b *eui.Button) {
		switch b.Text() {
		case data.QuitDemo:
			crudDemo.Hide()
			dialogSelect.Show()
		}
		log.Println("CrudDemo", b.Text())
	})

	circleDemo = circle.NewCircleDemo(func(b *eui.Button) {
		switch b.Text() {
		case data.QuitDemo:
			circleDemo.Hide()
			dialogSelect.Show()
		}
		log.Println("CircleDemo", b.Text())
	})

	cellsDemo = cells.NewCellsDemo(func(b *eui.Button) {
		switch b.Text() {
		case data.QuitDemo:
			cellsDemo.Hide()
			dialogSelect.Show()
		}
		log.Println("CellsDemo", b.Text())
	})

	c := eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 95}, 1))
	c.Add(eui.NewTopBar("7Guis", nil))

	dialogSelect = NewSelect7Guis(func(b *eui.Button) {
		switch b.Text() {
		case data.Counter:
			dialogSelect.Hide()
			counterDemo.Show()
		case data.TempConv:
			dialogSelect.Hide()
			temperatureDemo.Show()
		case data.FlightBooker:
			dialogSelect.Hide()
			flightBookerDemo.Show()
		case data.Timer:
			dialogSelect.Hide()
			timerDemo.Show()
		case data.Crud:
			dialogSelect.Hide()
			crudDemo.Show()
		case data.Circle:
			dialogSelect.Hide()
			circleDemo.Show()
		case data.Cells:
			dialogSelect.Hide()
			cellsDemo.Show()
		}
		log.Println("7Guis", b.Text())
	})

	board := eui.NewContainer(eui.NewStackLayout(5))
	board.Add(dialogSelect)
	board.Add(counterDemo)
	board.Add(temperatureDemo)
	board.Add(flightBookerDemo)
	board.Add(timerDemo)
	board.Add(crudDemo)
	board.Add(circleDemo)
	board.Add(cellsDemo)
	c.Add(board)
	return c
}

func NewSelect7Guis(fn func(*eui.Button)) *eui.Container {
	c := eui.NewContainer(eui.NewVBoxLayout(5))
	for _, txt := range []string{data.Counter, data.TempConv, data.FlightBooker, data.Timer, data.Crud, data.Circle, data.Cells} {
		c.Add(eui.NewButton(txt, fn))
	}
	return c
}
