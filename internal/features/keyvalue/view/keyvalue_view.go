package view

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui fyne.CanvasObject

	destroyButton *widget.Button
	enabledCheck  *widget.Check
	keyEntry      *widget.Entry
	valueEntry    *widget.Entry
}

type Bindables struct {
	Enabled    features.BoolBindable
	Key, Value features.StringBindable
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
	ValidateKey() error
	ValidateValue() error
	SetEnabled(enabled bool)
	SetDestroyHandler(handler func())
	SetKeyValidator(validator func(string) error)
	SetValueValidator(validator func(string) error)
}

func NewView() View {

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

	return &viewImpl{
		ui:            ui,
		destroyButton: destroyButton,
		enabledCheck:  enabledCheck,
		keyEntry:      keyEntry,
		valueEntry:    valueEntry,
	}
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{
		Enabled: v.enabledCheck,
		Key:     v.keyEntry,
		Value:   v.valueEntry,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) ValidateKey() error {
	return v.keyEntry.Validate()
}

func (v *viewImpl) ValidateValue() error {
	return v.valueEntry.Validate()
}

func (v *viewImpl) SetEnabled(enabled bool) {
	if enabled {
		v.keyEntry.Enable()
		v.valueEntry.Enable()
	} else {
		v.keyEntry.Disable()
		v.valueEntry.Disable()
	}
}

func (v *viewImpl) SetDestroyHandler(handler func()) {
	v.destroyButton.OnTapped = handler
}

func (v *viewImpl) SetKeyValidator(validator func(string) error) {
	v.keyEntry.Validator = validator
}

func (v *viewImpl) SetValueValidator(validator func(string) error) {
	v.keyEntry.Validator = validator
}
