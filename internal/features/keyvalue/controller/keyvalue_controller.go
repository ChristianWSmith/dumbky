package controller

import (
	"dumbky/internal/features/keyvalue/model"
	"dumbky/internal/features/keyvalue/view"
	"dumbky/internal/state"

	"fyne.io/fyne/v2"
)

type controllerImpl struct {
	model model.Model
	view  view.View
}

func NewController(keyValidator, valueValidator func(val string) error) *controllerImpl {
	model := model.NewModel()
	view := view.NewView()

	c := &controllerImpl{
		model: model,
		view:  view,
	}

	c.bindAll()

	c.view.SetKeyValidator(keyValidator)
	c.view.SetValueValidator(valueValidator)

	c.model.SetEnabledListener(func() {
		c.view.SetEnabled(c.IsEnabled())
	})

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) ToState() state.KeyValueState {
	return c.model.ToState()
}

func (c *controllerImpl) LoadState(state state.KeyValueState) {
	c.model.LoadState(state)
}

func (c *controllerImpl) SetDestroyHandler(handler func()) {
	c.view.SetDestroyHandler(handler)
}

func (c *controllerImpl) IsEnabled() bool {
	return c.model.IsEnabled()
}

func (c *controllerImpl) Validate() error {
	err := c.view.ValidateKey()
	if err != nil {
		return err
	}
	return c.view.ValidateValue()
}

func (c *controllerImpl) Get() (key, value string) {
	return c.model.GetKey(), c.model.GetValue()
}

func (c *controllerImpl) bindAll() {
	bindings := c.model.GetBindings()
	bindables := c.view.GetBindables()
	bindables.Enabled.Bind(bindings.Enabled)
	bindables.Key.Bind(bindings.Key)
	bindables.Value.Bind(bindings.Value)
}
