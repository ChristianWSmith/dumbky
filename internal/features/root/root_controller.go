package root

import (
	"dumbky/internal/features"
	"dumbky/internal/features/dashboard"

	"fyne.io/fyne/v2"
)

type controller struct {
	model               *model
	view                *view
	dashboardController dashboard.DashboardController
}

type RootController interface {
	features.Controller
}

var _ RootController = (*controller)(nil)

func New() RootController {
	dashboardController := dashboard.New()
	return newController(dashboardController)
}

func newController(dashboardController dashboard.DashboardController) *controller {
	model := newModel()
	view := newView(dashboardController.CanvasObject())

	return &controller{
		model:               model,
		view:                view,
		dashboardController: dashboardController,
	}
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}
