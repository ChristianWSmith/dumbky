package workspaceview

import (
	"dumbky/internal/constants"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/ui/views/exchangeview"
	"dumbky/internal/ui/views/workspaceheaderview"
	"encoding/json"
	"errors"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
)

type WorkspaceView struct {
	UI                  *fyne.Container
	workspaceHeader     workspaceheaderview.WorkspaceHeaderView
	exchangeTabs        *container.DocTabs
	tabsToExchangeViews map[*container.TabItem]exchangeview.ExchangeView
}

type Document struct {
	Title         string                     `json:"title"`
	ExchangeState exchangeview.ExchangeState `json:"exchange"`
}

func (wv WorkspaceView) addTab(document Document) {
	exchangeView := exchangeview.ComposeExchangeView()
	err := exchangeView.LoadState(document.ExchangeState)
	if err != nil {
		log.Error(err)
		return
	}
	exchangeViewTab := container.NewTabItem(document.Title, exchangeView.UI)
	wv.exchangeTabs.Append(exchangeViewTab)
	wv.tabsToExchangeViews[exchangeViewTab] = exchangeView
	wv.exchangeTabs.Select(exchangeViewTab)
}

func (wv WorkspaceView) saveTab(writer fyne.URIWriteCloser, document Document) {
	if writer == nil {
		log.Debug("nil writer")
		return
	}
	if document.Title == "" {
		document.Title = constants.UI_PLACEHOLDER_UNTITLED
	}
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Error(err)
		return
	}

	jsonString := string(jsonData)
	_, writeErr := writer.Write([]byte(jsonString))
	if writeErr != nil {
		log.Error(writeErr)
	}
	writerErr := writer.Close()
	if writerErr != nil {
		log.Error(writerErr)
	}
}

func (wv WorkspaceView) loadTab(reader fyne.URIReadCloser) {
	if reader == nil {
		log.Debug("nil reader")
		return
	}
	jsonData, err := io.ReadAll(reader)
	if err != nil {
		log.Error(err)
		return
	}
	document := Document{}
	unmarshalErr := json.Unmarshal(jsonData, &document)
	if unmarshalErr != nil {
		log.Error(unmarshalErr)
		return
	}
	fyne.Do(func() {
		wv.addTab(document)
	})
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
		UI:                  ui,
		workspaceHeader:     workspaceHeader,
		exchangeTabs:        exchangeTabs,
		tabsToExchangeViews: make(map[*container.TabItem]exchangeview.ExchangeView),
	}

	workspaceHeader.AddButton.OnTapped = func() {
		wv.addTab(Document{Title: constants.UI_PLACEHOLDER_UNTITLED})
	}

	workspaceHeader.SaveButton.OnTapped = func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				log.Error(err)
				return
			}
			title, titleErr := wv.workspaceHeader.TitleBinding.Get()
			if titleErr != nil {
				log.Error(titleErr)
				return
			}
			exchangeView, ok := wv.tabsToExchangeViews[wv.exchangeTabs.Selected()]
			if !ok {
				log.Error(errors.New("failed to locate selected tab"))
				return
			}
			exchangeState, exchangeStateErr := exchangeView.ToState()
			if exchangeStateErr != nil {
				log.Error(exchangeStateErr)
				return
			}
			go wv.saveTab(writer, Document{Title: title, ExchangeState: exchangeState})
		}, global.Window)
	}

	workspaceHeader.LoadButton.OnTapped = func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Error(err)
				return
			}
			go wv.loadTab(reader)
		}, global.Window)
	}

	exchangeTabs.OnClosed = func(tabItem *container.TabItem) {
		if len(exchangeTabs.Items) == 0 {
			wv.addTab(Document{Title: constants.UI_PLACEHOLDER_UNTITLED})
		}
	}

	wv.addTab(Document{Title: constants.UI_PLACEHOLDER_UNTITLED})
	return wv
}
