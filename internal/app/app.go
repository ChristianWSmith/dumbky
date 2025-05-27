package app

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/root"
	"dumbky/internal/global"
	"dumbky/internal/theme"

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
	a := app.NewWithID(constants.APP_ID)
	a.Settings().SetTheme(theme.AppTheme{})
	w := a.NewWindow("Dumbky")
	global.Window = w

	rootCtrl := root.NewController()
	w.SetContent(rootCtrl.CanvasObject())

	width, height := getWindowSize()
	w.Resize(fyne.NewSize(width, height))

	// TODO: remove?
	// defer glfw.Terminate()
	// w.Canvas().Refresh(rootCtrl.UI)

	w.ShowAndRun()
}
