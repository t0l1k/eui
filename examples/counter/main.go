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

func main() {
	eui.Init( // Перед запуском настроить приложение
		eui.GetUi(). // Получить экземляр gui
				SetTitle("Counter"). // Текст окна приложения
				SetSize(320, 200))   // Размер окна приложения <F12> тумблер окно на/из полный экран
	eui.Run(func() *eui.Scene {
		counter := NewCount()                                                             // Сигнал подписчикам передать автоматически оповещение при изменении счетчика
		s := eui.NewScene(eui.NewVBoxLayout(1))                                           // Контейнер сцены по вертикали
		lblCounter := eui.NewText("Count")                                                // Текстовая метка
		s.Add(lblCounter)                                                                 // Добавить в контейнер текстовую метку
		counter.ConnectAndFire(func(data int) { lblCounter.SetText(counter.Value()) }, 0) // Подписка на уведомления от счетчика
		s.Add(eui.NewContainer(                                                           // Добавить в контейнер сцены контейнер кнопок
			eui.NewHBoxLayout(1)).                                          // Контейнер кнопок по горизонтали
			Add(eui.NewButton("+", func(b *eui.Button) { counter.Inc() })). // Добавить в контейнер кнопку увеличить на единицу и передать подписчикам об этом
			Add(eui.NewButton("-", func(b *eui.Button) { counter.Dec() }))) // Добавить в контейнер кнопку уменьшить на единицу и передать подписчикам об этом
		eui.NewSnackBar("Test Counter! Click Escape to quit").Show(3 * time.Second) // Показать сообщение 3 секунды после запуска
		return s
	}())
	eui.Quit(func() {}) // Завершить приложение
}
