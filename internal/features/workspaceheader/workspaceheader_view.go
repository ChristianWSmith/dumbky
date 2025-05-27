package workspaceheader

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui               *fyne.Container
	requestNameEntry *widget.Entry
	addButton        *widget.Button
	saveButton       *widget.Button
}

func newView() *view {
	requestNameEntry := widget.NewEntry()

	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)

	saveButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), nil)

	controlsLeft := container.NewHBox(addButton)
	controlsRight := container.NewHBox(saveButton)

	ui := container.NewBorder(nil, nil, controlsLeft, controlsRight, requestNameEntry)
	return &view{
		ui:               ui,
		requestNameEntry: requestNameEntry,
		addButton:        addButton,
		saveButton:       saveButton,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *view) setAddHandler(handler func()) {
	v.addButton.OnTapped = handler
}

func (v *view) setSaveHandler(handler func()) {
	v.saveButton.OnTapped = handler
}
