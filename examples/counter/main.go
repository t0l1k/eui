// Пример показывает создание сцены с меткой текста и двумя кнопками увеличить или уменьшить число на единицу. Элементы сцены содержаться в двух контейнерах с автоматической разметкой элементов внутри контейнера, в этих это по горизонтали или вертикали.

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

type Count struct{ eui.SubjectBase }

func NewCount() *Count {
	c := &Count{}
	c.SetValue(0)
	return c
}

func (c *Count) Inc() {
	c.SetValue(c.Value().(int) + 1)
}

func (c *Count) Dec() {
	if c.Value().(int) > 0 {
		c.SetValue(c.Value().(int) - 1)
	}
}

type SceneCounter struct {
	eui.SceneBase
	lay1, lay2 *eui.BoxLayout
}

func NewSceneCounter() *SceneCounter {
	sc := &SceneCounter{}
	sc.lay1 = eui.NewVLayout() // Контейнер по вертикали
	sc.lay2 = eui.NewHLayout() // Контейнер по горизонтали

	sb := eui.NewSnackBar("Test Counter!!!").Show(3000)
	sc.Add(sb)

	count := NewCount()                     // Подписчикам передать оповещение при изменении переменной
	lblCount := eui.NewText(count.String()) // Текстовая метка
	count.Attach(lblCount)                  // Подписка на уведомления от этой переменной
	count.SetValue(count.Value())
	sc.lay1.Add(lblCount) // Добавить в контейнер метку

	btnInc := eui.NewButton("+", func(b *eui.Button) {
		count.Inc()
	}) // Кнопка увеличить на единицу и передать подписчикам об этом

	btnDec := eui.NewButton("-", func(b *eui.Button) {
		count.Dec()
	}) // Кнопка уменьшить на единицу и передать подписчикам об этом

	sc.lay2.Add(btnInc) // Добавить в контейнер кнопку увеличить
	sc.lay2.Add(btnDec) // Добавить в контейнер кнопку уменьшить
	sc.Resize()         // Метод обновить размеры сцены
	return sc
}

// Наверно многословно обновление и рисование, но пока так
func (s *SceneCounter) Update(dt int) {
	for _, v := range s.lay1.GetContainer() {
		v.Update(dt)
	}
	for _, v := range s.lay2.GetContainer() {
		v.Update(dt)
	}
	s.SceneBase.Update(dt)
}

func (s *SceneCounter) Draw(surface *ebiten.Image) {
	for _, v := range s.lay1.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range s.lay2.GetContainer() {
		v.Draw(surface)
	}
	s.SceneBase.Draw(surface)
}

func (s *SceneCounter) Resize() {
	s.SceneBase.Resize()
	w0, h0 := eui.GetUi().Size() // Получить размеры окна, и для сцены это всё окно в расспоряжении
	s.lay1.Resize([]int{0, 0, w0, h0 / 2})
	s.lay2.Resize([]int{0, h0 / 2, w0, h0 / 2})
}

// Тут можно настроить внешний вид, размеры окна, подготовить данные...
func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Counter")
	k := 1
	w, h := 320*k, 200*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneCounter())
	eui.Quit()
}
