package rootview

import (
	"dumbky/internal/ui/views/dashboardview"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type RootView struct {
	UI *fyne.Container
}

func ComposeRootView() RootView {
	dashboardView := dashboardview.ComposeDashboardView()

	ui := container.NewBorder(nil, nil, nil, nil, dashboardView.UI)

	return RootView{
		ui,
	}
}
