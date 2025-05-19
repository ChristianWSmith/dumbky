package common

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type KeyValueEditorView struct {
	UI *container.Scroll
	KeyValues map[KeyValueView]bool
}

func (kve KeyValueEditorView) GetMap(filter func(kv *KeyValueView) bool) map[string]string {
	out := make(map[string]string)
	for kv, _ := range kve.KeyValues {
		/* TODO: handle errors */
		enabled, err := kv.Enabled.Get()
		if err != nil || !enabled {
			continue
		}
		key, keyErr := kv.Key.Get()
		value, valueErr := kv.Value.Get()
		if keyErr == nil && valueErr == nil {
			out[key] = value
		}
	}
	return out
}

func ComposeKeyValueEditorView() KeyValueEditorView {

	keyValueViews := make(map[KeyValueView]bool)
	keyValueBox := container.NewVBox()

	addButton := widget.NewButton("➕", func() {
		keyValueView := ComposeKeyValueView()
		keyValueView.DestroyButton.OnTapped = func () {
			delete(keyValueViews, keyValueView)
			keyValueBox.Remove(keyValueView.UI)
		}
		keyValueViews[keyValueView] = true
		keyValueBox.Add(keyValueView.UI)
		keyValueBox.Refresh()
	})

	addButton.Tapped(nil)

	keyValueAddBox := container.NewVBox(keyValueBox, addButton)

	ui := container.NewVScroll(keyValueAddBox)
	return KeyValueEditorView{
		ui,
		keyValueViews,
	}
}