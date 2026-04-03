package eui

import (
	"bytes"
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
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

func (r *ResourceManager) LoadOGG(value []byte) (*audio.Player, time.Duration, error) {
	if audio.CurrentContext() == nil {
		audio.NewContext(48000)
	}
	d, err := vorbis.DecodeWithSampleRate(audio.CurrentContext().SampleRate(), bytes.NewReader(value))
	if err != nil {
		return nil, 0, err
	}
	p, err := audio.CurrentContext().NewPlayer(d)
	if err != nil {
		return nil, 0, err
	}
	// Конвертация байтов в Duration: bytes / (channels * bytesPerSample * sampleRate)
	// Ebiten использует 16-bit (2 bytes) stereo (2 channels) = 4 bytes per frame
	duration := time.Duration(d.Length()) * time.Second / (time.Duration(audio.CurrentContext().SampleRate()) * 4)
	return p, duration, nil
}

func (r *ResourceManager) LoadWAV(value []byte) (*audio.Player, time.Duration, error) {
	if audio.CurrentContext() == nil {
		audio.NewContext(48000)
	}
	d, err := wav.DecodeWithSampleRate(audio.CurrentContext().SampleRate(), bytes.NewReader(value))
	if err != nil {
		return nil, 0, err
	}
	p, err := audio.CurrentContext().NewPlayer(d)
	if err != nil {
		return nil, 0, err
	}
	duration := time.Duration(d.Length()) * time.Second / (time.Duration(audio.CurrentContext().SampleRate()) * 4)
	return p, duration, nil
}
