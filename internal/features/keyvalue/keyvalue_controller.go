package keyvalue

import (
	"dumbky/internal/features"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model *model
	view  *view
}

var _ features.Controller = (*Controller)(nil)

func NewController(keyValidator, valueValidator func(val string) error) *Controller {
	model := newModel()
	view := newView()

	c := &Controller{
		model: model,
		view:  view,
	}

	c.view.keyEntry.Bind(model.key)
	c.view.valueEntry.Bind(model.value)
	c.view.enabledCheck.Bind(model.enabled)

	c.view.setKeyValidator(keyValidator)
	c.view.setValueValidator(valueValidator)

	c.model.setEnabledListener(func() {
		c.view.setEnabled(c.IsEnabled())
	})

	return c
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) ToState() KeyValueState {
	return c.model.toState()
}

func (c *Controller) LoadState(state KeyValueState) {
	c.model.loadState(state)
}

func (c *Controller) SetDestroyHandler(handler func()) {
	c.view.setDestroyHandler(handler)
}

func (c *Controller) IsEnabled() bool {
	return c.model.isEnabled()
}

func (c *Controller) Validate() error {
	err := c.view.validateKey()
	if err != nil {
		return err
	}
	return c.view.validateValue()
}

func (c *Controller) Get() (key, value string) {
	return c.model.getKey(), c.model.getValue()
}
