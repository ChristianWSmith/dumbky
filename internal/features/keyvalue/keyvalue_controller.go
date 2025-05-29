package keyvalue

import (
	"dumbky/internal/features"
	"dumbky/internal/state"

	"fyne.io/fyne/v2"
)

type controller struct {
	model *model
	view  *view
}

type KeyValueController interface {
	features.Controller
	Get() (key, value string)
	Validate() error
	IsEnabled() bool
	SetDestroyHandler(handler func())
	LoadState(state state.KeyValueState)
	ToState() state.KeyValueState
}

var _ KeyValueController = (*controller)(nil)

func New(keyValidator, valueValidator func(val string) error) KeyValueController {
	return newController(keyValidator, valueValidator)
}

func newController(keyValidator, valueValidator func(val string) error) *controller {
	model := newModel()
	view := newView()

	c := &controller{
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

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) ToState() state.KeyValueState {
	return c.model.toState()
}

func (c *controller) LoadState(state state.KeyValueState) {
	c.model.loadState(state)
}

func (c *controller) SetDestroyHandler(handler func()) {
	c.view.setDestroyHandler(handler)
}

func (c *controller) IsEnabled() bool {
	return c.model.isEnabled()
}

func (c *controller) Validate() error {
	err := c.view.validateKey()
	if err != nil {
		return err
	}
	return c.view.validateValue()
}

func (c *controller) Get() (key, value string) {
	return c.model.getKey(), c.model.getValue()
}
