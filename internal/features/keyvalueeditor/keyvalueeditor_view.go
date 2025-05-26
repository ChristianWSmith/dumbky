package keyvalueeditor

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui                *fyne.Container
	keyValueContainer *fyne.Container
	addButton         *widget.Button
}

func newView() *view {
	keyValueBox := container.NewVBox()

	addButton := widget.NewButtonWithIcon(constants.UI_LABEL_KEY_VALUE_ADD, theme.ContentAddIcon(), nil)

	keyValueAddContainer := container.NewVBox(keyValueBox, addButton)

	scroll := container.NewVScroll(keyValueAddContainer)
	ui := container.NewBorder(nil, nil, nil, nil, scroll)

	return &view{
		ui:                ui,
		keyValueContainer: keyValueBox,
		addButton:         addButton,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}

func (v *view) hide() {
	v.ui.Hide()
}

func (v *view) show() {
	v.ui.Show()
}

func (v *view) setAddHandler(handler func()) {
	v.addButton.OnTapped = handler
}

func (v *view) clear() {
	v.keyValueContainer.RemoveAll()
	v.keyValueContainer.Refresh()
}

func (v *view) add(ui *fyne.Container) {
	v.keyValueContainer.Add(ui)
	v.keyValueContainer.Refresh()
}

func (v *view) remove(ui *fyne.Container) {
	v.keyValueContainer.Remove(ui)
	v.keyValueContainer.Refresh()
}
