// Пример использования библиотеки
package main

import (
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/counter/app"
)

func main() {
	ui.Init(app.NewGame())
	ui.Run(app.NewSceneCounter())
	ui.Quit()
}
