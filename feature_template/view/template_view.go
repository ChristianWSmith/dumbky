package view

import (
	"fyne.io/fyne/v2"
)

type viewImpl struct {
	ui fyne.CanvasObject
}

type Bindables struct {
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
}

func NewView() View {
	v := &viewImpl{
		ui: nil,
	}

	return v
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{}
}
