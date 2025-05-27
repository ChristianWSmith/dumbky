package root

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui *fyne.Container
}

func newView(dashboardUI *fyne.Container) *view {
	ui := container.NewBorder(nil, nil, nil, nil, dashboardUI)

	return &view{
		ui: ui,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}
