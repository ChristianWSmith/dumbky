package view

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui fyne.CanvasObject
}

type Bindables struct {
}

type View interface {
	CanvasObject() fyne.CanvasObject
	GetBindables() Bindables
}

func NewView(dashboardSidebarUI, workspaceUI fyne.CanvasObject) View {

	showHideButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), nil)

	mainUI := container.NewBorder(nil, nil, showHideButton, nil, workspaceUI)

	splitUI := container.NewHSplit(dashboardSidebarUI, mainUI)
	splitUI.SetOffset(constants.UI_DASHBOARD_SIDEBAR_OFFSET)
	splitUI.Hide()

	ui := container.NewBorder(nil, nil, nil, nil, mainUI)

	showHideButton.OnTapped = func() {
		if splitUI.Visible() {
			showHideButton.SetIcon(theme.NavigateNextIcon())
			splitUI.Hide()
			ui.RemoveAll()
			ui.Add(mainUI)
			ui.Refresh()
		} else {
			showHideButton.SetIcon(theme.NavigateBackIcon())
			splitUI.Show()
			ui.RemoveAll()
			ui.Add(splitUI)
			ui.Refresh()
		}
	}

	v := &viewImpl{
		ui: ui,
	}

	return v
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{}
}
