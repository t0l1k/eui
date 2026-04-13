package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/dodge_the_creeps/assets"
	"golang.org/x/image/colornames"
)

const title = "Увернись от Крипов!"

var screenWidth, screenHeight = 480, 720

type GameState int

const (
	StateMenu GameState = iota
	StateStarting
	StatePlaying
	StateGameOver
)

type Player struct {
	anim       *eui.Animation
	framesWalk []*ebiten.Image
	framesFly  []*ebiten.Image

	hit                   *eui.Signal[bool]
	lastTick              time.Time
	left, right, up, down bool
	hidden                bool
	speed                 float64
	rect                  eui.Rect[float64]
	isColliding           bool
}

func NewPlayer(frameWalk, framesFly []*ebiten.Image) *Player {
	scale := 0.5
	bounds := frameWalk[0].Bounds()
	visualWidth := float64(bounds.Dx()) * scale
	visualHeight := float64(bounds.Dy()) * scale
	a := eui.NewAnimation(frameWalk, 100*time.Millisecond)
	a.SetScale(scale).SetAngle(0).SetPos(eui.NewPoint(0.0, 0.0))
	p := &Player{
		anim:       a,
		framesWalk: frameWalk,
		framesFly:  framesFly,
		hit:        eui.NewSignal(func(a, b bool) bool { return a == b }),
		speed:      400,
		rect:       eui.NewRect([]float64{0, 0, visualWidth, visualHeight}),
	}
	p.Hide()
	return p
}
func (p *Player) Show() { p.hidden = false }
func (p *Player) Hide() { p.hidden = true }

func (p *Player) Reset() {
	p.hit.Emit(false)
	p.rect.X = float64(screenWidth)/2 - p.rect.W/2
	p.rect.Y = float64(screenHeight)/2 - p.rect.H/2
	p.anim.SetPos(eui.NewPoint(p.rect.X, p.rect.Y))
	now := time.Now()
	p.lastTick = now
	p.anim.Reset()
}

func (p *Player) Update() {
	p.isColliding = false
	now := time.Now()
	dt := now.Sub(p.lastTick).Seconds()

	if p.left || p.right { //walk
		p.anim.SetFrames(p.framesWalk)
	} else if p.up || p.down { //fly
		p.anim.SetFrames(p.framesFly)
	}

	if p.left {
		p.rect.X -= p.speed * dt
		p.anim.FlipY = false
		p.anim.FlipX = true // Зеркально для влево (по умолчанию вправо)
		p.left = false
	}
	if p.right {
		p.rect.X += p.speed * dt
		p.anim.FlipY = false
		p.anim.FlipX = false // Нормально для вправо
		p.right = false
	}
	if p.up {
		p.rect.Y -= p.speed * dt
		p.anim.FlipY = false // Нормально для вверх
		p.up = false
	}
	if p.down {
		p.rect.Y += p.speed * dt
		p.anim.FlipY = true // Зеркально для вниз
		p.down = false
	}

	if p.rect.X <= 0 {
		p.rect.X = 0
	}
	if p.rect.Y <= 0 {
		p.rect.Y = 0
	}

	if p.rect.X >= float64(screenWidth) {
		p.rect.X = float64(screenWidth)
	}
	if p.rect.Y >= float64(screenHeight) {
		p.rect.Y = float64(screenHeight)
	}

	p.anim.SetPos(eui.NewPoint(p.rect.X, p.rect.Y))
	p.anim.Update()
	p.lastTick = now
}
func (p *Player) Draw(surface *ebiten.Image) {
	if p.hidden {
		return
	}
	p.anim.Draw(surface)

	rectColor := colornames.Green
	if p.isColliding {
		rectColor = colornames.Red
	}
	vector.StrokeRect(surface, float32(p.rect.X), float32(p.rect.Y), float32(p.rect.W), float32(p.rect.H), 1, rectColor, true)
}

type Creep struct {
	*eui.Drawable
	rect        eui.Rect[float64]
	VX, VY      float64
	anim        *eui.Animation
	Hidden      bool
	isColliding bool
}

func NewCreep(creepTypes [][]*ebiten.Image) *Creep {
	// случайный набор кадров
	scale := 0.5
	frames := creepTypes[rand.Intn(len(creepTypes))]
	w, h := float64(frames[0].Bounds().Dx())*scale, float64(frames[0].Bounds().Dy())*scale

	// случайный край
	side := rand.Intn(4)
	var x, y float64

	switch side {
	case 0:
		x = rand.Float64() * float64(screenWidth)
		y = -h
	case 1:
		x = rand.Float64() * float64(screenWidth)
		y = float64(screenHeight) + h
	case 2:
		x = -w
		y = rand.Float64() * float64(screenHeight)
	case 3:
		x = float64(screenWidth) + w
		y = rand.Float64() * float64(screenHeight)
	}

	// направление на центр
	targetX := float64(screenWidth) / 2
	targetY := float64(screenHeight) / 2

	dx := targetX - x
	dy := targetY - y
	length := math.Hypot(dx, dy)

	nx := dx / length
	ny := dy / length

	// разброс направления ±30°
	spread := (rand.Float64()*60 - 30) * (math.Pi / 180)
	cos := math.Cos(spread)
	sin := math.Sin(spread)

	rx := nx*cos - ny*sin
	ry := nx*sin + ny*cos

	// случайная скорость
	speed := 1.5 + rand.Float64()*2.0

	vx := rx * speed
	vy := ry * speed

	angle := math.Atan2(vy, vx)

	a := eui.NewAnimation(frames, 100*time.Millisecond)
	a.SetAngle(angle).SetScale(scale).SetPos(eui.NewPoint(x, y))
	m := &Creep{
		rect: eui.NewRect([]float64{x, y, w, h}),
		VX:   vx,
		VY:   vy,
		anim: a,
	}
	return m
}
func (p *Creep) Update() {
	p.isColliding = false
	p.rect.X += p.VX
	p.rect.Y += p.VY
	p.anim.SetPos(eui.NewPoint(p.rect.X, p.rect.Y))
	p.anim.Update()
}
func (p *Creep) Draw(surface *ebiten.Image) {
	if p.Hidden {
		return
	}
	p.anim.Draw(surface)

	rectColor := colornames.Green
	if p.isColliding {
		rectColor = colornames.Red
	}
	vector.StrokeRect(surface, float32(p.rect.X), float32(p.rect.Y), float32(p.rect.W), float32(p.rect.H), 1, rectColor, true)
}

// GameArea — специальный UI элемент, который рисует игру внутри себя
type GameArea struct {
	*eui.Drawable
	player *Player
	mobs   map[int]*Creep
}

func NewGameArea(p *Player, mobs map[int]*Creep) *GameArea {
	ga := &GameArea{Drawable: eui.NewDrawable(), player: p, mobs: mobs}
	ga.SetViewType(eui.ViewBackground)
	return ga
}

func (g *GameArea) KeyPressed(kd eui.KeyboardData) {
	for _, v := range kd.GetKeysPressed() {
		switch v {
		case ebiten.KeyArrowLeft:
			g.player.left = true
		case ebiten.KeyArrowRight:
			g.player.right = true
		case ebiten.KeyArrowUp:
			g.player.up = true
		case ebiten.KeyArrowDown:
			g.player.down = true
		}
	}
}

func (g *GameArea) Draw(surface *ebiten.Image) {
	if g.IsHidden() {
		return
	}
	if g.IsDirty() {
		g.Layout()
	}
	// Рисуем игрока
	g.player.Draw(surface)
	// Рисуем всех крипов одним циклом без рекурсии Traverse
	if g.mobs != nil {
		for _, mob := range g.mobs {
			mob.Draw(surface)
		}
	}
}

func NewMain() *eui.Scene {
	var (
		lblStatus, lblScore              *eui.Label
		btnStart                         *eui.Button
		score                            *eui.Signal[int]
		startTimer, scoreTimer, mobTimer *eui.Timer
		state                            *eui.Signal[GameState]
		player                           *Player
		mobs                             map[int]*Creep
		mobId                            int
	)
	m := eui.NewScene(eui.NewLayoutVerticalPercent([]int{40, 20, 20, 20}, 10))
	lblScore = eui.NewLabel("0")
	lblScore.SetFontFace(assets.XoloniumRegular, assets.XoloniumRegular_ttf)
	lblScore.SetFontSize(50)
	lblScore.SetAlign(text.AlignCenter, text.AlignStart)
	lblScore.SetBg(color.Transparent)
	lblStatus = eui.NewLabel("")
	lblStatus.SetFontFace(assets.XoloniumRegular, assets.XoloniumRegular_ttf)
	lblStatus.SetBg(color.Transparent)
	btnStart = eui.NewButton("Старт", func(b *eui.Button) {
		state.Emit(StateStarting)
		log.Println("Start pressed")
	})
	btnStart.SetFontFace(assets.XoloniumRegular, assets.XoloniumRegular_ttf)
	framesWalk := []*ebiten.Image{
		eui.GetUi().RM().LoadImage(assets.PlayerGrey_walk1_png),
		eui.GetUi().RM().LoadImage(assets.PlayerGrey_walk2_png)}
	framesFly := []*ebiten.Image{
		eui.GetUi().RM().LoadImage(assets.PlayerGrey_up1_png),
		eui.GetUi().RM().LoadImage(assets.PlayerGrey_up2_png)}
	player = NewPlayer(framesWalk, framesFly)

	creepWalk := []*ebiten.Image{
		eui.GetUi().RM().LoadImage(assets.EnemyWalking_1_png),
		eui.GetUi().RM().LoadImage(assets.EnemyWalking_2_png)}
	creepFly := []*ebiten.Image{
		eui.GetUi().RM().LoadImage(assets.EnemyFlyingAlt_1_png),
		eui.GetUi().RM().LoadImage(assets.EnemyFlyingAlt_2_png)}
	creepSwim := []*ebiten.Image{
		eui.GetUi().RM().LoadImage(assets.EnemySwimming_1_png),
		eui.GetUi().RM().LoadImage(assets.EnemySwimming_2_png)}
	if mobs == nil {
		mobs = make(map[int]*Creep)
	}

	music, _, _ := eui.GetUi().RM().LoadOGG(assets.HouseInAForestLoop_ogg)
	music.Play()

	gameOverSfx, _, _ := eui.GetUi().RM().LoadWAV(assets.Gameover_wav)

	gameArea := NewGameArea(player, mobs)

	player.hit.Connect(func(data bool) {
		if data {
			state.Emit(StateGameOver)
			log.Println("Столкнулись с крипом! Конец!", len(mobs))
		}
	})
	score = eui.NewSignal(func(a, b int) bool { return a == b })
	score.Connect(func(data int) {
		lblScore.SetText(fmt.Sprintf("%v", data))
	})
	state = eui.NewSignal(func(a, b GameState) bool { return a == b })
	state.ConnectAndFire(func(data GameState) {
		switch data {
		case StateMenu:
			lblStatus.SetText(title)
			lblStatus.Show()
			lblScore.Show()
			btnStart.Show()
			score.Emit(0)
			log.Println("StateMenu")
		case StateStarting:
			for k := range mobs {
				delete(mobs, k)
			}
			mobId = 0
			player.Reset()
			startTimer.On()
			lblStatus.SetText("Приготовиться")
			btnStart.Hide()
			player.Show()
			music.Rewind()
			music.Play()
			log.Println("StateStarting")
		case StatePlaying:
			lblStatus.Hide()
			scoreTimer.On()
			mobTimer.On()
			log.Println("StatePlaying")
		case StateGameOver:
			if music.IsPlaying() {
				music.Pause()
			}
			gameOverSfx.Rewind()
			gameOverSfx.Play()
			player.Hide()
			scoreTimer.Off()
			mobTimer.Off()
			startTimer.On()
			lblStatus.SetText("Конец!").Show()
			log.Println("StateGameOver")
		}
	}, StateMenu)

	mobTimer = eui.NewTimer(500*time.Millisecond, func() {
		if state.Value() == StatePlaying {
			mob := NewCreep([][]*ebiten.Image{creepFly, creepSwim, creepWalk})
			mobs[mobId] = mob
			mobId++
			mobTimer.On()
			log.Println("Новый крип", mobId, len(mobs))
		}
	})
	scoreTimer = eui.NewTimer(1*time.Second, func() {
		score.Emit(score.Value() + 1)
		scoreTimer.On()
	})
	startTimer = eui.NewTimer(2*time.Second, func() {
		if state.Value() == StateStarting {
			state.Emit(StatePlaying)
			log.Println("Начало игры")
		}
		if state.Value() == StateGameOver {
			state.Emit(StateMenu)
			log.Println("Выбор меню")
		}
	})
	eui.GetUi().TickListener().Connect(func(data eui.Event) {
		if state.Value() != StateMenu {
			screenRect := eui.NewRect([]float64{0, 0, float64(screenWidth), float64(screenHeight)})
			player.Update()
			for id, mob := range mobs {
				mob.Update()
				if player.rect.Intersects(mob.rect) {
					player.isColliding = true
					mob.isColliding = true
					player.hit.Emit(true)
					break
				}

				// Если крип полностью покинул экран (с запасом в его ширину), удаляем его
				if mob.rect.IsOutside(screenRect.Inflated(mob.rect.W)) {
					delete(mobs, id)
					log.Println("Крип вышел за экран", id, mob.rect.X, mob.rect.Y, len(mobs), mobId)
				}
			}
			if state.Value() == StatePlaying {
				if !music.IsPlaying() {
					music.Rewind()
					music.Play()
				}
			}
		}
	})
	eui.GetUi().KeyboardListener().Connect(func(data eui.Event) {
		kd := data.Value.(eui.KeyboardData)
		if kd.IsReleased(ebiten.KeySpace) && state.Value() == StateMenu {
			state.Emit(StateStarting)
		}
	})
	eui.GetUi().ResizeListener().Connect(func(data eui.Event) {
		sz := data.Value.(eui.Rect[int])
		screenWidth = sz.Width()
		screenHeight = sz.Height()
	})
	m.SetBg(colornames.Teal)
	m.Add(lblScore)
	m.Add(gameArea)
	m.Add(lblStatus)
	m.Add(btnStart)
	m.Add(eui.NewDrawable())

	theme := eui.GetUi().Theme()
	theme.Set(eui.SceneBg, colornames.Teal)
	return m
}

func main() {
	eui.Init(eui.GetUi().SetTitle(title).SetSize(screenWidth, screenHeight))
	eui.Run(NewMain())
}
