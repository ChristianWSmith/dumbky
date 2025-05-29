package workspaceheader

import (
	"dumbky/internal/features"

	"fyne.io/fyne/v2"
)

type controller struct {
	model *model
	view  *view
}

type WorkspaceHeaderController interface {
	features.Controller
	GetRequestName() string
	SetRequestName(requestName string)
	SetRequestNameListener(handler func())
	SetAddHandler(handler func())
	SetSaveHandler(handler func())
}

var _ WorkspaceHeaderController = (*controller)(nil)

func New() WorkspaceHeaderController {
	return newController()
}

func newController() *controller {
	model := newModel()
	view := newView()

	view.requestNameEntry.Bind(model.requestNameBinding)

	return &controller{
		model: model,
		view:  view,
	}
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) GetRequestName() string {
	return c.model.getRequestName()
}
func (c *controller) SetRequestName(requestName string) {
	c.model.setRequestName(requestName)
}

func (c *controller) SetRequestNameListener(handler func()) {
	c.model.setRequestNameListener(handler)
}

func (c *controller) SetAddHandler(handler func()) {
	c.view.setAddHandler(handler)
}

func (c *controller) SetSaveHandler(handler func()) {
	c.view.setSaveHandler(handler)
}

// TODO: validate request name?
