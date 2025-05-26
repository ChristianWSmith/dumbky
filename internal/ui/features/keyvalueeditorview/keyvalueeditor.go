package keyvalueeditorview

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/ui/features/keyvalue"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type KeyValueEditorView struct {
	UI *fyne.Container

	keyValues      map[*keyvalue.Controller]bool
	keyValueBox    *fyne.Container
	keyValidator   func(val string) error
	valueValidator func(val string) error
}

type KeyValueEditorState struct {
	KeyValueStates []keyvalue.KeyValueState `json:"keyValueStates"`
}

func (kve KeyValueEditorView) ToState() (KeyValueEditorState, error) {
	keyValueStates := []keyvalue.KeyValueState{}
	for keyValue := range kve.keyValues {
		keyValueState := keyValue.ToState()
		keyValueStates = append(keyValueStates, keyValueState)
	}
	return KeyValueEditorState{
		KeyValueStates: keyValueStates,
	}, nil
}

func (kve KeyValueEditorView) clear() {
	kve.keyValues = make(map[*keyvalue.Controller]bool)
	fyne.Do(func() {
		kve.keyValueBox.RemoveAll()
		kve.keyValueBox.Refresh()
	})
}

func (kve KeyValueEditorView) LoadState(keyValueEditorState KeyValueEditorState) error {
	kve.clear()
	for _, keyValueState := range keyValueEditorState.KeyValueStates {
		err := kve.addKeyValue(keyValueState)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (kve KeyValueEditorView) Validate() error {
	for _, kv := range kve.collectEnabled() {
		err := kv.ValidateKey()
		if err != nil {
			log.Warn(err)
			return err
		}
		err = kv.ValidateValue()
		if err != nil {
			log.Warn(err)
			return err
		}
	}
	return nil
}

func (kve KeyValueEditorView) GetMap() (map[string]string, error) {
	out := make(map[string]string)
	for _, kv := range kve.collectEnabled() {
		key := kv.GetKey()
		value := kv.GetValue()
		out[key] = value
	}
	return out, nil
}

func (kve KeyValueEditorView) addKeyValue(keyValueState keyvalue.KeyValueState) error {
	keyValue := keyvalue.NewController(kve.keyValidator, kve.valueValidator)
	keyValue.LoadState(keyValueState)
	keyValue.SetDestroyOnTapped(func() {
		delete(kve.keyValues, keyValue)
		fyne.Do(func() {
			kve.keyValueBox.Remove(keyValue.GetUI())
			kve.keyValueBox.Refresh()
		})
	})
	kve.keyValues[keyValue] = true
	fyne.Do(func() {
		kve.keyValueBox.Add(keyValue.GetUI())
		kve.keyValueBox.Refresh()
	})
	return nil
}

func (kve KeyValueEditorView) collectEnabled() []*keyvalue.Controller {
	out := []*keyvalue.Controller{}
	for kv := range kve.keyValues {
		enabled := kv.IsEnabled()
		if !enabled {
			continue
		}
		out = append(out, kv)
	}
	return out
}

func ComposeKeyValueEditorView(keyValidator, valueValidator func(val string) error) KeyValueEditorView {

	keyValueViews := make(map[*keyvalue.Controller]bool)
	keyValueBox := container.NewVBox()

	addButton := widget.NewButtonWithIcon(constants.UI_LABEL_KEY_VALUE_ADD, nil, nil)
	addButton.Icon = addButton.Theme().Icon(theme.IconNameContentAdd)

	keyValueAddBox := container.NewVBox(keyValueBox, addButton)

	scroll := container.NewVScroll(keyValueAddBox)
	ui := container.NewBorder(nil, nil, nil, nil, scroll)
	kve := KeyValueEditorView{
		ui,
		keyValueViews,
		keyValueBox,
		keyValidator,
		valueValidator,
	}

	addButton.OnTapped = func() {
		go kve.addKeyValue(keyvalue.KeyValueState{Enabled: true, Key: "", Value: ""})
	}

	return kve
}
