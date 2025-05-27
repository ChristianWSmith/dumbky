package dashboard

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui *fyne.Container
}

func newView(dashboardSidebarUI, workspaceUI *fyne.Container) *view {

	split := container.NewHSplit(dashboardSidebarUI, workspaceUI)
	split.SetOffset(constants.UI_DASHBOARD_SIDEBAR_OFFSET)

	ui := container.NewBorder(nil, nil, nil, nil, split)
	return &view{
		ui: ui,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}
