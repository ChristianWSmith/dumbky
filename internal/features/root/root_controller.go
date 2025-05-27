package root

import (
	"dumbky/internal/features"
	"dumbky/internal/features/dashboard"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model               *model
	view                *view
	dashboardController *dashboard.Controller
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	dashboardController := dashboard.NewController()
	model := newModel()
	view := newView(dashboardController.CanvasObject())

	return &Controller{
		model:               model,
		view:                view,
		dashboardController: dashboardController,
	}
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}
