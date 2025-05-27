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
	documentSessionMap  map[string]*exchange.Controller
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
		documentSessionMap:  make(map[string]*exchange.Controller),
	}

	c.view.exchangeTabs.OnSelected = func(tabItem *container.TabItem) {
		tabId := c.view.documentViewMap[tabItem]
		documentData := c.model.documentDataMap[tabId]
		workspaceHeader.SetRequestName(documentData.RequestName)
	}

	c.view.exchangeTabs.OnClosed = func(tabItem *container.TabItem) {
		tabId := c.view.documentViewMap[tabItem]
		delete(c.view.documentViewMap, tabItem)
		delete(c.documentSessionMap, tabId)
		delete(c.model.documentDataMap, tabId)
		if len(c.view.exchangeTabs.Items) == 0 {
			c.OpenTab(DocumentState{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, RequestName: constants.UI_PLACEHOLDER_UNTITLED})
		}
	}

	c.OpenTab(DocumentState{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, RequestName: constants.UI_PLACEHOLDER_UNTITLED})

	workspaceHeader.SetRequestNameListener(func() {
		selectedTab := c.view.exchangeTabs.Selected()
		if selectedTab == nil {
			log.Error(errors.New("no selected tab"))
			return
		}
		tabId := c.view.documentViewMap[selectedTab]
		documentData := c.model.documentDataMap[tabId]
		requestName := workspaceHeader.GetRequestName()
		documentData.RequestName = requestName
		c.model.documentDataMap[tabId] = documentData
		if documentData.RequestName == "" {
			c.view.exchangeTabs.Selected().Text = formatTabText(documentData.CollectionName, constants.UI_PLACEHOLDER_UNTITLED)
		} else {
			c.view.exchangeTabs.Selected().Text = formatTabText(documentData.CollectionName, documentData.RequestName)
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

func (c *Controller) OpenTab(document DocumentState) {
	for tabItem, tabId := range c.view.documentViewMap {
		documentData := c.model.documentDataMap[tabId]
		if documentData.RequestName == document.RequestName && documentData.CollectionName == document.CollectionName {
			c.view.exchangeTabs.Select(tabItem)
			return
		}
	}

	exchangeCtrl := exchange.NewController()
	exchangeCtrl.LoadState(document.ExchangeState)
	exchangeViewTab := container.NewTabItem(formatTabText(document.CollectionName, document.RequestName), exchangeCtrl.GetUI())
	tabId := fmt.Sprintf("%s / %s", document.CollectionName, document.RequestName)
	c.view.documentViewMap[exchangeViewTab] = tabId
	c.model.documentDataMap[tabId] = documentData{
		CollectionName: document.CollectionName,
		RequestName:    document.RequestName,
	}
	c.documentSessionMap[tabId] = exchangeCtrl
	c.view.exchangeTabs.Append(exchangeViewTab)
	c.view.exchangeTabs.Select(exchangeViewTab)
}

func (c *Controller) SaveTab(callback func()) error {
	tabId := c.view.documentViewMap[c.view.exchangeTabs.Selected()]
	documentSession := c.documentSessionMap[tabId]
	documentData := c.model.documentDataMap[tabId]
	exchangeState := documentSession.ToState()

	collectionName := documentData.CollectionName
	requestName := documentData.RequestName

	document := DocumentState{CollectionName: collectionName, RequestName: requestName, ExchangeState: exchangeState}
	request, err := documentStateToRequest(document)
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

func (c *Controller) LoadTab(collectionName, requestName string) {
	go func() {
		request, err := db.LoadRequest(collectionName, requestName)
		if err != nil {
			log.Error(err)
		}
		document, err := requestToDocumentState(request)
		if err != nil {
			log.Error(err)
		}
		fyne.Do(func() {
			c.OpenTab(document)
		})
	}()
}

func documentStateToRequest(documentState DocumentState) (db.Request, error) {
	if documentState.RequestName == "" {
		documentState.RequestName = constants.UI_PLACEHOLDER_UNTITLED
	}

	jsonData, err := json.Marshal(documentState)
	if err != nil {
		log.Error(err)
		return db.Request{}, err
	}

	jsonString := string(jsonData)

	return db.Request{
		CollectionName: documentState.CollectionName,
		Name:           documentState.RequestName,
		Payload:        jsonString,
	}, nil
}

func requestToDocumentState(request db.Request) (DocumentState, error) {
	log.Info(fmt.Sprintf("%v", request))
	document := DocumentState{}
	err := json.Unmarshal([]byte(request.Payload), &document)
	if err != nil {
		log.Error(err)
		return DocumentState{}, err
	}

	return document, nil
}

func formatTabText(collectionName, requestName string) string {
	return fmt.Sprintf("%s / %s", collectionName, requestName)
}
