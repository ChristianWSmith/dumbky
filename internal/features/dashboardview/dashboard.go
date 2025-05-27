package dashboardview

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/dashboardsidebarview"
	"dumbky/internal/features/workspace"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type DashboardView struct {
	UI *fyne.Container
}

func ComposeDashboardView() DashboardView {
	dashboardSidebarView := dashboardsidebarview.ComposeDashboardSidebarView()
	workspaceCtrl := workspace.NewController()

	dashboardSidebarView.CollectionsBrowserCtrl.SetSelectedRequestListener(func() {
		collectionName := dashboardSidebarView.CollectionsBrowserCtrl.GetSelectedCollection()

		requestName := dashboardSidebarView.CollectionsBrowserCtrl.GetSelectedRequest()
		dashboardSidebarView.CollectionsBrowserCtrl.SetSelectedRequest("")
		if collectionName == "" || requestName == "" {
			return
		}
		workspaceCtrl.LoadTab(collectionName, requestName)
	})

	workspaceCtrl.SetAddHandler(func() {
		collectionName := dashboardSidebarView.CollectionsBrowserCtrl.GetSelectedCollection()

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
				err := dashboardSidebarView.CollectionsBrowserCtrl.RefreshRequests()
				if err != nil {
					log.Error(err)
				}
			})
		})
	})

	split := container.NewHSplit(dashboardSidebarView.UI, workspaceCtrl.GetUI())
	split.SetOffset(constants.UI_DASHBOARD_SIDEBAR_OFFSET)

	ui := container.NewBorder(nil, nil, nil, nil, split)

	return DashboardView{
		UI: ui,
	}
}
