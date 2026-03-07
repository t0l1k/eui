package crud

import (
	"log"
	"strings"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/data"
)

type Person struct {
	First string
	Last  string
}

func NewCrudDemo(fn func(*eui.Button)) *eui.Container {
	var (
		filterInp, firstInput, lastInput *eui.InputLine
		list                             *eui.ListView
		createBtn, updateBtn, deleteBtn  *eui.Button
		selectedIdx                      int = -1
	)

	people := []Person{
		{"Hans", "Emil"},
		{"Max", "Mustermann"},
		{"Roman", "Tisch"},
	}

	updateList := func() {
		f := filterInp.Text()
		list.Reset()
		for i, p := range people {
			if f == "" || strings.Contains(p.First, f) || strings.Contains(p.Last, f) {
				idx := i
				list.AddItem(eui.NewButton(p.First+" "+p.Last, func(b *eui.Button) {
					log.Println("Pressed:", b.Text())
					selectedIdx = idx
					firstInput.SetText(people[idx].First)
					lastInput.SetText(people[idx].Last)
				}))
			}
		}
	}

	filterInp = eui.NewInputLine(func(il *eui.InputLine) {}).SetPlaceholder("Filter")
	filterInp.TextChanged().Connect(func(data string) { updateList() })
	firstInput = eui.NewInputLine(func(il *eui.InputLine) {}).SetPlaceholder("First name")
	lastInput = eui.NewInputLine(func(il *eui.InputLine) {}).SetPlaceholder("Last name")

	list = eui.NewListView()

	createBtn = eui.NewButton("Create", func(b *eui.Button) {
		p := Person{First: "."}
		people = append(people, p)
		updateList()
	})
	updateBtn = eui.NewButton("Update", func(b *eui.Button) {
		if selectedIdx >= 0 && selectedIdx < len(people) {
			people[selectedIdx].First = firstInput.Text()
			people[selectedIdx].Last = lastInput.Text()
			updateList()
		}
	})
	deleteBtn = eui.NewButton("Delete", func(b *eui.Button) {
		if selectedIdx >= 0 && selectedIdx < len(people) {
			people = append(people[:selectedIdx], people[selectedIdx+1:]...)
			selectedIdx = -1
			firstInput.SetText("")
			lastInput.SetText("")
			updateList()
		}
	})

	updateList()

	btns := eui.NewContainer(eui.NewHBoxLayout(1)).Add(createBtn).Add(updateBtn).Add(deleteBtn)
	contB := eui.NewContainer(eui.NewVBoxLayout(1)).Add(firstInput).Add(lastInput)
	contA := eui.NewContainer(eui.NewHBoxLayout(1)).Add(list).Add(contB)
	cont := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{15, 80, 5}, 1)).Add(filterInp).Add(contA).Add(btns)

	scene := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 95}, 1))
	scene.Add(eui.NewTopBar(data.Crud, fn).SetButtonText(data.QuitDemo))
	scene.Add(cont)
	return scene
}
