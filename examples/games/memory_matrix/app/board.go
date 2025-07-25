package app

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/memory_matrix/mem"
	"golang.org/x/image/colornames"
)

type BoardMem struct {
	*eui.Container
	game      *mem.Game
	gamesData *mem.GamesData
	showTimer *eui.Timer
	fn        func(*eui.Button)
	varMsg    *eui.Signal[string]
	varColor  *eui.Signal[[]color.Color]
}

func NewBoardMem(fn func(*eui.Button)) *BoardMem {
	d := &BoardMem{Container: eui.NewContainer(eui.NewGridLayout(1, 1, 1))}
	d.varMsg = eui.NewSignal(func(a, b string) bool { return a == b })
	d.varColor = eui.NewSignal(func(a, b []color.Color) bool {
		// Если оба nil — считаем равными
		if a == nil && b == nil {
			return true
		}
		// Если только один nil — не равны
		if a == nil || b == nil {
			return false
		}
		// Приведение к []color.Color
		aArr := a
		bArr := b
		if len(aArr) != len(bArr) {
			return false
		}
		for i := range aArr {
			aR, aG, aB, aA := aArr[i].RGBA()
			bR, bG, bB, bA := bArr[i].RGBA()
			if aR != bR || aG != bG || aB != bB || aA != bA {
				return false
			}
		}
		return true
	})
	d.game = mem.NewGame(mem.Level(1))
	d.gamesData = mem.NewGamesData()
	// d.SetHidden(true)
	d.showTimer = eui.NewTimer(1500*time.Millisecond, func() {
		d.Game().SetNextStage()
		d.SetupRecolection()
	})
	d.fn = fn
	d.SetupPreparation()
	return d
}

func (d *BoardMem) Game() *mem.Game { return d.game }

func (d *BoardMem) SetupPreparation() {
	d.ResetContainer()
	d.SetLayout(eui.NewGridLayout(1, 1, 1))
	btn := eui.NewButton("Click to Start "+d.game.Level().String()+" "+d.game.Dim().String(), d.fn)
	d.Add(btn)
	d.Resize(d.Rect())
	str := d.gamesData.String()
	d.varMsg.Emit(str)
	d.varColor.Emit([]color.Color{colornames.Yellowgreen, colornames.Black})
	log.Println("Setup Preparation done", d.game.String())
}

func (d *BoardMem) SetupShow() {
	d.ResetContainer()
	w, h := d.game.Dim().Width(), d.Game().Dim().Height()
	d.SetLayout(eui.NewGridLayout(float64(w), float64(h), 1))
	for y := 0; y < d.Game().Dim().Height(); y++ {
		for x := 0; x < d.game.Dim().Width(); x++ {
			cell := d.Game().Cell(d.Game().Idx(x, y))
			btn := eui.NewButton(" ", d.fn)
			btn.Disable()
			if cell.IsReadOnly() {
				btn.Bg(colornames.Orange)
			}
			d.Add(btn)
		}
	}
	d.Resize(d.Rect())
	str := d.gamesData.String()
	d.varMsg.Emit(str)
	d.varColor.Emit([]color.Color{colornames.Red, colornames.Black})
	log.Println("Setup Show Done", d.game.String())
}

func (d *BoardMem) SetupRecolection() {
	d.ResetContainer()
	w, h := d.game.Dim().Width(), d.Game().Dim().Height()
	d.SetLayout(eui.NewGridLayout(float64(w), float64(h), 1))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			btn := eui.NewButton(" ", d.fn)
			btn.Enable()
			d.Add(btn)
		}
	}
	d.Resize(d.Rect())
	str := d.gamesData.String()
	d.varMsg.Emit(str)
	d.varColor.Emit([]color.Color{colornames.Blue, colornames.Yellow})
	log.Println("Setup Recolection Done", d.game.String())
}

func (d *BoardMem) SetupConclusion() {
	var str string
	d.ResetContainer()
	d.SetLayout(eui.NewSquareGridLayout(1, 1, 1))
	sb := eui.NewSnackBar("")
	if d.Game().Win {
		str = "Winner"
		sb.Bg(colornames.Blue)
	} else if d.Game().GameOver {
		str = "Game Over"
		sb.Bg(colornames.Red)
	}
	sb.SetText(str + " " + d.Game().String()).Show(3 * time.Second)
	d.varMsg.Emit(d.gamesData.String())
	d.varColor.Emit([]color.Color{colornames.Fuchsia, colornames.Black})
	btn := d.setupScoreBtn()
	d.Add(btn)
	d.Resize(d.Rect())
	log.Println("Setup Conclusion done", d.game.String(), d.Rect(), btn.Rect())
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
	log.Println("d.Rect()=", d.Rect())
	plot.Resize(d.Rect())
	log.Println("plot.Rect()=", plot.Rect())
	plot.Layout()
	log.Println("plot.Image().Bounds()=", plot.Image().Bounds())
	btn := eui.NewButtonIcon([]*ebiten.Image{plot.Image(), plot.Image()}, d.fn)
	log.Println("BoardMem:setupScoreBtn done", btn.Rect(), d.Rect(), plot.Rect())
	return btn
}

func (d *BoardMem) Update(dt int) {
	if d.IsHidden() {
		return
	}
	switch d.game.Stage() {
	case mem.Show:
	case mem.Conclusion:
		d.gamesData = d.gamesData.Add(d.game.GameData())
		d.Game().SetNextStage()
		d.SetupConclusion()
	}
	for _, v := range d.Childrens() {
		v.Update(dt)
	}
	d.Container.Update(dt)
}
