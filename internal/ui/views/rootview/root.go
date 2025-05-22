package rootview

import (
	"dumbky/internal/ui/views/workspaceview"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type RootView struct {
	UI *fyne.Container
}

func ComposeRootView() RootView {
	workspaceView := workspaceview.ComposeWorkspaceView()

	ui := container.NewBorder(nil, nil, nil, nil, workspaceView.UI)

	return RootView{
		ui,
	}
}
