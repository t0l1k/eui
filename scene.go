package eui

type Scene interface {
	Entered()
	Quit()
	Resize()
	Container
}
