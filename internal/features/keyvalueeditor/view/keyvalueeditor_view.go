package view

import (
	"dumbky/internal/components"
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui                *fyne.Container
	keyValueContainer *fyne.Container
	addButton         *widget.Button
}

type Bindables struct {
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
	SetVisible(visible bool)
	SetAddHandler(handler func())
	Clear()
	Add(ui fyne.CanvasObject)
	Remove(ui fyne.CanvasObject)
}

func NewView() View {
	keyValueBox := container.NewVBox()

	addButton := widget.NewButtonWithIcon(constants.UI_LABEL_KEY_VALUE_ADD, theme.ContentAddIcon(), nil)

	keyValueAddContainer := container.NewVBox(keyValueBox, addButton)

	scroll := components.NewScrollInterceptorWrapper(keyValueAddContainer)
	ui := container.NewBorder(nil, nil, nil, nil, scroll)

	return &viewImpl{
		ui:                ui,
		keyValueContainer: keyValueBox,
		addButton:         addButton,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{}
}

func (v *viewImpl) SetVisible(visible bool) {
	if visible {
		v.ui.Show()
	} else {
		v.ui.Hide()
	}
}

func (v *viewImpl) SetAddHandler(handler func()) {
	v.addButton.OnTapped = handler
}

func (v *viewImpl) Clear() {
	v.keyValueContainer.RemoveAll()
	v.keyValueContainer.Refresh()
}

func (v *viewImpl) Add(ui fyne.CanvasObject) {
	v.keyValueContainer.Add(ui)
	v.keyValueContainer.Refresh()
}

func (v *viewImpl) Remove(ui fyne.CanvasObject) {
	v.keyValueContainer.Remove(ui)
	v.keyValueContainer.Refresh()
}
