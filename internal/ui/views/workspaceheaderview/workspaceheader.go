package workspaceheaderview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type WorkspaceHeaderView struct {
	UI           *fyne.Container
	TitleBinding binding.String
	AddButton    *widget.Button
}

func ComposeWorkspaceHeaderView() WorkspaceHeaderView {
	titleBind := binding.NewString()
	titleEntry := widget.NewEntry()
	titleEntry.Bind(titleBind)

	addButton := widget.NewButtonWithIcon("", nil, nil)
	addButton.Icon = addButton.Theme().Icon(theme.IconNameContentAdd)

	controls := container.NewHBox(addButton)

	ui := container.NewBorder(nil, nil, nil, controls, titleEntry)
	return WorkspaceHeaderView{
		UI:           ui,
		TitleBinding: titleBind,
		AddButton:    addButton,
	}
}
