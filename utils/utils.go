package utils

func ValueFromPercent(percent, total float64) float64 { return percent * total / 100 }
func PercentOf(value, total float64) float64          { return 100 * value / total }

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
