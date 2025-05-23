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
	UI                  *fyne.Container
	workspaceHeader     workspaceheaderview.WorkspaceHeaderView
	exchangeTabs        *container.DocTabs
	tabsToExchangeViews map[*container.TabItem]exchangeview.ExchangeView
}

type Document struct {
	CollectionName string                     `json:"collection_name"`
	Title          string                     `json:"title"`
	ExchangeState  exchangeview.ExchangeState `json:"exchange"`
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

func saveTab(collectionName, title string, exchangeState exchangeview.ExchangeState) error {
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
	return nil
}

func loadTab(collectionName, title string, workspaceView WorkspaceView) error {
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
		workspaceView.addTab(document)
	})
	return nil
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
		go saveTab("", title, exchangeState) // TODO: pass collection name
	}

	workspaceHeader.LoadButton.OnTapped = func() {
		title, titleErr := wv.workspaceHeader.TitleBinding.Get()
		if titleErr != nil {
			log.Error(titleErr)
			return
		}
		go loadTab("", title, wv) // TODO: remove empty string collection name
	}

	exchangeTabs.OnClosed = func(tabItem *container.TabItem) {
		if len(exchangeTabs.Items) == 0 {
			wv.addTab(Document{Title: constants.UI_PLACEHOLDER_UNTITLED})
		}
	}

	wv.addTab(Document{Title: constants.UI_PLACEHOLDER_UNTITLED})
	return wv
}
