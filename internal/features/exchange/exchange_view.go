package exchange

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui *fyne.Container
}

func newView(headerUI, requestUI, responseUI *fyne.Container) *view {
	requestResponseView := container.NewHSplit(requestUI, responseUI)
	ui := container.NewBorder(headerUI, nil, nil, nil, requestResponseView)
	return &view{
		ui: ui,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}
