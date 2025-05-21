package workspaceview

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/ui/views/exchangeview"
	"dumbky/internal/ui/views/workspaceheaderview"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	wv.exchangeTabs.Select(exchangeViewTab)
}

func ComposeWorkspaceView() WorkspaceView {
	workspaceHeader := workspaceheaderview.ComposeWorkspaceHeaderView()

	exchangeTabs := container.NewDocTabs()

	workspaceHeader.TitleBinding.AddListener(binding.NewDataListener(func() {
		title, titleErr := workspaceHeader.TitleBinding.Get()
		if titleErr != nil {
			log.Error(titleErr)
			return
		}
		if exchangeTabs.Selected() != nil {
			if title == "" {
				exchangeTabs.Selected().Text = constants.UI_PLACEHOLDER_UNTITLED
			} else {
				exchangeTabs.Selected().Text = title
			}
			exchangeTabs.Refresh()
		}
	}))

	exchangeTabs.OnSelected = func(tabItem *container.TabItem) {
		text := tabItem.Text
		if tabItem.Text == constants.UI_PLACEHOLDER_UNTITLED {
			text = ""
		}
		titleErr := workspaceHeader.TitleBinding.Set(text)
		if titleErr != nil {
			log.Error(titleErr)
			return
		}
	}

	ui := container.NewBorder(workspaceHeader.UI, nil, nil, nil, exchangeTabs)
	wv := WorkspaceView{
		UI:              ui,
		workspaceHeader: workspaceHeader,
		exchangeTabs:    exchangeTabs,
	}
	workspaceHeader.AddButton.OnTapped = func() {
		wv.addTab(constants.UI_PLACEHOLDER_UNTITLED, exchangeview.ExchangeState{})
	}
	wv.addTab(constants.UI_PLACEHOLDER_UNTITLED, exchangeview.ExchangeState{})
	return wv
}
