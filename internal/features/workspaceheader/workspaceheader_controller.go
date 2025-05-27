package workspaceheader

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

	view.requestNameEntry.Bind(model.requestNameBinding)

	return &Controller{
		model: model,
		view:  view,
	}
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) GetRequestName() string {
	return c.model.getRequestName()
}
func (c *Controller) SetRequestName(requestName string) {
	c.model.setRequestName(requestName)
}

func (c *Controller) SetRequestNameListener(handler func()) {
	c.model.setRequestNameListener(handler)
}

func (c *Controller) SetAddHandler(handler func()) {
	c.view.setAddHandler(handler)
}

func (c *Controller) SetSaveHandler(handler func()) {
	c.view.setSaveHandler(handler)
}
