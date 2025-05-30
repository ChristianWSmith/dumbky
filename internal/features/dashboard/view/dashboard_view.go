package view

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type viewImpl struct {
	ui fyne.CanvasObject
}

type Bindables struct {
}

type View interface {
	CanvasObject() fyne.CanvasObject
}

func NewView(dashboardSidebarUI, workspaceUI fyne.CanvasObject) View {

	split := container.NewHSplit(dashboardSidebarUI, workspaceUI)
	split.SetOffset(constants.UI_DASHBOARD_SIDEBAR_OFFSET)

	ui := container.NewBorder(nil, nil, nil, nil, split)
	return &viewImpl{
		ui: ui,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}
