package keyvalueeditor

import (
	"dumbky/internal/features"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model *model
	view  *view
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	model := newModel()
	view := newView()
	return &Controller{
		model: model,
		view:  view,
	}
}

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
}

