package eui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Умею показать кнопку под мышкой выделенной, или нажатой, или отпущенной(после отпускания исполняется прикрепленный метод)
type Button struct {
	*View
	text                *Text
	onPressed           func(*Button)
	buttonPressed       bool
	left, right, middle bool
	margin              int
}

func NewButton(text string, f func(*Button)) *Button {
	b := &Button{View: NewView(), onPressed: f}
	b.text = NewText(text)
	b.SetupButton(text, f)
	b.Visible(true)
	return b
}

func (b *Button) SetupButton(text string, f func(*Button)) {
	b.onPressed = f
	b.text.SetText(text)
	theme := GetUi().theme
	bg := theme.Get(ButtonBg)
	fg := theme.Get(ButtonFg)
	b.Bg(bg)
	b.Fg(fg)
}

func (b *Button) UpdateData(value interface{}) {
	switch v := value.(type) {
	case string:
		b.text.SetText(v)
	case int:
		b.text.SetText(strconv.Itoa(v))
	}
}

func (b *Button) GetText() string { return b.text.GetText() }
func (b *Button) SetText(value string) {
	if b.text.GetText() == value {
		return
	}
	b.text.SetText(value)
	b.MarkDirty()
}

func (b *Button) Bg(bg color.Color) {
	if b.bg == bg {
		return
	}
	b.bg = bg
	b.text.Bg(bg)
	b.MarkDirty()
}

func (b *Button) Fg(fg color.Color) {
	if b.fg == fg {
		return
	}
	b.fg = fg
	b.text.Fg(fg)
	b.MarkDirty()
}

func (b *Button) SetFunc(f func(*Button)) {
	b.onPressed = f
}

func (b *Button) IsMouseDownLeft() bool {
	return b.left && b.buttonPressed || b.buttonPressed && b.state == ViewStateExec
}

func (b *Button) IsMouseDownRight() bool {
	return b.right
}

func (b *Button) IsMouseDownMiddle() bool {
	return b.middle
}

func (b *Button) Layout() {
	b.View.Layout()
	var fg color.Color
	theme := GetUi().theme
	switch b.state {
	case ViewStateHover:
		fg = theme.Get(ButtonHover)
	case ViewStateFocus:
		fg = theme.Get(ButtonFocus)
	case ViewStateNormal:
		fg = theme.Get(ButtonNormal)
	case ViewStateSelected:
		fg = theme.Get(ButtonSelected)
	case ViewStateDisabled:
		fg = theme.Get(ButtonDisabled)
	case ViewStateActive:
		fg = theme.Get(ButtonActive)
	}
	_, _, w, h := b.Rect().GetRectFloat()
	bold := b.margin
	if b.buttonPressed {
		bold = b.margin * 2
	}
	vector.StrokeRect(b.Image(), 0, 0, w, h, float32(bold), fg, true)
	b.ClearDirty()
}

func (b *Button) IsPressed() bool { return b.buttonPressed }

func (b *Button) Pressed(value bool) {
	b.buttonPressed = value
	if value {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.left = true
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			b.right = true
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			b.middle = true
		}
	} else {
		b.left = false
		b.right = false
		b.middle = false
	}
}

func (b *Button) Update(dt int) {
	if b.IsDisabled() {
		return
	}
	if b.state == ViewStateFocus && !b.buttonPressed {
		b.Pressed(true)
	}
	if (b.state == ViewStateHover || b.state == ViewStateExec) && b.buttonPressed {
		if b.onPressed != nil {
			b.onPressed(b)
		}
		b.Pressed(false)
		if b.state == ViewStateExec {
			b.state = ViewStateNormal
		}
	}
	if b.state == ViewStateNormal {
		b.Pressed(false)
	}
}

func (b *Button) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	if b.IsDirty() {
		b.Layout()
		b.text.Layout()
	}
	b.text.Draw(surface)
	op := &ebiten.DrawImageOptions{}
	x, y := b.Rect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(b.Image(), op)
}

func (b *Button) Resize(rect Rect) {
	b.View.Resize(rect)
	b.margin = int(float64(b.Rect().GetLowestSize()) * 0.03)
	x, y, w, h := b.Rect().GetRect()
	b.text.Resize(NewRect([]int{x + b.margin, y + b.margin, w - b.margin*2, h - b.margin*2}))
	b.ImageReset()
}
