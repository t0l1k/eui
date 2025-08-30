package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/t0l1k/eui"
)

func main() {
	eui.Init(func() *eui.Ui {
		k := 2
		w, h := 500*k, 200*k
		return eui.GetUi().SetTitle("Test ListView").SetSize(w, h)
	}())
	eui.Run(func() *eui.Scene {
		s := eui.NewScene(eui.NewLayoutHorizontalPercent([]int{20, 20, 20, 40}, 10))
		var list []string
		for i := 0; i < 54; i++ {
			list = append(list, "Item "+strconv.Itoa(i))
		}
		theme := eui.GetUi().Theme()
		bg := theme.Get(eui.ListViewItemBg)
		fg := theme.Get(eui.ListViewItemFg)
		lstText := eui.NewListView()
		lstText.SetupListViewText(list, 30, 2, bg, fg)

		lstButtons := eui.NewListView()
		lstButtons.SetupListViewButtons(list, 30, 1, bg, fg, func(b *eui.Button) {
			log.Println("pressed:", b.Text())
		})

		lstCheckBoxs := eui.NewListView()
		lstCheckBoxs.SetupListViewCheckBoxs(list, 30, 1, bg, fg, func(b *eui.Checkbox) {
			log.Println("pressed:", b.Text())
		})

		btnRemoveSelected := eui.NewButton("Remove Selected", func(b *eui.Button) {
			list = nil
			for _, v := range lstCheckBoxs.GetCheckBoxes() {
				if v.IsChecked() {
					fmt.Println("selected:", v.Text())
					continue
				}
				list = append(list, v.Text())
			}
			lstCheckBoxs.Reset()
			lstCheckBoxs.SetupListViewCheckBoxs(list, 30, 1, bg, fg, func(b *eui.Checkbox) {
				log.Println("pressed:", b.Text())
			})

		})

		lstLongText := eui.NewListView()
		lstLongText.AddLongText(eui.NewLabel("Вот мой секрет, он очень прост: зорко одно лишь сердце. Самого главного глазами не увидишь. Тщеславные люди глухи ко всему, кроме похвал. Встал поутру, умылся, привел себя в порядок — и сразу же приведи в порядок свою планету. Если идти все прямо да прямо, далеко не уйдешь.Когда даешь себя приручить, потом случается и плакать.Ведь все взрослые сначала были детьми, только мало кто из них об этом помнит. Уж такой народ эти взрослые. Не стоит на них сердиться. Дети должны быть очень снисходительны к взрослым. Когда очень хочешь сострить, иной раз поневоле приврешь. Себя судить куда труднее, чем других. Если ты сумеешь правильно судить себя, значит, ты поистине мудр. Слова только мешают понимать друг друга. Взрослые никогда ничего не понимают сами, а для детей очень утомительно без конца им все объяснять и растолковывать. С каждого надо спрашивать то, что он может дать. Власть прежде всего должна быть разумной. Люди забыли эту истину, но ты не забывай: ты навсегда в ответе за всех, кого приручил. Ты в ответе за твою розу. Должна же я стерпеть двух-трех гусениц, если хочу познакомиться с бабочками. Будь то дом, звезды или пустыня — самое прекрасное в них то, чего не увидишь глазами. Люди забираются в скорые поезда, но они уже сами не понимают, чего ищут, поэтому они не знают покоя и бросаются то в одну сторону, то в другую. И все напрасно. \n\nМаленький принц \nАнтуан де Сент-Экзюпери ").SetFontSize(18).SetAlign(eui.LabelAlignLeftUp))

		contCheck := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{90, 10}, 5))
		contCheck.Add(lstCheckBoxs)
		contCheck.Add(btnRemoveSelected)

		s.Add(lstText)
		s.Add(lstButtons)
		s.Add(contCheck)
		s.Add(lstLongText)
		return s
	}())
	eui.Quit(func() {})
}
