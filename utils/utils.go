package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
)

func DrawDebugLines(surface *ebiten.Image, rect *eui.Rect) {
	x0, y0 := rect.Pos()
	w0, h0 := rect.BottomRight()
	x, y := float32(x0), float32(y0)
	w, h := float32(w0), float32(h0)
	vector.StrokeLine(surface, x, y, w, h, 2, colors.Red, true)
	vector.StrokeLine(surface, x, h, w, y, 1, colors.Red, true)
	vector.StrokeLine(surface, x, h/2, w, h/2, 1, colors.Red, true)
	vector.StrokeLine(surface, w/2, h, w/2, y, 1, colors.Red, true)
}

func RemoveFromIntSliceValue(arr []int, value int) []int {
	if !IntSliceContains(arr, value) {
		return arr
	}
	idx := GetIdxValueFromIntSlice(arr, value)
	copy(arr[idx:], arr[idx+1:])
	arr[len(arr)-1] = value
	arr = arr[:len(arr)-1]
	return arr
}

func IntSliceContains(arr []int, value int) bool {
	for _, v := range arr {
		if value == v {
			return true
		}
	}
	return false
}

func GetIdxValueFromIntSlice(arr []int, value int) int {
	for i, v := range arr {
		if value == v {
			return i
		}
	}
	return -1
}

func IntSlicesIsEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func PopIntSlice(a []int) []int {
	a[len(a)-1] = -1
	a = a[:len(a)-1]
	return a
}
