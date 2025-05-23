package dashboardsidebarview

import (
	"dumbky/internal/ui/views/collectionsbrowserview"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type DashboardSidebarView struct {
	UI                     *fyne.Container
	CollectionsBrowserView collectionsbrowserview.CollectionsBrowserView
}

func ComposeDashboardSidebarView() DashboardSidebarView {
	collectionsBrowserView := collectionsbrowserview.ComposeCollectionsBrowserView()

	ui := container.NewBorder(nil, nil, nil, nil, collectionsBrowserView.UI)
	return DashboardSidebarView{
		UI:                     ui,
		CollectionsBrowserView: collectionsBrowserView,
	}
}
