package model

import (
	"testing"
)

func TestTemp(t *testing.T) {
	temps := map[float64]float64{
		5:    41,
		10:   50,
		36.6: 97.88,
		50:   122,
		100:  212,
		-40:  -40,
	}
	for c, f := range temps {
		got := GetFahrenheitFromCelsius(c)
		want := f
		if got != want {
			t.Errorf("Error temp convert got :%v C want:%v F", got, want)
		}
	}

	for c, f := range temps {
		got := GetCelsiusFromFahrenheit(f)
		want := c
		if got != want {
			t.Errorf("Error temp convert got:%v F want:%v C", got, want)
		}
	}
}
