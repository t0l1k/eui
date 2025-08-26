package res

import (
	"bytes"
	"image"

	_ "embed"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed mines.png
	MinesSprites_png                                               []byte
	AllMinesSprite                                                 *ebiten.Image
	SmileSprites, NumbersSprites, CellsUpSprites, CellsDownSprites []*ebiten.Image
)

func init() {
	AllMinesSprite = GetMinesPng()
	SmileSprites = LoadSmiles(AllMinesSprite)
	CellsUpSprites = LoadCellIUpcons(AllMinesSprite)
	CellsDownSprites = LoadCellIDowncons(AllMinesSprite)
	NumbersSprites = LoadNumbers(AllMinesSprite)
}

func GetMinesPng() *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(MinesSprites_png))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

func LoadSmiles(img *ebiten.Image) []*ebiten.Image {
	var icons []*ebiten.Image // x:0 y:24 w: 26 h: 26
	var nImg *ebiten.Image
	w, h := 26, 26
	for i := 0; i < 5; i++ {
		nImg = ebiten.NewImage(w, h)
		x, y := i*w, 24
		x += i
		nImg.DrawImage(img.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image), &ebiten.DrawImageOptions{})
		icons = append(icons, nImg)
	}
	return icons
}

func LoadCellIUpcons(img *ebiten.Image) []*ebiten.Image {
	var icons []*ebiten.Image // x:0 y:51 w: 16 h: 16
	var nImg *ebiten.Image
	w, h := 16, 16
	for i := 0; i < 8; i++ {
		nImg = ebiten.NewImage(w, h)
		x, y := i*w, 51
		x += i
		nImg.DrawImage(img.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image), &ebiten.DrawImageOptions{})
		icons = append(icons, nImg)
	}
	return icons
}

func LoadCellIDowncons(img *ebiten.Image) []*ebiten.Image {
	var icons []*ebiten.Image // x:0 y:68 w: 16 h: 16
	var nImg *ebiten.Image
	w, h := 16, 16
	for i := 0; i < 8; i++ {
		nImg = ebiten.NewImage(w, h)
		x, y := i*w, 68
		x += i
		nImg.DrawImage(img.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image), &ebiten.DrawImageOptions{})
		icons = append(icons, nImg)
	}
	return icons
}

func LoadNumbers(img *ebiten.Image) []*ebiten.Image {
	var numbers []*ebiten.Image // x:0 y:0 w: 13 h: 23
	var nImg *ebiten.Image
	w, h := 14, 24
	nImg = ebiten.NewImage(w, h)
	for i := 0; i < 10; i++ {
		nImg.Clear()
		op := &ebiten.DrawImageOptions{}
		// op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		x, y := i*w, 0
		nImg.DrawImage(img.SubImage(image.Rect(x, y, w, h)).(*ebiten.Image), op)
		numbers = append(numbers, nImg)
	}
	return numbers
}
