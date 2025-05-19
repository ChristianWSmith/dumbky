package common

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type KeyValueView struct {
	UI *fyne.Container
	DestroyButton *widget.Button
	Enabled binding.Bool
	Key binding.String
	Value binding.String
}

func ComposeKeyValueView() KeyValueView {
	keyBinding := binding.NewString()
	valueBinding := binding.NewString()
	enabledBinding := binding.NewBool()

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_KEY)
	keyEntry.TextStyle.Monospace = true
	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_VALUE)
	valueEntry.TextStyle.Monospace = true
	enabledCheck := widget.NewCheck("", nil)
	destroyButton := widget.NewButton(constants.UI_LABEL_REMOVE, nil)

	keyEntry.Bind(keyBinding)
	valueEntry.Bind(valueBinding)
	enabledCheck.Bind(enabledBinding)

	enabledCheck.OnChanged = func(checked bool) {
		// TODO: background
		if checked {
			keyEntry.Enable()
			valueEntry.Enable()
		} else {
			keyEntry.Disable()
			valueEntry.Disable()
		}
	}

	err := enabledBinding.Set(true)
	if err != nil { /* TODO: handle error */ }

	keyValue := container.NewGridWithColumns(2, keyEntry, valueEntry)

	ui := container.NewBorder(nil, nil, enabledCheck, destroyButton, keyValue)
	
	return KeyValueView{
		ui,
		destroyButton,
		enabledBinding,
		keyBinding,
		valueBinding,
	}
}