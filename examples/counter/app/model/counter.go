package model

func init() {
	counter = 0
}

var counter int

func Value() int {
	return counter
}

func Inc() {
	counter++
}

func Dec() {
	if counter > 0 {
		counter--
	}
}
