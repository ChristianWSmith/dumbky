package common

import (
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
	keyEntry.SetPlaceHolder("<key>")
	keyEntry.TextStyle.Monospace = true
	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("<value>")
	valueEntry.TextStyle.Monospace = true
	enabledCheck := widget.NewCheck("", nil)
	destroyButton := widget.NewButton("✖", nil)

	keyEntry.Bind(keyBinding)
	valueEntry.Bind(valueBinding)
	enabledCheck.Bind(enabledBinding)

	enabledCheck.SetChecked(true)

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