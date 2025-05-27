package keyvalueeditor

import (
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalue"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model *model
	view  *view

	keyValueCtrls  map[*keyvalue.Controller]bool
	keyValidator   func(val string) error
	valueValidator func(val string) error
}

var _ features.Controller = (*Controller)(nil)

func NewController(keyValidator, valueValidator func(string) error) *Controller {
	model := newModel()
	view := newView()
	keyValueControllers := make(map[*keyvalue.Controller]bool)
	c := &Controller{
		model:          model,
		view:           view,
		keyValueCtrls:  keyValueControllers,
		keyValidator:   keyValidator,
		valueValidator: valueValidator,
	}
	view.setAddHandler(func() {
		c.addKeyValue(keyvalue.KeyValueState{Enabled: true, Key: "", Value: ""})
	})
	return c
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) Hide() {
	c.view.hide()
}

func (c *Controller) Show() {
	c.view.show()
}

func (c *Controller) ToState() KeyValueEditorState {
	keyValueStates := []keyvalue.KeyValueState{}
	for keyValue := range c.keyValueCtrls {
		keyValueState := keyValue.ToState()
		keyValueStates = append(keyValueStates, keyValueState)
	}
	return KeyValueEditorState{
		KeyValueStates: keyValueStates,
	}
}

func (c *Controller) LoadState(keyValueEditorState KeyValueEditorState) {
	c.clear()
	for _, keyValueState := range keyValueEditorState.KeyValueStates {
		c.addKeyValue(keyValueState)
	}
}

func (c *Controller) Validate() error {
	for _, keyValueCtrl := range c.collectEnabled() {
		err := keyValueCtrl.ValidateKey()
		if err != nil {
			log.Warn(err)
			return err
		}
		err = keyValueCtrl.ValidateValue()
		if err != nil {
			log.Warn(err)
			return err
		}
	}
	return nil
}

func (c *Controller) GetMap() map[string]string {
	out := make(map[string]string)
	for _, keyValueCtrl := range c.collectEnabled() {
		key := keyValueCtrl.GetKey()
		value := keyValueCtrl.GetValue()
		out[key] = value
	}
	return out
}

func (c *Controller) clear() {
	c.keyValueCtrls = make(map[*keyvalue.Controller]bool)
	c.view.clear()
}

func (c *Controller) addKeyValue(keyValueState keyvalue.KeyValueState) {
	keyValueCtrl := keyvalue.NewController(c.keyValidator, c.valueValidator)
	keyValueCtrl.LoadState(keyValueState)
	keyValueCtrl.SetDestroyHandler(func() {
		delete(c.keyValueCtrls, keyValueCtrl)
		c.view.remove(keyValueCtrl.CanvasObject())
	})
	c.keyValueCtrls[keyValueCtrl] = true
	c.view.add(keyValueCtrl.CanvasObject())
}

func (c *Controller) collectEnabled() []*keyvalue.Controller {
	out := []*keyvalue.Controller{}
	for kv := range c.keyValueCtrls {
		enabled := kv.IsEnabled()
		if !enabled {
			continue
		}
		out = append(out, kv)
	}
	return out
}
