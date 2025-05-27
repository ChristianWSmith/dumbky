package root

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui fyne.CanvasObject
}

func newView(dashboardUI fyne.CanvasObject) *view {
	ui := container.NewBorder(nil, nil, nil, nil, dashboardUI)

	return &view{
		ui: ui,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}
