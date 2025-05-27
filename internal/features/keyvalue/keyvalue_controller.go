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

func (c *Controller) SetDestroyHandler(handler func()) {
	c.view.setDestroyHandler(handler)
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) ToState() KeyValueState {
	return c.model.toState()
}

func (c *Controller) LoadState(kvs KeyValueState) {
	c.model.loadState(kvs)
}

func (c *Controller) IsEnabled() bool {
	return c.model.isEnabled()
}

func (c *Controller) ValidateKey() error {
	return c.view.validateKey()
}

func (c *Controller) ValidateValue() error {
	return c.view.validateValue()
}

func (c *Controller) GetKey() string {
	return c.model.getKey()
}

func (c *Controller) GetValue() string {
	return c.model.getValue()
}
