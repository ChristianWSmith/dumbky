package response

import (
	"dumbky/internal/ui/features"

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

	view.statusLabel.Bind(model.status)
	view.timeLabel.Bind(model.time)
	view.bodyLabel.Bind(model.body)

	return &Controller{
		model: model,
		view:  view,
	}
}

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
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
