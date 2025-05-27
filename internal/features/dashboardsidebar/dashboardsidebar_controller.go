package dashboardsidebar

import (
	"dumbky/internal/features"
	"dumbky/internal/features/collectionsbrowser"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model                  *model
	view                   *view
	collectionsBrowserCtrl *collectionsbrowser.Controller
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	collectionsBrowserCtrl := collectionsbrowser.NewController()
	model := newModel()
	view := newView(collectionsBrowserCtrl.CanvasObject())
	return &Controller{
		model:                  model,
		view:                   view,
		collectionsBrowserCtrl: collectionsBrowserCtrl,
	}
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) SetSelectedRequestListener(handler func()) {
	c.collectionsBrowserCtrl.SetSelectedRequestListener(handler)
}

func (c *Controller) GetSelectedCollection() string {
	return c.collectionsBrowserCtrl.GetSelectedCollection()
}

func (c *Controller) GetSelectedRequest() string {
	return c.collectionsBrowserCtrl.GetSelectedRequest()
}

func (c *Controller) SetSelectedRequest(requestName string) {
	c.collectionsBrowserCtrl.SetSelectedRequest(requestName)
}

func (c *Controller) LazyRefreshAndShowRequests() {
	c.collectionsBrowserCtrl.LazyRefreshAndShowRequests()
}
