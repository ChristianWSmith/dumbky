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

	view.titleEntry.Bind(model.titleBinding)

	return &Controller{
		model: model,
		view:  view,
	}
}

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
}

func (c *Controller) GetTitle() string {
	return c.model.getTitle()
}
func (c *Controller) SetTitle(title string) {
	c.model.setTitle(title)
}

func (c *Controller) SetTitleListener(handler func()) {
	c.model.setTitleListener(handler)
}

func (c *Controller) SetAddHandler(handler func()) {
	c.view.setAddHandler(handler)
}

func (c *Controller) SetSaveHandler(handler func()) {
	c.view.setSaveHandler(handler)
}
