package dashboardsidebarview

import (
	"dumbky/internal/features/collectionsbrowser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type DashboardSidebarView struct {
	UI                     *fyne.Container
	CollectionsBrowserCtrl *collectionsbrowser.Controller
}

func ComposeDashboardSidebarView() DashboardSidebarView {
	collectionsbrowserCtrl := collectionsbrowser.NewController()

	ui := container.NewBorder(nil, nil, nil, nil, collectionsbrowserCtrl.GetUI())
	return DashboardSidebarView{
		UI:                     ui,
		CollectionsBrowserCtrl: collectionsbrowserCtrl,
	}
}
