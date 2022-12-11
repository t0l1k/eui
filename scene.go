package eui

type Scene interface {
	Entered()
	Resize()
	Container
}
