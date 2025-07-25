package eui

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Numbers interface {
	constraints.Integer | constraints.Float
}

type Point[T Numbers] struct{ X, Y T }

func NewPoint[T Numbers](x, y T) Point[T]      { return Point[T]{X: x, Y: y} }
func (p Point[T]) Get() (T, T)                 { return p.X, p.Y }
func (p Point[T]) GetX() T                     { return p.X }
func (p Point[T]) GetY() T                     { return p.Y }
func (p Point[T]) Equal(value Point[T]) bool   { return p.X == value.X && p.Y == value.Y }
func (p Point[T]) Offset(a Point[T]) *Point[T] { return &Point[T]{p.X - a.X, p.Y - a.Y} }
func (p Point[T]) String() string              { return fmt.Sprintf("[%.2v, %.2v]", p.X, p.Y) }

type Rect[T Numbers] struct{ X, Y, W, H T }

func NewRect[T Numbers](arr []T) Rect[T] { return Rect[T]{X: arr[0], Y: arr[1], W: arr[2], H: arr[3]} }
func (r Rect[T]) InRect(x, y T) bool {
	return r.Left() <= x && r.Right() >= x && r.Top() <= y && r.Bottom() >= y
}
func (r Rect[T]) Pos() (T, T)           { return r.X, r.Y }
func (r Rect[T]) Size() (T, T)          { return r.W, r.H }
func (r Rect[T]) GetArr() []T           { return []T{r.X, r.Y, r.W, r.H} }
func (r Rect[T]) GetRect() (T, T, T, T) { return r.X, r.Y, r.W, r.H }
func (r Rect[T]) GetRectFloat() (float32, float32, float32, float32) {
	return float32(r.X), float32(r.Y), float32(r.W), float32(r.H)
}
func (r Rect[T]) GetRectFloat64() (float64, float64, float64, float64) {
	return float64(r.X), float64(r.Y), float64(r.W), float64(r.H)
}
func (r Rect[T]) Left() T              { return r.X }
func (r Rect[T]) Right() T             { return r.X + r.W }
func (r Rect[T]) Top() T               { return r.Y }
func (r Rect[T]) Bottom() T            { return r.Y + r.H }
func (r Rect[T]) CenterX() T           { return (r.Right() - r.X) / 2 }
func (r Rect[T]) CenterY() T           { return (r.Bottom() - r.Y) / 2 }
func (r Rect[T]) Center() (T, T)       { return r.CenterX(), r.CenterY() }
func (r Rect[T]) TopLeft() (T, T)      { return r.X, r.Y }
func (r Rect[T]) TopRight() (T, T)     { return r.X + r.W, r.Y }
func (r Rect[T]) BottomLeft() (T, T)   { return r.X, r.Y + r.H }
func (r Rect[T]) BottomRight() (T, T)  { return r.X + r.W, r.Y + r.H }
func (r Rect[T]) Width() T             { return r.W }
func (r Rect[T]) Height() T            { return r.H }
func (r Rect[T]) GetLowestSize() T     { return min(r.W, r.H) }
func (a Rect[T]) Equal(b Rect[T]) bool { return a.X == b.X && a.Y == b.Y && a.W == b.W && a.H == b.H }
func (a Rect[T]) IsEmpty() bool        { return a.W == 0 && a.H == 0 }
func (r Rect[T]) String() string       { return fmt.Sprintf("[%v, %v, %v, %v]", r.X, r.Y, r.W, r.H) }
