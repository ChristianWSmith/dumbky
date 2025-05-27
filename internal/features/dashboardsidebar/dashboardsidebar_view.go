package dashboardsidebar

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui *fyne.Container
}

func newView(collectionsBrowserUI *fyne.Container) *view {
	ui := container.NewBorder(nil, nil, nil, nil, collectionsBrowserUI)
	return &view{
		ui: ui,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}
