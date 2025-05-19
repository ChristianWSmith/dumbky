package app

import (
	"dumbky/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/fyne-io/glfw-js"
)

var Canvas fyne.Canvas

func getWindowSize() (float32, float32) {
	width := float32(800.0)
	height := float32(600.0)

	monitor := glfw.GetPrimaryMonitor()
	if monitor != nil {
		mode := monitor.GetVideoMode()
		if mode != nil {
			width = float32(mode.Width * 2 / 3)
			height = float32(mode.Height * 2 / 3)
		}
	}

	return width, height
}

func Run() {
	a := app.New()
	w := a.NewWindow("Dumbky")
	Canvas = w.Canvas()

	rootView := ui.ComposeRootView()
	w.SetContent(rootView.UI)

	width, height := getWindowSize()
	w.Resize(fyne.NewSize(width, height))

	// TODO: remove?
	// defer glfw.Terminate()
	// w.Canvas().Refresh(rootView.UI)

	w.ShowAndRun()
}
