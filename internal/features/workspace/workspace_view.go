package workspace

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type view struct {
	ui           *fyne.Container
	exchangeTabs *container.DocTabs
	tabMap       map[*container.TabItem]string
}

func newView(workspaceHeaderUI *fyne.Container) *view {

	exchangeTabs := container.NewDocTabs()

	ui := container.NewBorder(workspaceHeaderUI, nil, nil, nil, exchangeTabs)
	return &view{
		ui:           ui,
		exchangeTabs: exchangeTabs,
		tabMap:       make(map[*container.TabItem]string),
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}
