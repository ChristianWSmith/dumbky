package keyvalue

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui *fyne.Container

	destroy *widget.Button
	enabled *widget.Check
	key     *widget.Entry
	value   *widget.Entry
}

func newView() *view {

	key := widget.NewEntry()
	key.SetPlaceHolder(constants.UI_PLACEHOLDER_KEY)
	key.TextStyle.Monospace = true

	value := widget.NewEntry()
	value.SetPlaceHolder(constants.UI_PLACEHOLDER_VALUE)
	value.TextStyle.Monospace = true

	enabled := widget.NewCheck(constants.UI_LABEL_KEY_VALUE_ENABLE, nil)

	destroy := widget.NewButtonWithIcon("", nil, nil)
	destroy.Icon = destroy.Theme().Icon(theme.IconNameContentClear)

	grid := container.NewGridWithColumns(2, key, value)

	ui := container.NewBorder(nil, nil, enabled, destroy, grid)

	return &view{
		ui:      ui,
		destroy: destroy,
		enabled: enabled,
		key:     key,
		value:   value,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}

func (v *view) validateKey() error {
	return v.key.Validate()
}

func (v *view) validateValue() error {
	return v.value.Validate()
}

func (v *view) setEnabled(enabled bool) {
	if enabled {
		v.key.Enable()
		v.value.Enable()
	} else {
		v.key.Disable()
		v.value.Disable()
	}
}

func (v *view) setDestroyHandler(handler func()) {
	v.destroy.OnTapped = handler
}

func (v *view) setKeyValidator(validator func(string) error) {
	v.key.Validator = validator
}

func (v *view) setValueValidator(validator func(string) error) {
	v.key.Validator = validator
}
