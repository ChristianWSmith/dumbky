package controller

import (
	"dumbky/internal/features/dashboard"
	"dumbky/internal/features/root/model"
	"dumbky/internal/features/root/view"

	"fyne.io/fyne/v2"
)

type controllerImpl struct {
	model               model.Model
	view                view.View
	dashboardController dashboard.DashboardController
}

func NewController(dashboardController dashboard.DashboardController) *controllerImpl {
	model := model.NewModel()
	view := view.NewView(dashboardController.CanvasObject())

	return &controllerImpl{
		model:               model,
		view:                view,
		dashboardController: dashboardController,
	}
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}
