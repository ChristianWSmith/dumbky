package controller

import (
	"dumbky/internal/features/template/model"
	"dumbky/internal/features/template/view"

	"fyne.io/fyne/v2"
)

type controllerImpl struct {
	model model.Model
	view  view.View
}

func NewController() *controllerImpl {
	c := &controllerImpl{
		model: model.NewModel(),
		view:  view.NewView(),
	}

	c.bindAll()

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) bindAll() {
	// bindings := c.model.GetBindings()
	// bindables := c.view.GetBindables()
}
