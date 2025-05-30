package controller

import (
	"dumbky/internal/features/keyvalue"
	"dumbky/internal/features/keyvalueeditor/model"
	"dumbky/internal/features/keyvalueeditor/view"
	"dumbky/internal/log"
	"dumbky/internal/state"

	"fyne.io/fyne/v2"
)

type controllerImpl struct {
	model model.Model
	view  view.View

	keyValueCtrls  map[keyvalue.KeyValueController]bool
	keyValidator   func(val string) error
	valueValidator func(val string) error
}

func NewController(keyValidator, valueValidator func(string) error) *controllerImpl {
	model := model.NewModel()
	view := view.NewView()
	keyValueControllers := make(map[keyvalue.KeyValueController]bool)
	c := &controllerImpl{
		model:          model,
		view:           view,
		keyValueCtrls:  keyValueControllers,
		keyValidator:   keyValidator,
		valueValidator: valueValidator,
	}
	view.SetAddHandler(func() {
		c.addKeyValue(state.KeyValueState{Enabled: true, Key: "", Value: ""})
	})
	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) ToState() state.KeyValueEditorState {
	keyValueStates := []state.KeyValueState{}
	for keyValue := range c.keyValueCtrls {
		keyValueState := keyValue.ToState()
		keyValueStates = append(keyValueStates, keyValueState)
	}
	return state.KeyValueEditorState{
		KeyValueStates: keyValueStates,
	}
}

func (c *controllerImpl) LoadState(keyValueEditorState state.KeyValueEditorState) {
	c.clear()
	for _, keyValueState := range keyValueEditorState.KeyValueStates {
		c.addKeyValue(keyValueState)
	}
}

func (c *controllerImpl) SetVisible(visible bool) {
	c.view.SetVisible(visible)
}

func (c *controllerImpl) Validate() error {
	for _, keyValueCtrl := range c.collectEnabled() {
		err := keyValueCtrl.Validate()
		if err != nil {
			log.Warn(err)
			return err
		}
	}
	return nil
}

func (c *controllerImpl) Get() map[string]string {
	out := make(map[string]string)
	for _, keyValueCtrl := range c.collectEnabled() {
		key, value := keyValueCtrl.Get()
		out[key] = value
	}
	return out
}

func (c *controllerImpl) clear() {
	c.keyValueCtrls = make(map[keyvalue.KeyValueController]bool)
	c.view.Clear()
}

func (c *controllerImpl) addKeyValue(keyValueState state.KeyValueState) {
	keyValueCtrl := keyvalue.New(c.keyValidator, c.valueValidator)
	keyValueCtrl.LoadState(keyValueState)
	keyValueCtrl.SetDestroyHandler(func() {
		delete(c.keyValueCtrls, keyValueCtrl)
		c.view.Remove(keyValueCtrl.CanvasObject())
	})
	c.keyValueCtrls[keyValueCtrl] = true
	c.view.Add(keyValueCtrl.CanvasObject())
}

func (c *controllerImpl) collectEnabled() []keyvalue.KeyValueController {
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
