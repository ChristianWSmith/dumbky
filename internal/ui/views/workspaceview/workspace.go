package workspaceview

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/ui/views/exchangeview"
	"dumbky/internal/ui/views/workspaceheaderview"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type WorkspaceView struct {
	UI              *fyne.Container
	workspaceHeader workspaceheaderview.WorkspaceHeaderView
	exchangeTabs    *container.DocTabs
}

func (wv WorkspaceView) addTab(title string, evs exchangeview.ExchangeState) {
	exchangeView := exchangeview.ComposeExchangeView()
	err := exchangeView.LoadState(evs)
	if err != nil {
		log.Error(err)
		return
	}
	exchangeViewTab := container.NewTabItem(title, exchangeView.UI)
	wv.exchangeTabs.Append(exchangeViewTab)
}

func ComposeWorkspaceView() WorkspaceView {
	workspaceHeader := workspaceheaderview.ComposeWorkspaceHeaderView()

	exchangeTabs := container.NewDocTabs()

	ui := container.NewBorder(workspaceHeader.UI, nil, nil, nil, exchangeTabs)
	wv := WorkspaceView{
		UI:              ui,
		workspaceHeader: workspaceHeader,
		exchangeTabs:    exchangeTabs,
	}
	wv.addTab(constants.UI_PLACEHOLDER_UNTITLED, exchangeview.ExchangeState{})
	return wv
}
