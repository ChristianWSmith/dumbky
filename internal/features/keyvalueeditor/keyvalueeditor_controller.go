package keyvalueeditor

import (
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalue"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
)

type controller struct {
	model *model
	view  *view

	keyValueCtrls  map[keyvalue.KeyValueController]bool
	keyValidator   func(val string) error
	valueValidator func(val string) error
}

type KeyValueEditorController interface {
	features.Controller
	ToState() KeyValueEditorState
	LoadState(keyValueEditorState KeyValueEditorState)
	SetVisible(visible bool)
	Validate() error
	Get() map[string]string
}

var _ KeyValueEditorController = (*controller)(nil)

func New(keyValidator, valueValidator func(string) error) KeyValueEditorController {
	return newController(keyValidator, valueValidator)
}

func newController(keyValidator, valueValidator func(string) error) *controller {
	model := newModel()
	view := newView()
	keyValueControllers := make(map[keyvalue.KeyValueController]bool)
	c := &controller{
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

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) ToState() KeyValueEditorState {
	keyValueStates := []keyvalue.KeyValueState{}
	for keyValue := range c.keyValueCtrls {
		keyValueState := keyValue.ToState()
		keyValueStates = append(keyValueStates, keyValueState)
	}
	return KeyValueEditorState{
		KeyValueStates: keyValueStates,
	}
}

func (c *controller) LoadState(keyValueEditorState KeyValueEditorState) {
	c.clear()
	for _, keyValueState := range keyValueEditorState.KeyValueStates {
		c.addKeyValue(keyValueState)
	}
}

func (c *controller) SetVisible(visible bool) {
	if visible {
		c.view.show()
	} else {
		c.view.hide()
	}
}

func (c *controller) Validate() error {
	for _, keyValueCtrl := range c.collectEnabled() {
		err := keyValueCtrl.Validate()
		if err != nil {
			log.Warn(err)
			return err
		}
	}
	return nil
}

func (c *controller) Get() map[string]string {
	out := make(map[string]string)
	for _, keyValueCtrl := range c.collectEnabled() {
		key, value := keyValueCtrl.Get()
		out[key] = value
	}
	return out
}

func (c *controller) clear() {
	c.keyValueCtrls = make(map[keyvalue.KeyValueController]bool)
	c.view.clear()
}

func (c *controller) addKeyValue(keyValueState keyvalue.KeyValueState) {
	keyValueCtrl := keyvalue.New(c.keyValidator, c.valueValidator)
	keyValueCtrl.LoadState(keyValueState)
	keyValueCtrl.SetDestroyHandler(func() {
		delete(c.keyValueCtrls, keyValueCtrl)
		c.view.remove(keyValueCtrl.CanvasObject())
	})
	c.keyValueCtrls[keyValueCtrl] = true
	c.view.add(keyValueCtrl.CanvasObject())
}

func (c *controller) collectEnabled() []keyvalue.KeyValueController {
	out := []keyvalue.KeyValueController{}
	for kv := range c.keyValueCtrls {
		enabled := kv.IsEnabled()
		if !enabled {
			continue
		}
		out = append(out, kv)
	}
	return out
}
