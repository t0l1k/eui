package eui

import (
	"fmt"
)

type Point struct {
	X, Y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

func (p Point) Get() (float64, float64) {
	return p.X, p.Y
}

func (p Point) GetX() float64 {
	return p.X
}

func (p Point) GetY() float64 {
	return p.Y
}

func (p Point) String() string {
	return fmt.Sprintf("[%.2f, %.2f]", p.X, p.Y)
}

type PointInt struct {
	X, Y int
}

func NewPointInt(x, y int) *PointInt {
	return &PointInt{
		X: x,
		Y: y,
	}
}

func (p PointInt) Get() (int, int) {
	return p.X, p.Y
}

func (p PointInt) GetX() int {
	return p.X
}

func (p PointInt) GetY() int {
	return p.Y
}

func (p PointInt) Equal(a PointInt) bool {
	return p.X == a.X && p.Y == a.Y
}

func (p PointInt) Offset(a PointInt) *PointInt {
	return &PointInt{p.X - a.X, p.Y - a.Y}
}

func (p PointInt) String() string {
	return fmt.Sprintf("[%.2d, %.2d]", p.X, p.Y)
}

type Rect struct {
	X, Y, W, H int
}

func NewRect(arr []int) *Rect {
	return &Rect{
		X: arr[0],
		Y: arr[1],
		W: arr[2],
		H: arr[3],
	}
}
func (r Rect) InRect(x, y int) bool {
	return r.Left() <= x && r.Right() >= x && r.Top() <= y && r.Bottom() >= y
}

func (r Rect) Pos() (int, int) {
	return r.X, r.Y
}

func (r Rect) Size() (int, int) {
	return r.W, r.H
}

func (r Rect) GetArr() []int {
	return []int{r.X, r.Y, r.W, r.H}
}

func (r Rect) GetRect() (int, int, int, int) {
	return r.X, r.Y, r.W, r.H
}

func (r Rect) GetRectFloat() (float32, float32, float32, float32) {
	return float32(r.X), float32(r.Y), float32(r.W), float32(r.H)
}

func (r Rect) GetRectFloat64() (float64, float64, float64, float64) {
	return float64(r.X), float64(r.Y), float64(r.W), float64(r.H)
}

func (r Rect) Left() int {
	return r.X
}

func (r Rect) Right() int {
	return r.X + r.W
}

func (r Rect) Top() int {
	return r.Y
}

func (r Rect) Bottom() int {
	return r.Y + r.H
}
func (r Rect) CenterX() int {
	return (r.Right() - r.X) / 2
}

func (r Rect) CenterY() int {
	return (r.Bottom() - r.Y) / 2
}
func (r Rect) Center() (int, int) {
	return r.CenterX(), r.CenterY()
}

func (r Rect) TopLeft() (int, int) {
	return r.X, r.Y
}

func (r Rect) TopRight() (int, int) {
	return r.X + r.W, r.Y
}

func (r Rect) BottomLeft() (int, int) {
	return r.X, r.Y + r.H
}

func (r Rect) BottomRight() (int, int) {
	return r.X + r.W, r.Y + r.H
}

func (r Rect) GetLowestSize() int {
	result := r.W
	if r.W > r.H {
		result = r.H
	}
	return result
}

func (r Rect) String() string {
	return fmt.Sprintf("[%v, %v, %v, %v]", r.X, r.Y, r.W, r.H)
}
