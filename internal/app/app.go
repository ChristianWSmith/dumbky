package app

import (
	"dumbky/internal/ui"

	"fyne.io/fyne/v2/app"
)

func Run() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(ui.Hello())
	w.ShowAndRun()
}
