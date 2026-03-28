package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
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
	hit                   *eui.Signal[bool]
	speed                 float64
	x, y                  float64
	width                 float64
	color                 color.Color
	last                  time.Time
	left, right, up, down bool
	hidden                bool
}

func NewPlayer() *Player {
	p := &Player{
		hit:   eui.NewSignal(func(a, b bool) bool { return a == b }),
		width: 32,
		color: colornames.Aqua,
		speed: 400,
	}
	p.Hide()
	return p
}
func (p *Player) Show() { p.hidden = false }
func (p *Player) Hide() { p.hidden = true }

func (p *Player) Reset() {
	p.hit.Emit(false)
	p.x = float64(screenWidth)/2 - 16
	p.y = float64(screenHeight)/2 - 16
	p.last = time.Now()
}

func (p *Player) Update() {
	now := time.Now()
	dt := now.Sub(p.last).Seconds()
	p.last = now
	if p.left {
		p.x -= p.speed * dt
		p.left = false
	}
	if p.right {
		p.x += p.speed * dt
		p.right = false
	}
	if p.up {
		p.y -= p.speed * dt
		p.up = false
	}
	if p.down {
		p.y += p.speed * dt
		p.down = false
	}

	if p.x <= 0 {
		p.x = 0
	}
	if p.y <= 0 {
		p.y = 0
	}

	if p.x >= float64(screenWidth) {
		p.x = float64(screenWidth)
	}
	if p.y >= float64(screenHeight) {
		p.y = float64(screenHeight)
	}
}
func (p *Player) Draw(surface *ebiten.Image) {
	if p.hidden {
		return
	}
	x := p.x - p.width/2
	y := p.y - p.width/2
	// ebitenutil.DrawRect(surface, x, y, p.width, p.width, p.color)
	vector.FillRect(surface, float32(x), float32(y), float32(p.width), float32(p.width), p.color, true)
}

type Mob struct {
	*eui.Drawable
	X, Y   float64
	VX, VY float64
	Frames []*ebiten.Image
	Frame  int
	Size   float64
	Angle  float64
	Hidden bool
	color  color.Color
}

func NewMob() *Mob {
	size := 32.0

	// случайный набор кадров
	// frames := creepTypes[rand.Intn(len(creepTypes))]

	// случайный край
	side := rand.Intn(4)
	var x, y float64

	switch side {
	case 0:
		x = rand.Float64() * float64(screenWidth)
		y = -size
	case 1:
		x = rand.Float64() * float64(screenWidth)
		y = float64(screenHeight) + size
	case 2:
		x = -size
		y = rand.Float64() * float64(screenHeight)
	case 3:
		x = float64(screenWidth) + size
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

	m := &Mob{
		X:  x,
		Y:  y,
		VX: vx,
		VY: vy,
		// Frames: frames,
		Frame: 0,
		Size:  size,
		Angle: angle,
		color: colornames.Red,
	}
	return m
}
func (p *Mob) Update() {
	p.X += p.VX
	p.Y += p.VY
}
func (p *Mob) Draw(surface *ebiten.Image) {
	if p.Hidden {
		return
	}
	x := p.X - 32/2
	y := p.Y - 32/2
	vector.FillRect(surface, float32(x), float32(y), float32(32), float32(32), p.color, true)
}

// GameArea — специальный UI элемент, который рисует игру внутри себя
type GameArea struct {
	*eui.Drawable
	player *Player
	mobs   map[int]*Mob
}

func NewGameArea(p *Player, mobs map[int]*Mob) *GameArea {
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

func isOut(mob *Mob, sz float64) bool {
	return mob.X < -sz || mob.X > float64(screenWidth)+sz || mob.Y < -sz || mob.Y > float64(screenHeight)+sz
}

func intersects(ax, ay, as, bx, by, bs float64) bool {
	return ax < bx+bs &&
		ax+as > bx &&
		ay < by+bs &&
		ay+as > by
}

func NewMain() *eui.Scene {
	var (
		lblStatus, lblScore              *eui.Label
		btnStart                         *eui.Button
		score                            *eui.Signal[int]
		startTimer, scoreTimer, mobTimer *eui.Timer
		state                            *eui.Signal[GameState]
		player                           *Player
		mobs                             map[int]*Mob
		mobId                            int
	)
	m := eui.NewScene(eui.NewLayoutVerticalPercent([]int{40, 20, 20, 20}, 10))
	lblScore = eui.NewLabel("0")
	lblScore.SetFontSize(50)
	lblScore.SetAlign(eui.LabelAlignUp)
	lblScore.SetBg(color.Transparent)
	lblStatus = eui.NewLabel("")
	lblStatus.SetBg(color.Transparent)
	btnStart = eui.NewButton("Старт", func(b *eui.Button) {
		state.Emit(StateStarting)
		log.Println("Start pressed")
	})
	player = NewPlayer()
	if mobs == nil {
		mobs = make(map[int]*Mob)
	}

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
			lblStatus.SetText(title).SetFontSize(50)
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
			lblStatus.SetText("Приготовиться").SetFontSize(40)
			btnStart.Hide()
			player.Show()
			log.Println("StateStarting")
		case StatePlaying:
			lblStatus.Hide()
			scoreTimer.On()
			mobTimer.On()
			log.Println("StatePlaying")
		case StateGameOver:
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
			mob := NewMob()
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
			player.Update()
			for id, mob := range mobs {
				mob.Update()
				if intersects(player.x, player.y, player.width, mob.X, mob.Y, mob.Size) {
					player.hit.Emit(true)
					break
				}
				sz := mob.Size
				if isOut(mob, sz) {
					delete(mobs, id)
					log.Println("Крип вышел за экран", id, mob.X, mob.Y, len(mobs), mobId)
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
	m.Add(gameArea) // Теперь это единственный игровой объект в дереве UI
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
