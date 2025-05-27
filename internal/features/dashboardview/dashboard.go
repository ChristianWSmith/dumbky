package dashboardview

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/dashboardsidebar"
	"dumbky/internal/features/workspace"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type DashboardView struct {
	UI *fyne.Container
}

func ComposeDashboardView() DashboardView {
	dashboardSidebarCtrl := dashboardsidebar.NewController()
	workspaceCtrl := workspace.NewController()

	dashboardSidebarCtrl.SetSelectedRequestListener(func() {
		collectionName := dashboardSidebarCtrl.GetSelectedCollection()

		requestName := dashboardSidebarCtrl.GetSelectedRequest()
		dashboardSidebarCtrl.SetSelectedRequest("")
		if collectionName == "" || requestName == "" {
			return
		}
		workspaceCtrl.LoadTab(collectionName, requestName)
	})

	workspaceCtrl.SetAddHandler(func() {
		collectionName := dashboardSidebarCtrl.GetSelectedCollection()

		if collectionName == "" {
			collectionName = constants.DB_DEFAULT_COLLECTION_NAME
		}
		workspaceCtrl.OpenTab(workspace.DocumentState{
			CollectionName: collectionName,
			RequestName:    constants.UI_PLACEHOLDER_UNTITLED})
	})

	workspaceCtrl.SetSaveHandler(func() {
		go workspaceCtrl.SaveTab(func() {
			fyne.Do(func() {
				// TODO: the idea here is that the user might have a request open from
				// collection A before deleting collection A.  if they're on the
				// collections view in the browser at that time, they should see
				// the updated collections list.  this is an edge case and maybe
				// shouldn't even be addressed.
				dashboardSidebarCtrl.LazyRefreshAndShowRequests()
			})
		})
	})

	split := container.NewHSplit(dashboardSidebarCtrl.GetUI(), workspaceCtrl.GetUI())
	split.SetOffset(constants.UI_DASHBOARD_SIDEBAR_OFFSET)

	ui := container.NewBorder(nil, nil, nil, nil, split)

	return DashboardView{
		UI: ui,
	}
}
