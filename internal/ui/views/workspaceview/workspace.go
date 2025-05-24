package workspaceview

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/log"
	"dumbky/internal/ui/views/exchangeview"
	"dumbky/internal/ui/views/workspaceheaderview"
	"encoding/json"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

type WorkspaceView struct {
	UI              *fyne.Container
	WorkspaceHeader workspaceheaderview.WorkspaceHeaderView
	exchangeTabs    *container.DocTabs
	tabMap          map[*container.TabItem]WorkspaceTab
}

type Document struct {
	CollectionName string                     `json:"collection_name"`
	Title          string                     `json:"title"`
	ExchangeState  exchangeview.ExchangeState `json:"exchange"`
}

type WorkspaceTab struct {
	CollectionName string
	Title          string
	ExchangeView   exchangeview.ExchangeView
}

func DocumentToRequest(document Document) (db.Request, error) {
	if document.Title == "" {
		document.Title = constants.UI_PLACEHOLDER_UNTITLED
	}

	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Error(err)
		return db.Request{}, err
	}

	jsonString := string(jsonData)

	return db.Request{
		CollectionName: document.CollectionName,
		Name:           document.Title,
		Payload:        jsonString,
	}, nil
}

func RequestToDocument(request db.Request) (Document, error) {
	log.Info(fmt.Sprintf("%v", request))
	document := Document{}
	unmarshalErr := json.Unmarshal([]byte(request.Payload), &document)
	if unmarshalErr != nil {
		log.Error(unmarshalErr)
		return Document{}, unmarshalErr
	}

	return document, nil
}

func formatTabText(collectionName, title string) string {
	return fmt.Sprintf("%s / %s", collectionName, title)
}

func (wv WorkspaceView) OpenTab(document Document) {
	for fyneTab, workspaceTab := range wv.tabMap {
		if workspaceTab.Title == document.Title && workspaceTab.CollectionName == document.CollectionName {
			log.Debug("tab already exists, selecting")
			wv.exchangeTabs.Select(fyneTab)
			return
		}
	}

	exchangeView := exchangeview.ComposeExchangeView()
	err := exchangeView.LoadState(document.ExchangeState)
	if err != nil {
		log.Error(err)
		return
	}
	exchangeViewTab := container.NewTabItem(formatTabText(document.CollectionName, document.Title), exchangeView.UI)
	wv.tabMap[exchangeViewTab] = WorkspaceTab{
		CollectionName: document.CollectionName,
		Title:          document.Title,
		ExchangeView:   exchangeView,
	}
	wv.exchangeTabs.Append(exchangeViewTab)
	wv.exchangeTabs.Select(exchangeViewTab)
}

func (wv WorkspaceView) SaveTab(callback func()) error {
	workspaceTab, ok := wv.tabMap[wv.exchangeTabs.Selected()]
	if !ok {
		err := errors.New("failed to locate selected tab")
		log.Error(err)
		return err
	}
	exchangeState, exchangeStateErr := workspaceTab.ExchangeView.ToState()
	if exchangeStateErr != nil {
		log.Error(exchangeStateErr)
		return exchangeStateErr
	}

	collectionName := workspaceTab.CollectionName
	title := workspaceTab.Title

	document := Document{CollectionName: collectionName, Title: title, ExchangeState: exchangeState}
	request, err := DocumentToRequest(document)
	if err != nil {
		log.Error(err)
		return err
	}
	saveRequestErr := db.SaveRequest(request)
	if saveRequestErr != nil {
		log.Error(saveRequestErr)
		return err
	}
	callback()
	return nil
}

func (wv WorkspaceView) LoadTab(collectionName, title string) error {
	request, err := db.LoadRequest(collectionName, title)
	if err != nil {
		log.Error(err)
		return err
	}
	document, err := RequestToDocument(request)
	if err != nil {
		log.Error(err)
		return err
	}
	fyne.Do(func() {
		wv.OpenTab(document)
	})
	return nil
}

func ComposeWorkspaceView() WorkspaceView {
	workspaceHeader := workspaceheaderview.ComposeWorkspaceHeaderView()

	exchangeTabs := container.NewDocTabs()

	ui := container.NewBorder(workspaceHeader.UI, nil, nil, nil, exchangeTabs)
	wv := WorkspaceView{
		UI:              ui,
		WorkspaceHeader: workspaceHeader,
		exchangeTabs:    exchangeTabs,
		tabMap:          make(map[*container.TabItem]WorkspaceTab),
	}

	exchangeTabs.OnSelected = func(tabItem *container.TabItem) {
		workspaceTab, ok := wv.tabMap[tabItem]
		if !ok {
			log.Error(errors.New("selected tab not in tabMap (OnSelected)"))
			return
		}
		titleErr := workspaceHeader.TitleBinding.Set(workspaceTab.Title)
		if titleErr != nil {
			log.Error(titleErr)
			return
		}
	}

	exchangeTabs.OnClosed = func(tabItem *container.TabItem) {
		delete(wv.tabMap, tabItem)
		if len(exchangeTabs.Items) == 0 {
			wv.OpenTab(Document{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, Title: constants.UI_PLACEHOLDER_UNTITLED})
		}
	}

	wv.OpenTab(Document{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, Title: constants.UI_PLACEHOLDER_UNTITLED})

	workspaceHeader.TitleBinding.AddListener(binding.NewDataListener(func() {
		selectedTab := wv.exchangeTabs.Selected()
		if selectedTab == nil {
			log.Error(errors.New("no selected tab"))
			return
		}
		workspaceTab, ok := wv.tabMap[selectedTab]
		if !ok {
			log.Error(errors.New("selected tab not in tabMap"))
			return
		}
		title, titleErr := workspaceHeader.TitleBinding.Get()
		if titleErr != nil {
			log.Error(titleErr)
			return
		}
		workspaceTab.Title = title
		wv.tabMap[selectedTab] = workspaceTab
		if workspaceTab.Title == "" {
			exchangeTabs.Selected().Text = formatTabText(workspaceTab.CollectionName, constants.UI_PLACEHOLDER_UNTITLED)
		} else {
			exchangeTabs.Selected().Text = formatTabText(workspaceTab.CollectionName, workspaceTab.Title)
		}
		exchangeTabs.Refresh()

	}))

	return wv
}
