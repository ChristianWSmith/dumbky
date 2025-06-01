package controller

import (
	"dumbky/internal/constants"
	"dumbky/internal/events"
	"dumbky/internal/features/collectionsbrowser"
	"dumbky/internal/features/dashboard/model"
	"dumbky/internal/features/dashboard/view"
	"dumbky/internal/features/workspace"
	"dumbky/internal/state"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2"
)

type controllerImpl struct {
	model                  model.Model
	view                   view.View
	collectionsBrowserCtrl collectionsbrowser.CollectionsBrowserController
	workspaceCtrl          workspace.WorkspaceController
}

func NewController(collectionsBrowserCtrl collectionsbrowser.CollectionsBrowserController, workspaceCtrl workspace.WorkspaceController) *controllerImpl {
	model := model.NewModel()
	view := view.NewView(collectionsBrowserCtrl.CanvasObject(), workspaceCtrl.CanvasObject())

	c := &controllerImpl{
		model:                  model,
		view:                   view,
		collectionsBrowserCtrl: collectionsBrowserCtrl,
		workspaceCtrl:          workspaceCtrl,
	}

	events.Subscribe(func(requestSelectedEvent events.RequestSelected) {
		if !requestSelectedEvent.IsSelected {
			return
		}
		collectionSelectedEvent := events.Current[events.CollectionSelected]()
		if !collectionSelectedEvent.IsSelected {
			return
		}
		c.workspaceCtrl.LoadTab(collectionSelectedEvent.CollectionName, requestSelectedEvent.RequestName)
	})

	c.workspaceCtrl.SetAddHandler(func() {
		collectionSelectedEvent := events.Current[events.CollectionSelected]()
		collectionName := collectionSelectedEvent.CollectionName
		if !collectionSelectedEvent.IsSelected {
			collectionName = constants.DB_DEFAULT_COLLECTION_NAME
		}
		c.workspaceCtrl.OpenTab(state.DocumentState{
			CollectionName: collectionName,
			RequestName:    utils.SillyName()})
	})

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}
