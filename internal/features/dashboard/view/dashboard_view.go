package view

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui     fyne.CanvasObject
	offset float64
}

type Bindables struct {
}

type View interface {
	CanvasObject() fyne.CanvasObject
	GetBindables() Bindables
}

func NewView(dashboardSidebarUI, workspaceUI fyne.CanvasObject) View {

	showHideButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), nil)

	sidebar := container.NewBorder(nil, nil, nil, showHideButton, dashboardSidebarUI)

	split := container.NewHSplit(sidebar, workspaceUI)
	split.SetOffset(constants.UI_DASHBOARD_SIDEBAR_OFFSET)

	ui := container.NewBorder(nil, nil, nil, nil, split)
	v := &viewImpl{
		ui:     ui,
		offset: constants.UI_DASHBOARD_SIDEBAR_OFFSET,
	}

	showHideButton.OnTapped = func() {
		if dashboardSidebarUI.Visible() {
			showHideButton.Icon = theme.NavigateNextIcon()
			v.offset = split.Offset
			dashboardSidebarUI.Hide()
			split.SetOffset(0.0)
			split.Refresh()
		} else {
			showHideButton.Icon = theme.NavigateBackIcon()
			dashboardSidebarUI.Show()
			split.SetOffset(v.offset)
			split.Refresh()
		}
	}

	showHideButton.Tapped(nil)

	return v
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{}
}
