package dashboardsidebar

import (
	"dumbky/internal/features"
	"dumbky/internal/features/collectionsbrowser"

	"fyne.io/fyne/v2"
)

type controller struct {
	model                  *model
	view                   *view
	collectionsBrowserCtrl collectionsbrowser.CollectionsBrowserController
}

type DashboardSidebarController interface {
	features.Controller
	GetSelectedCollection() string
	GetSelectedRequest() string
	SetSelectedRequestCallback(callback func())
	LazyRefreshAndShowRequests()
}

var _ DashboardSidebarController = (*controller)(nil)

func New() DashboardSidebarController {
	collectionsBrowserCtrl := collectionsbrowser.New()
	return newController(collectionsBrowserCtrl)
}

func newController(collectionsBrowserCtrl collectionsbrowser.CollectionsBrowserController) *controller {
	model := newModel()
	view := newView(collectionsBrowserCtrl.CanvasObject())
	return &controller{
		model:                  model,
		view:                   view,
		collectionsBrowserCtrl: collectionsBrowserCtrl,
	}
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) GetSelectedCollection() string {
	return c.collectionsBrowserCtrl.GetSelectedCollection()
}

func (c *controller) GetSelectedRequest() string {
	return c.collectionsBrowserCtrl.GetSelectedRequest()
}

func (c *controller) SetSelectedRequestCallback(callback func()) {
	c.collectionsBrowserCtrl.SetSelectedRequestCallback(callback)
}

func (c *controller) LazyRefreshAndShowRequests() {
	c.collectionsBrowserCtrl.LazyRefreshAndShowRequests()
}
