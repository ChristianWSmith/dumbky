package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type viewImpl struct {
	ui fyne.CanvasObject
}

type Bindables struct{}
type View interface {
	CanvasObject() fyne.CanvasObject
	GetBindables() Bindables
}

func NewView(dashboardUI fyne.CanvasObject) View {
	ui := container.NewBorder(nil, nil, nil, nil, dashboardUI)

	return &viewImpl{
		ui: ui,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{}
}
