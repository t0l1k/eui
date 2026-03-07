package flightbooker

import (
	"log"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/data"
	"golang.org/x/image/colornames"
)

func NewFlightBookerDemo(fn func(*eui.Button)) *eui.Container {
	const (
		FlightOneWay eui.ViewState = iota + 200
		FlightReturn
	)

	var (
		cont, contA, contB *eui.Container
		opt                *OptionButton
		depDate, retDate   *eui.InputLine
		bookBtn            *eui.Button
	)
	opts := []string{"One-way flight", "Return flight"}
	timeLayout := "2006.01.02"

	validate := func() {
		var (
			d1, d2 time.Time
			err    error
		)
		d1, err = time.Parse(timeLayout, depDate.Text())
		if err != nil {
			depDate.SetBg(colornames.Red)
			depDate.SetFg(colornames.Black)
		} else {
			depDate.SetBg(colornames.Blue)
			depDate.SetFg(colornames.Yellow)
		}
		if opt.SelectedText() == opts[0] && err == nil {
			bookBtn.Show()
			return
		}
		d2, err = time.Parse(timeLayout, retDate.Text())
		if err != nil {
			retDate.SetBg(colornames.Red)
			retDate.SetFg(colornames.Black)
		} else {
			retDate.SetBg(colornames.Blue)
			retDate.SetFg(colornames.Yellow)
		}
		if d1.Before(d2) && err == nil {
			bookBtn.Show()
		} else {
			bookBtn.Hide()
		}
	}
	depDate = eui.NewInputLine(func(ib *eui.InputLine) { validate() }).SetPlaceholder(time.Now().Format(timeLayout))
	retDate = eui.NewInputLine(func(ib *eui.InputLine) { validate() }).SetPlaceholder(time.Now().Format(timeLayout))

	flightState := eui.NewSignal(func(a, b eui.ViewState) bool { return a == b })

	opt = NewOptionButton(opts, 0, func(ob *OptionButton) {
		if ob.SelectedText() == opts[1] { // "Return flight"
			flightState.Emit(FlightReturn)
		} else {
			flightState.Emit(FlightOneWay)
		}
	})

	bookBtn = eui.NewButton("Book", func(b *eui.Button) {
		log.Println("Booked:", depDate.Text(), retDate.Text())
	})

	contA = eui.NewContainer(eui.NewVBoxLayout(1))
	contA.Add(eui.NewLabel("Departure date")).Add(depDate)

	contB = eui.NewContainer(eui.NewVBoxLayout(1))
	contB.Add(eui.NewLabel("Return date")).Add(retDate)

	// Источник всех возможных полей ввода
	source := eui.NewContainer(eui.NewVBoxLayout(0))
	source.Add(contA)
	source.Add(contB)

	// Контейнер, который показывает поля в зависимости от состояния
	visibleInputs := eui.NewContainerVisibleByFilter(eui.NewVBoxLayout(1), source, func(d eui.Drawabler, s eui.ViewState) bool {
		if d == contA {
			return true // Дата вылета нужна всегда
		}
		return s == FlightReturn // Дата возврата только для полета туда-обратно
	})

	flightState.ConnectAndFire(func(s eui.ViewState) {
		visibleInputs.UpdateBy(s)
		validate()
	}, FlightOneWay)

	cont = eui.NewContainer(eui.NewLayoutVerticalPercent([]int{10, 80, 10}, 1))
	cont.Add(opt)
	cont.Add(visibleInputs)
	cont.Add(bookBtn)

	scene := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 95}, 1))
	scene.Add(eui.NewTopBar(data.FlightBooker, fn).SetButtonText(data.QuitDemo))
	scene.Add(cont)
	return scene
}
