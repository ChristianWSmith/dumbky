package common

import (
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type KeyValueEditorView struct {
	UI        *container.Scroll
	KeyValues map[KeyValueView]bool
}

func destroyButtonOnTappedWorker(keyValueViews map[KeyValueView]bool, keyValueView KeyValueView, keyValueBox *fyne.Container) {
	delete(keyValueViews, keyValueView)
	fyne.Do(func() {
		keyValueBox.Remove(keyValueView.UI)
		keyValueBox.Refresh()
	})
}

func addButtonOnTappedWorker(keyValueViews map[KeyValueView]bool, keyValueBox *fyne.Container, keyValidator, valueValidator func(val string) error) {
	keyValueView := ComposeKeyValueView(keyValidator, valueValidator)
	keyValueView.DestroyButton.OnTapped = func() {
		go destroyButtonOnTappedWorker(keyValueViews, keyValueView, keyValueBox)
	}
	keyValueViews[keyValueView] = true
	fyne.Do(func() {
		keyValueBox.Add(keyValueView.UI)
		keyValueBox.Refresh()
	})
}

func (kve KeyValueEditorView) collectEnabled() []KeyValueView {
	out := []KeyValueView{}
	for kv := range kve.KeyValues {
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

func (kve KeyValueEditorView) Validate() error {
	for _, kv := range kve.collectEnabled() {
		err := kv.KeyEntry.Validate()
		if err != nil {
			log.Warn("failed to validate key", err.Error())
			return err
		}
		err = kv.ValueEntry.Validate()
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

func ComposeKeyValueEditorView(keyValidator, valueValidator func(val string) error) KeyValueEditorView {

	keyValueViews := make(map[KeyValueView]bool)
	keyValueBox := container.NewVBox()

	addButton := widget.NewButtonWithIcon("", nil, nil)
	addButton.Icon = addButton.Theme().Icon(theme.IconNameContentAdd)

	addButton.OnTapped = func() {
		go addButtonOnTappedWorker(keyValueViews, keyValueBox, keyValidator, valueValidator)
	}

	keyValueAddBox := container.NewVBox(keyValueBox, addButton)

	ui := container.NewVScroll(keyValueAddBox)
	return KeyValueEditorView{
		ui,
		keyValueViews,
	}
}
