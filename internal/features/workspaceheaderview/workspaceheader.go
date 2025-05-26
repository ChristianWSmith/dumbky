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
	SaveButton   *widget.Button
}

func ComposeWorkspaceHeaderView() WorkspaceHeaderView {
	titleBind := binding.NewString()
	titleEntry := widget.NewEntry()
	titleEntry.Bind(titleBind)

	addButton := widget.NewButtonWithIcon("", nil, nil)
	addButton.Icon = addButton.Theme().Icon(theme.IconNameContentAdd)

	saveButton := widget.NewButtonWithIcon("", nil, nil)
	saveButton.Icon = addButton.Theme().Icon(theme.IconNameDocumentSave)

	controlsLeft := container.NewHBox(addButton)
	controlsRight := container.NewHBox(saveButton)

	ui := container.NewBorder(nil, nil, controlsLeft, controlsRight, titleEntry)
	return WorkspaceHeaderView{
		UI:           ui,
		TitleBinding: titleBind,
		AddButton:    addButton,
		SaveButton:   saveButton,
	}
}
