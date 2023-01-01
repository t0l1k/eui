package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type SceneMain struct {
	ui.ContainerDefault
	topBar   *TopBar
	lblDiff  *ui.Label
	btnsDiff []*ui.Button
}

func NewSceneMain() *SceneMain {
	s := &SceneMain{}
	s.topBar = NewTopBar(ui.GetUi().GetTitle())
	s.Add(s.topBar)
	rect := []int{0, 0, 1, 1}
	s.lblDiff = ui.NewLabel("Select Difficult:", rect, ui.Aqua, ui.Black)
	s.Add(s.lblDiff)
	strsDiff := []string{"Novice", "Begginer", "Intermediate", "Expert", "Custom"}
	for _, str := range strsDiff {
		btn := ui.NewButton(str, rect, ui.GreenYellow, ui.Black, s.selectLogic)
		s.btnsDiff = append(s.btnsDiff, btn)
		s.Add(btn)
	}
	return s
}

func (s *SceneMain) selectLogic(b *ui.Button) {
	switch b.GetText() {
	case "Novice":
		sc := NewSceneGame(5, 5, 5)
		ui.Push(sc)
		sc.topBar.lblTitle.SetText(ui.GetUi().GetTitle() + " Novice(5x5x5)")
		ui.GetUi().ShowNotification("Started game 5 columns 5 rows 5 mines")
	case "Begginer":
		sc := NewSceneGame(8, 8, 10)
		ui.Push(sc)
		sc.topBar.lblTitle.SetText(ui.GetUi().GetTitle() + " Begginer")
		ui.GetUi().ShowNotification("Started game 8 columns 8 rows 10 mines")
	case "Intermediate":
		sc := NewSceneGame(16, 16, 40)
		ui.Push(sc)
		sc.topBar.lblTitle.SetText(ui.GetUi().GetTitle() + " Intermediate")
		ui.GetUi().ShowNotification("Started game 16 columns 16 rows 40 mines")
	case "Expert":
		sc := NewSceneGame(40, 16, 99)
		ui.Push(sc)
		sc.topBar.lblTitle.SetText(ui.GetUi().GetTitle() + " Expert")
		ui.GetUi().ShowNotification("Started game 40 columns 16 rows 99 mines")
	case "Custom":
		sc := NewSceneCustom()
		// sc := NewSceneGame(10, 15, 20)
		ui.Push(sc)
	}
}

func (s *SceneMain) Entered() {
	s.Resize()
}

func (s *SceneMain) Update(dt int) {
	for _, c := range s.Container {
		c.Update(dt)
	}
}

func (s *SceneMain) Draw(surface *ebiten.Image) {
	for _, c := range s.Container {
		c.Draw(surface)
	}
}

func (s *SceneMain) Resize() {
	s.topBar.Resize()
	w, h := ebiten.WindowSize()
	hTop := int(float64(h) * 0.05)
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/(len(s.btnsDiff)+2)
	x := rect.CenterX() - w1/2
	y := rect.Y + h1/2
	s.lblDiff.Resize([]int{x, y, w1, h1 - 2})
	for _, btn := range s.btnsDiff {
		y += h1
		btn.Resize([]int{x, y, w1, h1 - 2})
	}
}

func (s *SceneMain) Close() {
	for _, c := range s.Container {
		c.Close()
	}
}
