package response

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

	c := &Controller{
		model: model,
		view:  view,
	}

	c.view.status.Bind(c.model.status)
	c.view.time.Bind(c.model.time)
	c.view.body.Bind(c.model.body)

	return c
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) SetStatus(s string) {
	c.model.setStatus(s)
}

func (c *Controller) SetTime(s string) {
	c.model.setTime(s)
}

func (c *Controller) SetBody(s string) {
	c.model.setBody(s)
}

func (c *Controller) SetLoading(loading bool) {
	c.view.setLoading(loading)
}

func (c *Controller) SetResponse(status, time, body string) {
	c.model.setResponse(status, time, body)
}
