package dashboard

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/features/dashboardsidebar"
	"dumbky/internal/features/workspace"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model                *model
	view                 *view
	dashboardSidebarCtrl *dashboardsidebar.Controller
	workspaceCtrl        *workspace.Controller
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {

	dashboardSidebarCtrl := dashboardsidebar.NewController()
	workspaceCtrl := workspace.NewController()
	model := newModel()
	view := newView(dashboardSidebarCtrl.GetUI(), workspaceCtrl.GetUI())

	c := &Controller{
		model:                model,
		view:                 view,
		dashboardSidebarCtrl: dashboardSidebarCtrl,
		workspaceCtrl:        workspaceCtrl,
	}

	c.dashboardSidebarCtrl.SetSelectedRequestListener(func() {
		collectionName := c.dashboardSidebarCtrl.GetSelectedCollection()

		requestName := c.dashboardSidebarCtrl.GetSelectedRequest()
		c.dashboardSidebarCtrl.SetSelectedRequest("")
		if collectionName == "" || requestName == "" {
			return
		}
		c.workspaceCtrl.LoadTab(collectionName, requestName)
	})

	c.workspaceCtrl.SetAddHandler(func() {
		collectionName := c.dashboardSidebarCtrl.GetSelectedCollection()

		if collectionName == "" {
			collectionName = constants.DB_DEFAULT_COLLECTION_NAME
		}
		c.workspaceCtrl.OpenTab(workspace.DocumentState{
			CollectionName: collectionName,
			RequestName:    constants.UI_PLACEHOLDER_UNTITLED})
	})

	c.workspaceCtrl.SetSaveHandler(func() {
		go c.workspaceCtrl.SaveTab(func() {
			fyne.Do(func() {
				// TODO: the idea here is that the user might have a request open from
				// collection A before deleting collection A.  if they're on the
				// collections view in the browser at that time, they should see
				// the updated collections list.  this is an edge case and maybe
				// shouldn't even be addressed.
				c.dashboardSidebarCtrl.LazyRefreshAndShowRequests()
			})
		})
	})

	return c
}

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
}
