// Пример показывает создание сцены с меткой текста и двумя кнопками увеличить или уменьшить число на единицу. Элементы сцены содержаться в двух контейнерах с автоматической разметкой элементов внутри контейнера.

package main

import (
	"strconv"
	"time"

	"github.com/t0l1k/eui"
)

type Count struct{ *eui.Signal[int] }

func NewCount() *Count { return &Count{Signal: eui.NewSignal(func(a, b int) bool { return a == b })} }
func (c *Count) Inc()  { c.Emit(c.Signal.Value() + 1) }
func (c *Count) Dec() {
	if c.Signal.Value() > 0 {
		c.Emit(c.Signal.Value() - 1)
	}
}
func (c *Count) Value() string { return strconv.Itoa(c.Signal.Value()) }

type SceneCounter struct{ *eui.Scene }

func NewSceneCounter() *SceneCounter {
	sc := &SceneCounter{Scene: eui.NewScene(eui.NewVBoxLayout(1))}              // Контейнер сцены по вертикали
	lay2 := eui.NewContainer(eui.NewHBoxLayout(1))                              // Контейнер кнопок по горизонтали
	lblCount := eui.NewText("")                                                 // Текстовая метка
	count := NewCount()                                                         // Сигнал подписчикам передать автоматически оповещение при изменении переменной
	count.ConnectAndFire(func(data int) { lblCount.SetText(count.Value()) }, 0) // Подписка на уведомления от этой переменной
	sc.Add(lblCount)                                                            // Добавить в контейнер метку
	btnInc := eui.NewButton("+", func(b *eui.Button) { count.Inc() })           // Кнопка увеличить на единицу и передать подписчикам об этом
	btnDec := eui.NewButton("-", func(b *eui.Button) { count.Dec() })           // Кнопка уменьшить на единицу и передать подписчикам об этом
	lay2.Add(btnInc)                                                            // Добавить в контейнер кнопку увеличить
	lay2.Add(btnDec)                                                            // Добавить в контейнер кнопку уменьшить
	sc.Add(lay2)                                                                // Добавить в контейнер сцены контейнер кнопок
	eui.NewSnackBar("Test Counter!!!").Show(3 * time.Second)
	return sc
}

func main() {
	eui.Init(eui.GetUi().SetTitle("Counter").SetSize(320, 200))
	eui.Run(NewSceneCounter())
	eui.Quit(func() {})
}
