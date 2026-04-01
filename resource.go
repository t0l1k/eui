package eui

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FontDefault string = "default"
)

type ResourceManager struct {
	fonts map[string]*Font
}

func NewResourceManager() *ResourceManager { return &ResourceManager{fonts: make(map[string]*Font)} }
func (r *ResourceManager) LoadFont(name string, data []byte, size int) *ResourceManager {
	var err error
	r.fonts[name], err = NewFont(name, data, size)
	if err != nil {
		panic(err)
	}
	log.Println("ResourceManager:LoadFont:", name, size)
	return r
}
func (r *ResourceManager) GetFont(name string) *Font { return r.fonts[name] }
func (r *ResourceManager) SystemFont() *Font         { return r.fonts[FontDefault] }

func (r *ResourceManager) LoadImage(value []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(value))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}
