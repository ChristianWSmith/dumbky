package app

import (
	"dumbky/internal/global"
	"dumbky/internal/ui/theme"
	"dumbky/internal/ui/views/rootview"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/fyne-io/glfw-js"
)

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
	a := app.NewWithID("com.christianwsmith.dumbky")
	a.Settings().SetTheme(theme.DumbkyTheme{})
	w := a.NewWindow("Dumbky")
	global.Window = w

	rootView := rootview.ComposeRootView()
	w.SetContent(rootView.UI)

	width, height := getWindowSize()
	w.Resize(fyne.NewSize(width, height))

	// TODO: remove?
	// defer glfw.Terminate()
	// w.Canvas().Refresh(rootView.UI)

	w.ShowAndRun()
}
