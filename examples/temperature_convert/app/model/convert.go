package model

func GetCelsiusFromFahrenheit(f float64) float64 {
	return (5.0 / 9) * (f - 32)
}

func GetFahrenheitFromCelsius(c float64) float64 {
	return c*9/5 + 32
}
