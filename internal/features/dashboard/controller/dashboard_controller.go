package controller

import (
	"dumbky/internal/features/collectionsbrowser"
	"dumbky/internal/features/dashboard/model"
	"dumbky/internal/features/dashboard/view"
	"dumbky/internal/features/workspace"

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

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}
