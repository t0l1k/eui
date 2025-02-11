package app

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
	"github.com/t0l1k/eui/examples/games/memory_matrix/mem"
)

type BoardMem struct {
	eui.DrawableBase
	layout           *eui.GridLayoutRightDown
	game             *mem.Game
	gamesData        *mem.GamesData
	showTimer        *eui.Timer
	fn               func(*eui.Button)
	varMsg, varColor *eui.SubjectBase
}

func NewBoardMem(fn func(*eui.Button)) *BoardMem {
	d := &BoardMem{}
	d.varMsg = eui.NewSubject()
	d.varColor = eui.NewSubject()
	d.game = mem.NewGame(mem.Level(1))
	d.gamesData = mem.NewGamesData()
	d.Visible(true)
	d.layout = eui.NewGridLayoutRightDown(1, 1)
	d.showTimer = eui.NewTimer(1500)
	d.fn = fn
	d.SetupPreparation()
	return d
}

func (d *BoardMem) Game() *mem.Game { return d.game }

func (d *BoardMem) SetupPreparation() {
	d.layout.ResetContainerBase()
	d.layout.SetDim(1, 1)
	d.layout.FitToDim(false)
	btn := eui.NewButton("Click to Start "+d.game.Level().String()+" "+d.game.Dim().String(), d.fn)
	d.layout.Add(btn)
	d.layout.Resize(d.GetRect().GetArr())
	str := d.gamesData.String()
	d.varMsg.SetValue(str)
	d.varColor.SetValue([]color.Color{colors.YellowGreen, colors.Black})
	log.Println("Setup Preparation done", d.game.String())
}

func (d *BoardMem) SetupShow() {
	d.layout.ResetContainerBase()
	w, h := d.game.Dim().Width(), d.Game().Dim().Height()
	d.layout.SetDim(float64(w), float64(h))
	for y := 0; y < d.Game().Dim().Height(); y++ {
		for x := 0; x < d.game.Dim().Width(); x++ {
			cell := d.Game().Cell(d.Game().Idx(x, y))
			btn := eui.NewButton(" ", d.fn)
			btn.Disable()
			if cell.IsReadOnly() {
				btn.Bg(colors.Orange)
			}
			d.layout.Add(btn)
		}
	}
	d.layout.Resize(d.GetRect().GetArr())
	str := d.gamesData.String()
	d.varMsg.SetValue(str)
	d.varColor.SetValue([]color.Color{colors.Red, colors.Black})
	log.Println("Setup Show Done", d.game.String())
}

func (d *BoardMem) SetupRecolection() {
	d.layout.ResetContainerBase()
	w, h := d.game.Dim().Width(), d.Game().Dim().Height()
	d.layout.SetDim(float64(w), float64(h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			btn := eui.NewButton(" ", d.fn)
			btn.Enable()
			d.layout.Add(btn)
		}
	}
	d.layout.Resize(d.GetRect().GetArr())
	str := d.gamesData.String()
	d.varMsg.SetValue(str)
	d.varColor.SetValue([]color.Color{colors.Blue, colors.Yellow})
	log.Println("Setup Recolection Done", d.game.String())
}

func (d *BoardMem) SetupConclusion() {
	var str string
	d.layout.ResetContainerBase()
	d.layout.SetDim(1, 1)
	d.layout.FitToDim(true)
	sb := eui.NewSnackBar("")
	if d.Game().Win {
		str = "Winner"
		sb.Bg(colors.Blue)
	} else if d.Game().GameOver {
		str = "Game Over"
		sb.Bg(colors.Red)
	}
	sb.SetText(str + " " + d.Game().String()).Show(3000)
	btn := d.setupScoreBtn()
	d.layout.Add(btn)
	d.layout.Resize(d.GetRect().GetArr())
	d.varMsg.SetValue(d.gamesData.String())
	d.varColor.SetValue([]color.Color{colors.Fuchsia, colors.Black})
	d.Add(sb)
	log.Println("Setup Conclusion done", d.game.String())
}

func (d *BoardMem) setupScoreBtn() *eui.ButtonIcon {
	level := d.gamesData.Max()
	count := d.gamesData.Size()
	var xArr, yArr, levels []int
	for i := 0; i < count; i++ {
		xArr = append(xArr, i+1)
	}
	for i := 0; i < level; i++ {
		yArr = append(yArr, i+1)
	}
	levels = d.gamesData.Levels()
	plot := eui.NewPlot(xArr, yArr, levels, "Memory Matrix", "Game", "Level")
	plot.Resize(d.GetRect().GetArr())
	plot.Layout()
	btn := eui.NewButtonIcon([]*ebiten.Image{plot.Image(), plot.Image()}, d.fn)
	return btn
}

func (d *BoardMem) Update(dt int) {
	if !d.IsVisible() {
		return
	}
	if d.showTimer.IsOn() {
		d.showTimer.Update(dt)
	}
	switch d.game.Stage() {
	case mem.Show:
		if d.showTimer.IsDone() && d.showTimer.IsOn() {
			d.showTimer.Off()
			d.Game().SetNextStage()
			d.SetupRecolection()
		}
	case mem.Conclusion:
		d.gamesData = d.gamesData.Add(d.game.GameData())
		d.Game().SetNextStage()
		d.SetupConclusion()
	}

	for _, v := range d.layout.GetContainer() {
		v.Update(dt)
	}
	d.DrawableBase.Update(dt)
}

func (d *BoardMem) Draw(surface *ebiten.Image) {
	if !d.IsVisible() {
		return
	}
	for _, v := range d.layout.GetContainer() {
		v.Draw(surface)
	}
	d.DrawableBase.Draw(surface)
}

func (d *BoardMem) Resize(rect []int) {
	d.Rect(eui.NewRect(rect))
	d.layout.Resize(rect)
	d.ImageReset()
}
