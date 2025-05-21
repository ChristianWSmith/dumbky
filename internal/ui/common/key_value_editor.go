package common

import (
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type KeyValueEditorView struct {
	UI *container.Scroll

	keyValueViews  map[KeyValueView]bool
	keyValueBox    *fyne.Container
	keyValidator   func(val string) error
	valueValidator func(val string) error
}

type KeyValueEditorState struct {
	KeyValueStates []KeyValueState `json:"keyValueStates"`
}

func (kve KeyValueEditorView) ToState() (KeyValueEditorState, error) {
	keyValueStates := []KeyValueState{}
	for keyValue := range kve.keyValueViews {
		keyValueState, err := keyValue.ToState()
		if err != nil {
			log.Error("", err.Error())
			return KeyValueEditorState{}, err
		}
		keyValueStates = append(keyValueStates, keyValueState)
	}
	return KeyValueEditorState{
		KeyValueStates: keyValueStates,
	}, nil
}

func (kve KeyValueEditorView) clear() {
	kve.keyValueViews = make(map[KeyValueView]bool)
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
			log.Error("", err.Error())
			return err
		}
	}
	return nil
}

func (kve KeyValueEditorView) Validate() error {
	for _, kv := range kve.collectEnabled() {
		err := kv.ValidateKey()
		if err != nil {
			log.Warn("failed to validate key", err.Error())
			return err
		}
		err = kv.ValidateValue()
		if err != nil {
			log.Warn("failed to validate value", err.Error())
			return err
		}
	}
	return nil
}

func (kve KeyValueEditorView) GetMap() (map[string]string, error) {
	out := make(map[string]string)
	for _, kv := range kve.collectEnabled() {
		key, keyErr := kv.KeyBinding.Get()
		if keyErr != nil {
			log.Error("Failed to get KeyValue Key in GetMap", keyErr.Error())
			return out, keyErr
		}
		value, valueErr := kv.ValueBinding.Get()
		if valueErr != nil {
			log.Error("Failed to get KeyValue Value in GetMap", valueErr.Error())
			return out, valueErr
		}
		out[key] = value
	}
	return out, nil
}

func (kve KeyValueEditorView) addKeyValue(keyValueState KeyValueState) error {
	keyValueView := ComposeKeyValueView(kve.keyValidator, kve.valueValidator)
	err := keyValueView.LoadState(keyValueState)
	if err != nil {
		log.Error("", err.Error())
		return err
	}
	keyValueView.DestroyButton.OnTapped = func() {
		delete(kve.keyValueViews, keyValueView)
		fyne.Do(func() {
			kve.keyValueBox.Remove(keyValueView.UI)
			kve.keyValueBox.Refresh()
		})
	}
	kve.keyValueViews[keyValueView] = true
	fyne.Do(func() {
		kve.keyValueBox.Add(keyValueView.UI)
		kve.keyValueBox.Refresh()
	})
	return nil
}

func (kve KeyValueEditorView) collectEnabled() []KeyValueView {
	out := []KeyValueView{}
	for kv := range kve.keyValueViews {
		enabled, enabledErr := kv.EnabledBinding.Get()
		if enabledErr != nil {
			log.Error("Failed to get EnabledBinding in GetMap", enabledErr.Error())
			continue
		}
		if !enabled {
			log.Debug("Skipping disabled KeyValue")
			continue
		}
		out = append(out, kv)
	}
	return out
}

func ComposeKeyValueEditorView(keyValidator, valueValidator func(val string) error) KeyValueEditorView {

	keyValueViews := make(map[KeyValueView]bool)
	keyValueBox := container.NewVBox()

	addButton := widget.NewButtonWithIcon("", nil, nil)
	addButton.Icon = addButton.Theme().Icon(theme.IconNameContentAdd)

	keyValueAddBox := container.NewVBox(keyValueBox, addButton)

	ui := container.NewVScroll(keyValueAddBox)
	kve := KeyValueEditorView{
		ui,
		keyValueViews,
		keyValueBox,
		keyValidator,
		valueValidator,
	}

	addButton.OnTapped = func() {
		go kve.addKeyValue(KeyValueState{Enabled: true, Key: "", Value: ""})
	}

	return kve
}
