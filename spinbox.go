package eui

import (
	"math"
	"strconv"
)

type SpinBox[T any] struct {
	*Container
	values        []T
	index         *Signal[int]
	SelectedValue *Signal[T]
	toStr         func(T) string
	compare       func(a, b T) int // сравнение для поиска ближайшего
	min, max      int              // индексы min/max (для проверки ввода)
}

func NewSpinBox[T any](
	txt string,
	values []T,
	index int,
	toStr func(T) string,
	compare func(a, b T) int,
	equals func(a, b T) bool,
	enableInput bool,
	width int,
) *SpinBox[T] {
	s := &SpinBox[T]{
		Container:     NewContainer(NewLayoutHorizontalPercent([]int{20, 20, 20, 40}, 3)),
		values:        values,
		index:         NewSignal(func(a, b int) bool { return a == b }),
		SelectedValue: NewSignal(equals),
		toStr:         toStr,
		compare:       compare,
		min:           0,
		max:           len(values) - 1,
	}
	lbl := NewLabel(txt)
	btnInc := NewButton(string('\u25B2'), func(b *Button) {
		i := s.index.Value()
		if i < s.max {
			s.index.Emit(i + 1)
		}
	})
	btnDec := NewButton(string('\u25BC'), func(b *Button) {
		i := s.index.Value()
		if i > s.min {
			s.index.Emit(i - 1)
		}
	})

	var tf Drawabler
	if enableInput {
		tf = NewTextInputLine(func(ti *TextInputLine) {
			// Пример для int, для других типов — свой парсер
			var input T
			ok := false
			switch any(input).(type) {
			case int:
				val, err := strconv.Atoi(ti.Text())
				if err == nil {
					input = any(val).(T)
					ok = true
				}
			}
			if ok {
				closestIdx := s.min
				minDiff := math.MaxInt
				for i, v := range s.values {
					diff := int(math.Abs(float64(s.compare(v, input))))
					if diff < minDiff {
						minDiff = diff
						closestIdx = i
					}
				}
				// Ограничение min/max
				if closestIdx < s.min {
					closestIdx = s.min
				}
				if closestIdx > s.max {
					closestIdx = s.max
				}
				s.index.Emit(closestIdx)
			}
		})
		tf.(*TextInputLine).SetMaxLen(width)
	} else {
		tf = NewLabel("--")
	}

	s.index.ConnectAndFire(func(idx int) {
		val := s.values[idx]
		s.SelectedValue.Emit(val)
		switch d := tf.(type) {
		case *Label:
			d.SetText(s.toStr(val))
		case *TextInputLine:
			d.SetText(s.toStr(val))
		}
	}, index)

	s.Add(btnInc)
	s.Add(tf)
	s.Add(btnDec)
	s.Add(lbl)
	return s
}
