package keyvalue

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui fyne.CanvasObject

	destroyButton *widget.Button
	enabledCheck  *widget.Check
	keyEntry      *widget.Entry
	valueEntry    *widget.Entry
}

func newView() *view {

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_KEY)
	keyEntry.TextStyle.Monospace = true

	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_VALUE)
	valueEntry.TextStyle.Monospace = true

	enabledCheck := widget.NewCheck(constants.UI_LABEL_KEY_VALUE_ENABLE, nil)

	destroyButton := widget.NewButtonWithIcon("", theme.ContentClearIcon(), nil)

	grid := container.NewGridWithColumns(2, keyEntry, valueEntry)

	ui := container.NewBorder(nil, nil, enabledCheck, destroyButton, grid)

	return &view{
		ui:            ui,
		destroyButton: destroyButton,
		enabledCheck:  enabledCheck,
		keyEntry:      keyEntry,
		valueEntry:    valueEntry,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *view) validateKey() error {
	return v.keyEntry.Validate()
}

func (v *view) validateValue() error {
	return v.valueEntry.Validate()
}

func (v *view) setEnabled(enabled bool) {
	if enabled {
		v.keyEntry.Enable()
		v.valueEntry.Enable()
	} else {
		v.keyEntry.Disable()
		v.valueEntry.Disable()
	}
}

func (v *view) setDestroyHandler(handler func()) {
	v.destroyButton.OnTapped = handler
}

func (v *view) setKeyValidator(validator func(string) error) {
	v.keyEntry.Validator = validator
}

func (v *view) setValueValidator(validator func(string) error) {
	v.keyEntry.Validator = validator
}
