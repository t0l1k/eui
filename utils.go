package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawDebugLines(surface *ebiten.Image, rect *Rect) {
	x0, y0 := rect.Pos()
	w0, h0 := rect.BottomRight()
	x, y := float32(x0), float32(y0)
	w, h := float32(w0), float32(h0)
	vector.StrokeLine(surface, x, y, w, h, 2, Red, true)
	vector.StrokeLine(surface, x, h, w, y, 1, Red, true)
	vector.StrokeLine(surface, x, h/2, w, h/2, 1, Red, true)
	vector.StrokeLine(surface, w/2, h, w/2, y, 1, Red, true)
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
