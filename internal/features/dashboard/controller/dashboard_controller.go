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
	"fmt"

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

	events.GetBus[events.RequestSelected]().Subscribe(func(e events.RequestSelected) {
		fmt.Println("Collection selected:", e.RequestName)
		collectionName := c.collectionsBrowserCtrl.GetSelectedCollection()
		if collectionName == "" || e.RequestName == "" {
			return
		}
		c.workspaceCtrl.LoadTab(collectionName, e.RequestName)
	})

	c.workspaceCtrl.SetAddHandler(func() {
		collectionName := c.collectionsBrowserCtrl.GetSelectedCollection()

		if collectionName == "" {
			collectionName = constants.DB_DEFAULT_COLLECTION_NAME
		}
		c.workspaceCtrl.OpenTab(state.DocumentState{
			CollectionName: collectionName,
			RequestName:    utils.SillyName()})
	})

	c.workspaceCtrl.SetSaveHandler(func() {
		go c.workspaceCtrl.SaveTab(func() {
			fyne.Do(func() {
				// TODO: the idea here is that the user might have a request open from
				// collection A before deleting collection A.  if they're on the
				// collections view in the browser at that time, they should see
				// the updated collections list.  this is an edge case and maybe
				// shouldn't even be addressed.
				c.collectionsBrowserCtrl.LazyRefreshAndShowRequests()
			})
		})
	})

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}
