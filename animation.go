package eui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Animation управляет последовательностью кадров и их отображением.
type Animation struct {
	*Drawable
	frames  []*ebiten.Image
	index   int
	tick    time.Duration
	last    time.Time
	playing bool
	FlipX   bool
	FlipY   bool

	angle, scale float64
	pos          Point[float64]
}

// NewAnimation создает новую анимацию с заданным интервалом смены кадров.
func NewAnimation(frames []*ebiten.Image, tick time.Duration) *Animation {
	a := &Animation{
		Drawable: NewDrawable(),
		frames:   frames,
		tick:     tick,
		playing:  true,
		scale:    1.0,
		last:     time.Now(),
	}
	a.SetViewType(ViewBackground)
	return a
}

// SetFrames меняет текущий набор кадров. Если набор новый, индекс сбрасывается.
func (a *Animation) SetFrames(frames []*ebiten.Image) {
	if len(a.frames) > 0 && len(frames) > 0 && &a.frames[0] == &frames[0] {
		return
	}
	a.frames = frames
	a.index = 0
}

// Update обновляет состояние анимации на основе прошедшего времени.
func (a *Animation) Update() {
	if !a.playing || len(a.frames) <= 1 {
		return
	}
	if time.Since(a.last) > a.tick {
		a.index = (a.index + 1) % len(a.frames)
		a.last = time.Now()
	}
}

// Play запускает проигрывание.
func (a *Animation) Play() { a.playing = true }

// Stop останавливает проигрывание.
func (a *Animation) Stop() { a.playing = false }

// Reset сбрасывает анимацию к началу.
func (a *Animation) Reset() {
	a.index = 0
	a.last = time.Now()
}

func (a *Animation) SetAngle(value float64) *Animation      { a.angle = value; return a }
func (a *Animation) SetScale(value float64) *Animation      { a.scale = value; return a }
func (a *Animation) SetPos(value Point[float64]) *Animation { a.pos = value; return a }

// Draw отрисовывает текущий кадр с учетом зеркалирования, масштаба и поворота.
func (a *Animation) Draw(surface *ebiten.Image) {
	if len(a.frames) == 0 {
		return
	}
	img := a.frames[a.index]
	w, h := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	sx, sy := a.scale, a.scale
	if a.FlipX {
		sx = -sx
	}
	if a.FlipY {
		sy = -sy
	}
	op.GeoM.Scale(sx, sy)
	if a.angle != 0 {
		op.GeoM.Rotate(a.angle)
	}
	x, y := a.pos.Get()
	op.GeoM.Translate(x+(w*a.scale)/2, y+(h*a.scale)/2)
	surface.DrawImage(img, op)
}
