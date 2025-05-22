package dashboardsidebarview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type DashboardSidebarView struct {
	UI *fyne.Container
}

func ComposeDashboardSidebarView() DashboardSidebarView {
	box := container.NewVBox()
	return DashboardSidebarView{
		UI: box,
	}
}
