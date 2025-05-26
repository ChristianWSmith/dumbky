package keyvalue

import (
	"dumbky/internal/features"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Controller struct {
	model *model
	view  *view
}

var _ features.Controller = (*Controller)(nil)

func NewController(keyValidator, valueValidator func(val string) error) *Controller {
	model := newModel()
	view := newView()

	view.key.Bind(model.key)
	view.value.Bind(model.value)
	view.enabled.Bind(model.enabled)

	view.setKeyValidator(keyValidator)
	view.setValueValidator(valueValidator)

	model.enabled.AddListener(binding.NewDataListener(func() {
		enabled, _ := model.enabled.Get()
		view.setEnabled(enabled)
	}))

	return &Controller{
		model: model,
		view:  view,
	}
}

func (c *Controller) SetDestroyHandler(handler func()) {
	c.view.setDestroyHandler(handler)
}

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
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
