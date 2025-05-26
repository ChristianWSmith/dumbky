package workspaceheader

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui         *fyne.Container
	titleEntry *widget.Entry
	addButton  *widget.Button
	saveButton *widget.Button
}

func newView() *view {
	titleEntry := widget.NewEntry()

	addButton := widget.NewButtonWithIcon("", nil, nil)
	addButton.Icon = addButton.Theme().Icon(theme.IconNameContentAdd)

	saveButton := widget.NewButtonWithIcon("", nil, nil)
	saveButton.Icon = addButton.Theme().Icon(theme.IconNameDocumentSave)

	controlsLeft := container.NewHBox(addButton)
	controlsRight := container.NewHBox(saveButton)

	ui := container.NewBorder(nil, nil, controlsLeft, controlsRight, titleEntry)
	return &view{
		ui:         ui,
		titleEntry: titleEntry,
		addButton:  addButton,
		saveButton: saveButton,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}

func (v *view) setAddHandler(handler func()) {
	v.addButton.OnTapped = handler
}

func (v *view) setSaveHandler(handler func()) {
	v.saveButton.OnTapped = handler
}
