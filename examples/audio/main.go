package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	riaudio "github.com/hajimehoshi/ebiten/v2/examples/resources/images/audio"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/dodge_the_creeps/assets"
)

func main() {
	eui.Init(eui.GetUi().SetTitle("Audio Player").SetSize(400, 600))
	eui.Run(func() *eui.Scene {
		s := eui.NewScene(eui.NewVBoxLayout(10))
		playerLoop, totalDuration, _ := eui.GetUi().RM().LoadOGG(assets.HouseInAForestLoop_ogg)
		playerSfx, _, _ := eui.GetUi().RM().LoadWAV(assets.Gameover_wav)

		progressBar := eui.NewProgress(0, 1, 0, eui.Horizontal, func(data float64) {})
		lblTime := eui.NewLabel("00:00 / 00:00")

		sliderSeek := eui.NewSlider(0, 1, 0, eui.Horizontal, func(data float64) {
			playerLoop.SetPosition(time.Duration(data * float64(totalDuration)))
		})

		lblVolume := eui.NewLabel("")
		sliderVolume := eui.NewSlider(0, 1, playerLoop.Volume(), eui.Horizontal, func(data float64) {
			playerLoop.SetVolume(data)
			lblVolume.SetText(fmt.Sprintf("Volume:%.2v", playerLoop.Volume()))
		})

		var (
			btnPlay, btnSfx *eui.Button
		)

		playIcon := eui.GetUi().RM().LoadImage(riaudio.Play_png)
		pauseIcon := eui.GetUi().RM().LoadImage(riaudio.Pause_png)
		alertIcon := eui.GetUi().RM().LoadImage(riaudio.Alert_png)

		play := false
		loop := false

		playMusic := func() {
			if play {
				playerLoop.Pause()
				btnPlay.SetIcons([]*ebiten.Image{playIcon, playIcon})
			} else if !play {
				playerLoop.Play()
				btnPlay.SetIcons([]*ebiten.Image{pauseIcon, pauseIcon})
			}
			play = !play
		}

		playSfx := func() {
			playerSfx.Rewind()
			playerSfx.Play()
		}

		btnPlay = eui.NewButtonIcon([]*ebiten.Image{playIcon, playIcon}, func(b *eui.Button) {
			playMusic()
		})
		btnSfx = eui.NewButtonIcon([]*ebiten.Image{alertIcon, alertIcon}, func(b *eui.Button) {
			playSfx()
		})
		checkbox := eui.NewCheckbox("Loop", func(c *eui.Checkbox) {
			loop = c.IsChecked()
		})
		checkbox.SetChecked(loop)

		eui.GetUi().KeyboardListener().Connect(func(data eui.Event) {
			kd := data.Value.(eui.KeyboardData)
			if kd.IsReleased(ebiten.KeySpace) {
				playMusic()
			}
			if kd.IsReleased(ebiten.KeyEnter) {
				playSfx()
			}
			if kd.IsReleased(ebiten.KeyEscape) {
				eui.GetUi().Pop()
			}
		})
		eui.GetUi().TickListener().Connect(func(data eui.Event) {
			if !playerLoop.IsPlaying() && loop && play {
				playerLoop.Rewind()
				playerLoop.Play()
			}
			if !playerLoop.IsPlaying() && !loop && play {
				playerLoop.Rewind()
				playMusic()
			}
			if totalDuration > 0 && !progressBar.State().IsFocused() {
				progressBar.SetValue(float64(playerLoop.Position()) / float64(totalDuration))
			}
			lblTime.SetText(fmt.Sprintf("%v / %vs",
				eui.FormatSmartDuration(playerLoop.Position(), false),
				eui.FormatSmartDuration(totalDuration, false)))

		})
		s.Add(progressBar)
		s.Add(lblTime)
		s.Add(sliderSeek)
		s.Add(sliderVolume)
		s.Add(lblVolume)

		contBtns := eui.NewContainer(eui.NewHBoxLayout(10)).
			Add(btnPlay).
			Add(btnSfx)
		s.Add(contBtns)
		s.Add(checkbox)
		return s
	}())
}
