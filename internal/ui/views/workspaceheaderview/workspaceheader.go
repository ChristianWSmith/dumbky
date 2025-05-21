package workspaceheaderview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WorkspaceHeaderView struct {
	UI *fyne.Container
}

func ComposeWorkspaceHeaderView() WorkspaceHeaderView {
	titleEntry := widget.NewEntry()
	ui := container.NewBorder(nil, nil, nil, nil, titleEntry)
	return WorkspaceHeaderView{
		UI: ui,
	}
}
