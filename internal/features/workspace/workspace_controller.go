package workspace

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/features"
	"dumbky/internal/features/exchange"
	"dumbky/internal/features/workspaceheader"
	"dumbky/internal/log"
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Controller struct {
	model               *model
	view                *view
	workspaceHeaderCtrl *workspaceheader.Controller
	exchangeCtrlMap     map[documentId]*exchange.Controller
}

var _ features.Controller = (*Controller)(nil)

type documentId string

func NewController() *Controller {
	model := newModel()
	workspaceHeaderCtrl := workspaceheader.NewController()
	view := newView(workspaceHeaderCtrl.GetUI())

	c := &Controller{
		model:               model,
		view:                view,
		workspaceHeaderCtrl: workspaceHeaderCtrl,
		exchangeCtrlMap:     make(map[documentId]*exchange.Controller),
	}

	c.view.setDocumentSelectedHandler(func(tabItem *container.TabItem) {
		tabId := c.view.getDocumentId(tabItem)
		documentData := c.model.getDocumentData(tabId)
		c.workspaceHeaderCtrl.SetRequestName(documentData.RequestName)
	})

	c.view.setDocumentClosedHandler(func(tabItem *container.TabItem) {
		id := c.view.getDocumentId(tabItem)
		c.view.destroyDocumentTab(tabItem)
		c.model.destroyDocumentData(id)
		delete(c.exchangeCtrlMap, id)
		if c.model.documentCount() == 0 {
			c.OpenTab(DocumentState{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, RequestName: constants.UI_PLACEHOLDER_UNTITLED})
		}
	})

	c.OpenTab(DocumentState{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, RequestName: constants.UI_PLACEHOLDER_UNTITLED})

	c.workspaceHeaderCtrl.SetRequestNameListener(func() {
		id := c.view.getSelectedDocumentId()
		c.model.updateRequestName(id, workspaceHeaderCtrl.GetRequestName())
		documentData := c.model.getDocumentData(id)
		if documentData.RequestName == "" {
			c.view.setSelectedDocumentText(formatTabText(documentData.CollectionName, constants.UI_PLACEHOLDER_UNTITLED))
		} else {
			c.view.setSelectedDocumentText(formatTabText(documentData.CollectionName, documentData.RequestName))
		}
		c.view.refreshTabs()

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
		documentData := c.model.getDocumentData(tabId)
		if documentData.RequestName == document.RequestName && documentData.CollectionName == document.CollectionName {
			c.view.exchangeTabs.Select(tabItem)
			return
		}
	}

	exchangeCtrl := exchange.NewController()
	exchangeCtrl.LoadState(document.ExchangeState)
	exchangeViewTab := container.NewTabItem(formatTabText(document.CollectionName, document.RequestName), exchangeCtrl.GetUI())
	id := documentId(fmt.Sprintf("%s / %s", document.CollectionName, document.RequestName)) // TODO: uuid?
	c.view.documentViewMap[exchangeViewTab] = id
	c.model.documentDataMap[id] = documentData{
		CollectionName: document.CollectionName,
		RequestName:    document.RequestName,
	}
	c.exchangeCtrlMap[id] = exchangeCtrl
	c.view.exchangeTabs.Append(exchangeViewTab)
	c.view.exchangeTabs.Select(exchangeViewTab)
}

func (c *Controller) SaveTab(callback func()) error {
	tabId := c.view.documentViewMap[c.view.exchangeTabs.Selected()]
	documentSession := c.exchangeCtrlMap[tabId]
	documentData := c.model.getDocumentData(tabId)
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
