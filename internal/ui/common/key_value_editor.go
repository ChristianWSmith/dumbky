package common

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

func addButtonOnTappedWorker(keyValueViews map[KeyValueView]bool, keyValueBox *fyne.Container) {
	keyValueView := ComposeKeyValueView()
	keyValueView.DestroyButton.OnTapped = func() {
		go destroyButtonOnTappedWorker(keyValueViews, keyValueView, keyValueBox)
	}
	keyValueViews[keyValueView] = true
	fyne.Do(func() {
		keyValueBox.Add(keyValueView.UI)
		keyValueBox.Refresh()
	})
}

func (kve KeyValueEditorView) GetMap(filter func(kv *KeyValueView) bool) map[string]string {
	out := make(map[string]string)
	for kv, _ := range kve.KeyValues {
		enabled, enabledErr := kv.EnabledBinding.Get()
		if enabledErr != nil {
			log.Error("")
			continue
		}
		if !enabled {
			log.Debug("Skipping disabled KeyValue")
			continue
		}
		key, keyErr := kv.KeyBinding.Get()
		if keyErr != nil {
			log.Error("Failed to get KeyValue Key in GetMap")
			continue
		}
		value, valueErr := kv.ValueBinding.Get()
		if valueErr != nil {
			log.Error("Failed to get KeyValue Value in GetMap")
			continue
		}
		out[key] = value
	}
	return out
}

func ComposeKeyValueEditorView() KeyValueEditorView {

	keyValueViews := make(map[KeyValueView]bool)
	keyValueBox := container.NewVBox()

	addButton := widget.NewButton(constants.UI_LABEL_ADD, nil)

	addButton.OnTapped = func() {
		go addButtonOnTappedWorker(keyValueViews, keyValueBox)
	}

	addButton.Tapped(nil)

	keyValueAddBox := container.NewVBox(keyValueBox, addButton)

	ui := container.NewVScroll(keyValueAddBox)
	return KeyValueEditorView{
		ui,
		keyValueViews,
	}
}
