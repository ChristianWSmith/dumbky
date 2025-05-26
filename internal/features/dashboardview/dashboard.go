package dashboardview

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/dashboardsidebarview"
	"dumbky/internal/features/workspace"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

type DashboardView struct {
	UI *fyne.Container
}

func ComposeDashboardView() DashboardView {
	dashboardSidebarView := dashboardsidebarview.ComposeDashboardSidebarView()
	workspaceCtrl := workspace.NewController()

	dashboardSidebarView.CollectionsBrowserView.SelectedRequestBinding.AddListener(binding.NewDataListener(func() {
		collectionName, err := dashboardSidebarView.CollectionsBrowserView.SelectedCollectionBinding.Get()
		if err != nil {
			log.Error(err)
			return
		}
		requestName, err := dashboardSidebarView.CollectionsBrowserView.SelectedRequestBinding.Get()
		if err != nil {
			log.Error(err)
			return
		}
		err = dashboardSidebarView.CollectionsBrowserView.SelectedRequestBinding.Set("")
		if err != nil {
			log.Error(err)
		}
		if collectionName == "" || requestName == "" {
			return
		}
		workspaceCtrl.LoadTab(collectionName, requestName)
	}))

	workspaceCtrl.WorkspaceHeader.SetAddHandler(func() {
		collectionName, err := dashboardSidebarView.CollectionsBrowserView.SelectedCollectionBinding.Get()
		if err != nil {
			log.Error(err)
			return
		}
		if collectionName == "" {
			collectionName = constants.DB_DEFAULT_COLLECTION_NAME
		}
		workspaceCtrl.OpenTab(workspace.Document{
			CollectionName: collectionName,
			Title:          constants.UI_PLACEHOLDER_UNTITLED})
	})

	workspaceCtrl.WorkspaceHeader.SetSaveHandler(func() {
		go workspaceCtrl.SaveTab(func() {
			fyne.Do(func() {
				err := dashboardSidebarView.CollectionsBrowserView.RefreshRequests()
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
