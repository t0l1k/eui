package app

import "strconv"

type Counter int

func NewCounter() Counter {
	return 0
}

func (c Counter) Get() Counter {
	return c
}

func (c *Counter) Inc() {
	*c++
}

func (c *Counter) Dec() {
	if *c > 0 {
		*c--
	}
}

func (c Counter) String() string {
	return strconv.Itoa(int(c))
}
