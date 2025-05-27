package workspace

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/features"
	"dumbky/internal/features/exchange"
	"dumbky/internal/features/workspaceheader"
	"dumbky/internal/log"
	"encoding/json"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Controller struct {
	model               *model
	view                *view
	workspaceHeaderCtrl *workspaceheader.Controller
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	model := newModel()
	workspaceHeader := workspaceheader.NewController()
	view := newView(workspaceHeader.GetUI())

	c := &Controller{
		model:               model,
		view:                view,
		workspaceHeaderCtrl: workspaceHeader,
	}

	c.view.exchangeTabs.OnSelected = func(tabItem *container.TabItem) {
		tabId, _ := c.view.tabMap[tabItem]
		workspaceTab, _ := c.model.tabMap[tabId]
		workspaceHeader.SetTitle(workspaceTab.Title)
	}

	c.view.exchangeTabs.OnClosed = func(tabItem *container.TabItem) {
		tabId, _ := c.view.tabMap[tabItem]
		delete(c.view.tabMap, tabItem)
		delete(c.model.tabMap, tabId)
		if len(c.view.exchangeTabs.Items) == 0 {
			c.OpenTab(Document{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, Title: constants.UI_PLACEHOLDER_UNTITLED})
		}
	}

	c.OpenTab(Document{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, Title: constants.UI_PLACEHOLDER_UNTITLED})

	workspaceHeader.SetTitleListener(func() {
		selectedTab := c.view.exchangeTabs.Selected()
		if selectedTab == nil {
			log.Error(errors.New("no selected tab"))
			return
		}
		tabId, _ := c.view.tabMap[selectedTab]
		workspaceTab, _ := c.model.tabMap[tabId]
		title := workspaceHeader.GetTitle()
		workspaceTab.Title = title
		c.model.tabMap[tabId] = workspaceTab
		if workspaceTab.Title == "" {
			c.view.exchangeTabs.Selected().Text = formatTabText(workspaceTab.CollectionName, constants.UI_PLACEHOLDER_UNTITLED)
		} else {
			c.view.exchangeTabs.Selected().Text = formatTabText(workspaceTab.CollectionName, workspaceTab.Title)
		}
		c.view.exchangeTabs.Refresh()

	})

	return c
}

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
}

func (c *Controller) SetAddHandler(handler func()) {
	c.workspaceHeaderCtrl.SetAddHandler(handler)
}

func (c *Controller) SetSaveHandler(handler func()) {
	c.workspaceHeaderCtrl.SetSaveHandler(handler)
}

func (c *Controller) OpenTab(document Document) {
	for fyneTab, tabId := range c.view.tabMap {
		workspaceTab, _ := c.model.tabMap[tabId]
		if workspaceTab.Title == document.Title && workspaceTab.CollectionName == document.CollectionName {
			c.view.exchangeTabs.Select(fyneTab)
			return
		}
	}

	exchangeCtrl := exchange.NewController()
	exchangeCtrl.LoadState(document.ExchangeState)
	exchangeViewTab := container.NewTabItem(formatTabText(document.CollectionName, document.Title), exchangeCtrl.GetUI())
	tabId, _ := c.view.tabMap[exchangeViewTab]
	c.model.tabMap[tabId] = WorkspaceTab{
		CollectionName: document.CollectionName,
		Title:          document.Title,
		ExchangeView:   exchangeCtrl,
	}
	c.view.exchangeTabs.Append(exchangeViewTab)
	c.view.exchangeTabs.Select(exchangeViewTab)
}

func (c *Controller) SaveTab(callback func()) error {
	tabId, _ := c.view.tabMap[c.view.exchangeTabs.Selected()]
	workspaceTab, _ := c.model.tabMap[tabId]
	exchangeState := workspaceTab.ExchangeView.ToState()

	collectionName := workspaceTab.CollectionName
	title := workspaceTab.Title

	document := Document{CollectionName: collectionName, Title: title, ExchangeState: exchangeState}
	request, err := documentToRequest(document)
	if err != nil {
		log.Error(err)
		return err
	}
	go func() {
		saveRequestErr := db.SaveRequest(request)
		if saveRequestErr != nil {
			log.Error(saveRequestErr)
		}
		callback()
	}()
	return nil
}

func (c *Controller) LoadTab(collectionName, title string) {
	go func() {
		request, err := db.LoadRequest(collectionName, title)
		if err != nil {
			log.Error(err)
		}
		document, err := requestToDocument(request)
		if err != nil {
			log.Error(err)
		}
		fyne.Do(func() {
			c.OpenTab(document)
		})
	}()
}

func documentToRequest(document Document) (db.Request, error) {
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

func requestToDocument(request db.Request) (Document, error) {
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
