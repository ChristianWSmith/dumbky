package common

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type KeyValueView struct {
	UI             *fyne.Container
	DestroyButton  *widget.Button
	KeyEntry       *widget.Entry
	ValueEntry     *widget.Entry
	EnabledBinding binding.Bool
	KeyBinding     binding.String
	ValueBinding   binding.String
}

func enabledCheckOnChangedWorker(checked bool, keyEntry, valueEntry *widget.Entry) {
	if checked {
		fyne.Do(func() {
			keyEntry.Enable()
			valueEntry.Enable()
		})
	} else {
		fyne.Do(func() {
			keyEntry.Disable()
			valueEntry.Disable()
		})
	}
}

func ComposeKeyValueView(keyValidator, valueValidator func(val string) error) KeyValueView {
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
	//destroyButton := widget.NewButton(constants.UI_LABEL_REMOVE, nil)

	destroyButton := widget.NewButtonWithIcon("", nil, nil)
	destroyButton.Icon = destroyButton.Theme().Icon(theme.IconNameContentClear)

	keyEntry.Bind(keyBinding)
	valueEntry.Bind(valueBinding)
	enabledCheck.Bind(enabledBinding)

	keyEntry.Validator = keyValidator
	valueEntry.Validator = valueValidator

	enabledCheckOnChangedOld := enabledCheck.OnChanged
	enabledCheck.OnChanged = func(checked bool) {
		enabledCheckOnChangedOld(checked)
		go enabledCheckOnChangedWorker(checked, keyEntry, valueEntry)
	}

	err := enabledBinding.Set(true)
	if err != nil {
		log.Error("Failed to set value for KeyValueView enabledBinding", err.Error())
	}

	keyValue := container.NewGridWithColumns(2, keyEntry, valueEntry)

	ui := container.NewBorder(nil, nil, enabledCheck, destroyButton, keyValue)

	return KeyValueView{
		ui,
		destroyButton,
		keyEntry,
		valueEntry,
		enabledBinding,
		keyBinding,
		valueBinding,
	}
}
