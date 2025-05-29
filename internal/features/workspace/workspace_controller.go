package workspace

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/features"
	"dumbky/internal/features/exchange"
	"dumbky/internal/log"
	"dumbky/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/google/uuid"
)

type controller struct {
	model           *model
	view            *view
	exchangeCtrlMap map[documentId]exchange.ExchangeController
}

type WorkspaceController interface {
	features.Controller
	SetAddHandler(handler func())
	SetSaveHandler(handler func())
	OpenTab(document state.DocumentState)
	SaveTab(callback func()) error
	LoadTab(collectionName, requestName string)
}

var _ WorkspaceController = (*controller)(nil)

type documentId string

func New() WorkspaceController {
	return newController()
}

func newController() *controller {
	model := newModel()
	view := newView()

	view.requestNameEntry.Bind(model.requestNameBinding)

	c := &controller{
		model:           model,
		view:            view,
		exchangeCtrlMap: make(map[documentId]exchange.ExchangeController),
	}

	c.view.setDocumentSelectedHandler(func(tabItem *container.TabItem) {
		tabId := c.view.getDocumentId(tabItem)
		documentData := c.model.getDocumentData(tabId)
		c.model.setRequestName(documentData.requestName)
	})

	c.view.setDocumentClosedHandler(func(tabItem *container.TabItem) {
		id := c.view.getDocumentId(tabItem)
		c.view.destroyDocumentTab(tabItem)
		c.model.destroyDocumentData(id)
		delete(c.exchangeCtrlMap, id)
		if c.model.documentCount() == 0 {
			c.OpenTab(state.DocumentState{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, RequestName: constants.UI_PLACEHOLDER_UNTITLED})
		}
	})

	c.OpenTab(state.DocumentState{CollectionName: constants.DB_DEFAULT_COLLECTION_NAME, RequestName: constants.UI_PLACEHOLDER_UNTITLED})

	c.model.setRequestNameListener(func() {
		id := c.view.getSelectedDocumentId()
		c.model.updateRequestName(id, c.model.getRequestName())
		documentData := c.model.getDocumentData(id)
		if documentData.requestName == "" {
			c.view.setSelectedDocumentText(documentData.collectionName, constants.UI_PLACEHOLDER_UNTITLED)
		} else {
			c.view.setSelectedDocumentText(documentData.collectionName, documentData.requestName)
		}
		c.view.refreshTabs()

	})

	return c
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) OpenTab(document state.DocumentState) {
	if c.view.selectDocumentTabOnCondition(func(id documentId) bool {
		documentData := c.model.getDocumentData(id)
		return documentData.requestName == document.RequestName && documentData.collectionName == document.CollectionName
	}) {
		return
	}

	id := documentId(uuid.New().String())

	exchangeCtrl := exchange.New()
	exchangeCtrl.LoadState(document.ExchangeState)
	c.exchangeCtrlMap[id] = exchangeCtrl

	c.model.setDocumentData(id, documentData{
		collectionName: document.CollectionName,
		requestName:    document.RequestName,
	})

	c.view.addDocumentTab(id, document.CollectionName, document.RequestName, exchangeCtrl.CanvasObject())
}

func (c *controller) SaveTab(callback func()) error {
	id := c.view.getSelectedDocumentId()
	documentData := c.model.getDocumentData(id)
	exchangeState := c.exchangeCtrlMap[id].ToState()

	document := state.DocumentState{
		CollectionName: documentData.collectionName,
		RequestName:    documentData.requestName,
		ExchangeState:  exchangeState}
	request, err := documentStateToRequest(document)
	if err != nil {
		log.Error(err)
		return err
	}
	go func() {
		err := db.SaveRequest(request)
		if err != nil {
			log.Error(err)
		}
		callback()
	}()
	return nil
}

func (c *controller) LoadTab(collectionName, requestName string) {
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

func (c *controller) SetAddHandler(handler func()) {
	c.view.setAddHandler(handler)
}

func (c *controller) SetSaveHandler(handler func()) {
	c.view.setSaveHandler(handler)
}
