package dashboardview

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/ui/views/dashboardsidebarview"
	"dumbky/internal/ui/views/workspaceview"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

type DashboardView struct {
	UI *fyne.Container
}

func ComposeDashboardView() DashboardView {
	dashboardSidebarView := dashboardsidebarview.ComposeDashboardSidebarView()
	workspaceView := workspaceview.ComposeWorkspaceView()

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
		log.Debug("selected request reset")
		if err != nil {
			log.Error(err)
		}
		if collectionName == "" || requestName == "" {
			log.Debug("empty collection or request name for load")
			return
		}
		log.Debug(fmt.Sprintf("loading tab %s %s", collectionName, requestName))
		err = workspaceView.LoadTab(collectionName, requestName)
		if err != nil {
			log.Error(err)
			return
		}
	}))

	workspaceView.WorkspaceHeader.AddButton.OnTapped = func() {
		collectionName, err := dashboardSidebarView.CollectionsBrowserView.SelectedCollectionBinding.Get()
		if err != nil {
			log.Error(err)
			return
		}
		if collectionName == "" {
			collectionName = constants.DB_DEFAULT_COLLECTION_NAME
		}
		workspaceView.OpenTab(workspaceview.Document{
			CollectionName: collectionName,
			Title:          constants.UI_PLACEHOLDER_UNTITLED})
	}

	workspaceView.WorkspaceHeader.SaveButton.OnTapped = func() {
		go workspaceView.SaveTab(func() {
			fyne.Do(func() {
				err := dashboardSidebarView.CollectionsBrowserView.RefreshRequests()
				if err != nil {
					log.Error(err)
				}
			})
		})
	}

	split := container.NewHSplit(dashboardSidebarView.UI, workspaceView.UI)
	split.SetOffset(constants.UI_DASHBOARD_SIDEBAR_OFFSET)

	ui := container.NewBorder(nil, nil, nil, nil, split)

	return DashboardView{
		UI: ui,
	}
}
