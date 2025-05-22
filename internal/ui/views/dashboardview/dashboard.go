package dashboardview

import (
	"dumbky/internal/constants"
	"dumbky/internal/ui/views/dashboardsidebarview"
	"dumbky/internal/ui/views/workspaceview"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type DashboardView struct {
	UI *fyne.Container
}

func ComposeDashboardView() DashboardView {
	dashboardSidebarView := dashboardsidebarview.ComposeDashboardSidebarView()
	workspaceView := workspaceview.ComposeWorkspaceView()

	split := container.NewHSplit(dashboardSidebarView.UI, workspaceView.UI)
	split.SetOffset(constants.UI_DASHBOARD_SIDEBAR_OFFSET)

	ui := container.NewBorder(nil, nil, nil, nil, split)

	return DashboardView{
		UI: ui,
	}
}
