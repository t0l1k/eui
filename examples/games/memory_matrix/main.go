package main

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/memory_matrix/mem"
	"golang.org/x/image/colornames"
)

const (
	title = "Вспомни Матрицу"
)

type BoardMem struct {
	*eui.Container
	game                       *mem.Game
	gamesData                  *mem.GamesData
	showTimer, conclusionTimer *eui.Timer
	fn                         func(*eui.Button)
	varMsg                     *eui.Signal[string]
	varColor                   *eui.Signal[[]color.Color]
	conclusionBoard, scoreBtn  eui.Drawabler
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
	d.showTimer = eui.NewTimer(1500*time.Millisecond, func() {
		d.Game().SetNextStage()
		d.SetupRecolection()
	})
	d.conclusionBoard = eui.NewDrawable()
	d.conclusionBoard.Hide()
	d.scoreBtn = eui.NewDrawable()
	d.scoreBtn.Hide()
	d.conclusionTimer = eui.NewTimer(3000*time.Millisecond, func() {
		d.conclusionBoard.Hide()
		d.scoreBtn.Show()
	})
	d.fn = fn
	d.SetupPreparation()
	return d
}

func (d *BoardMem) Game() *mem.Game { return d.game }

func (d *BoardMem) SetupPreparation() {
	d.ResetContainer()
	d.SetLayout(eui.NewSquareGridLayout(1, 1, 1))
	btn := eui.NewButton("Click to Start "+d.game.Level().String()+" "+d.game.Dim().String(), d.fn)
	d.Add(btn)
	str := d.gamesData.String()
	d.varMsg.Emit(str)
	d.varColor.Emit([]color.Color{colornames.Yellowgreen, colornames.Black})
	if d.Rect().IsEmpty() {
		return
	}
	d.SetRect(d.Rect())
	d.Layout()
	log.Println("Setup Preparation done", d.game.String())
}

func (d *BoardMem) SetupShow() {
	d.ResetContainer()
	w, h := d.game.Dim().Width(), d.Game().Dim().Height()
	d.SetLayout(eui.NewSquareGridLayout(w, h, 1))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cell := d.Game().Cell(d.Game().Idx(x, y))
			btn := eui.NewButton(" ", d.fn)
			btn.Disable()
			if cell.IsReadOnly() {
				btn.SetBg(colornames.Orange)
			}
			d.Add(btn)
		}
	}
	d.SetRect(d.Rect())
	str := d.gamesData.String()
	d.varMsg.Emit(str)
	d.varColor.Emit([]color.Color{colornames.Red, colornames.Black})
	d.Layout()
	log.Println("Setup Show Done", d.game.String())
}

func (d *BoardMem) SetupRecolection() {
	d.ResetContainer()
	w, h := d.game.Dim().Width(), d.Game().Dim().Height()
	d.SetLayout(eui.NewSquareGridLayout(w, h, 1))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			btn := eui.NewButton(" ", d.fn)
			btn.Enable()
			d.Add(btn)
		}
	}
	d.SetRect(d.Rect())
	str := d.gamesData.String()
	d.varMsg.Emit(str)
	d.varColor.Emit([]color.Color{colornames.Blue, colornames.Yellow})
	d.Layout()
	log.Println("Setup Recolection Done", d.game.String())
}

func (d *BoardMem) SetupConclusion() {
	var str string
	d.ResetContainer()
	d.SetLayout(eui.NewGridLayout(1, 1, 1))
	sb := eui.NewSnackBar("")
	if d.Game().Win {
		str = "Winner"
		sb.SetBg(colornames.Blue)
	} else if d.Game().GameOver {
		str = "Game Over"
		sb.SetBg(colornames.Red)
	}
	sb.SetText(str + " " + d.Game().String()).ShowTime(3 * time.Second)
	d.varMsg.Emit(d.gamesData.String())
	d.varColor.Emit([]color.Color{colornames.Fuchsia, colornames.Black})
	contConclusion := eui.NewContainer(eui.NewStackLayout(5))
	d.conclusionBoard = d.setupConclusionBoard()
	d.scoreBtn = d.setupScoreBtn()
	contConclusion.Add(d.conclusionBoard)
	contConclusion.Add(d.scoreBtn)
	d.Add(contConclusion)
	d.Layout()
	d.conclusionTimer.Reset().On()
	log.Println("Setup Conclusion done", d.game.StringFull())
}

func (d *BoardMem) setupConclusionBoard() *eui.Button {
	w0, h0 := eui.GetUi().Size()
	w, h := d.game.Dim().Width(), d.Game().Dim().Height()
	cont := eui.NewContainer(eui.NewSquareGridLayout(w, h, 1))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cell := d.Game().Cell(d.Game().Idx(x, y))
			btn := eui.NewButton(" ", d.fn)
			btn.Disable()
			switch {
			case cell.IsMarked():
				btn.SetBg(colornames.Orchid)
			case cell.IsFail():
				btn.SetBg(colornames.Red)
			case !cell.IsMarked() && cell.IsReadOnly():
				btn.SetBg(colornames.Orange)
			default:
				btn.SetBg(colornames.Silver)
			}
			cont.Add(btn)
		}
	}

	cont.SetRect(d.Rect())
	r := cont.Rect()
	cameraRect := image.Rect(r.X, r.Y, r.Right(), r.Bottom())
	contentImg := ebiten.NewImage(w0, h0)
	cont.Layout()
	cont.Draw(contentImg)
	img := contentImg.SubImage(cameraRect).(*ebiten.Image)

	return eui.NewButtonIcon([]*ebiten.Image{img, img}, d.fn)
}

func (d *BoardMem) setupScoreBtn() *eui.Button {
	level := d.gamesData.Max()
	count := float64(d.gamesData.Size())
	levels := d.gamesData.Levels()
	plot := eui.NewPlot(
		func() (result []float64) {
			for i := 0.0; i < count; i++ {
				result = append(result, i+1)
			}
			return result
		}(),
		func() (result []float64) {
			for i := 0.0; i < level; i++ {
				result = append(result, i+1)
			}
			return result
		}(),
		d.gamesData.Levels(),
		"Memory Matrix",
		"Game",
		"Level").
		AddValues(func() (result []float64) {
			sum := 0.0
			for i, v := range levels {
				sum += v
				result = append(result, sum/float64(i+1))
			}
			return result
		}())
	plot.SetRect(d.Rect())
	plot.Layout()
	btn := eui.NewButtonIcon([]*ebiten.Image{plot.Image(), plot.Image()}, d.fn)
	log.Println("BoardMem:setupScoreBtn done", btn.Rect(), d.Rect(), plot.Rect(), d.gamesData.String())
	return btn
}

func (d *BoardMem) Update() {
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
	for _, v := range d.Children() {
		v.Update()
	}
	d.Container.Update()
}

type SceneMain struct{ *eui.Scene }

func NewSceneMain() *SceneMain {
	s := &SceneMain{Scene: eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 90, 5}, 5))}
	topBar := eui.NewTopBar(title, nil)
	topBar.SetUseStopwatch()
	topBar.SetShowStoppwatch(true)
	lblStatus := eui.NewLabel("")
	var board *BoardMem
	board = NewBoardMem(func(btn *eui.Button) {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
			switch board.Game().Stage() {
			case mem.Preparation:
				board.game.Move(0)
				board.showTimer.On()
				board.SetupShow()
			case mem.Restart:
				board.game.NextLevel()
				board.SetupPreparation()
			}
		}
		for i, v := range board.Children() {
			switch vv := v.(type) {
			case *eui.Button:
				if vv == btn {
					switch board.Game().Stage() {
					case mem.Recollection:
						if board.game.Move(i) {
							v.(*eui.Button).SetBg(colornames.Aqua)
						} else {
							v.(*eui.Button).SetBg(colornames.Orange)
						}
					}
				}
			}
		}
	})
	board.varMsg.Connect(func(data string) { lblStatus.SetText(data) })
	board.varColor.Connect(func(arr []color.Color) {
		lblStatus.SetBg(arr[0])
		lblStatus.SetFg(arr[1])
	})
	lblStatus.SetText(board.Game().String())
	s.Add(topBar)
	s.Add(board)
	s.Add(lblStatus)
	return s
}

func main() {
	eui.Init(func() *eui.Ui {
		u := eui.GetUi().SetTitle(title).SetSize(800, 600)
		u.Theme().Set(eui.ViewBg, colornames.Navy)
		return u
	}())
	eui.Run(NewSceneMain())
	eui.Quit(func() {})
}
