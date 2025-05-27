package exchange

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui fyne.CanvasObject
}

func newView(headerUI, requestUI, responseUI fyne.CanvasObject) *view {
	requestResponseView := container.NewHSplit(requestUI, responseUI)
	ui := container.NewBorder(headerUI, nil, nil, nil, requestResponseView)
	return &view{
		ui: ui,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}
