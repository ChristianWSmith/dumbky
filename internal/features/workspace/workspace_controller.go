package workspace

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/features"
	"dumbky/internal/features/exchange"
	"dumbky/internal/features/workspaceheader"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/google/uuid"
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
		c.workspaceHeaderCtrl.SetRequestName(documentData.requestName)
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
		if documentData.requestName == "" {
			c.view.setSelectedDocumentText(documentData.collectionName, constants.UI_PLACEHOLDER_UNTITLED)
		} else {
			c.view.setSelectedDocumentText(documentData.collectionName, documentData.requestName)
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
	if c.view.selectDocumentTabOnCondition(func(id documentId) bool {
		documentData := c.model.getDocumentData(id)
		return documentData.requestName == document.RequestName && documentData.collectionName == document.CollectionName
	}) {
		return
	}

	id := documentId(uuid.New().String())

	exchangeCtrl := exchange.NewController()
	exchangeCtrl.LoadState(document.ExchangeState)
	c.exchangeCtrlMap[id] = exchangeCtrl

	c.model.setDocumentData(id, documentData{
		collectionName: document.CollectionName,
		requestName:    document.RequestName,
	})

	c.view.addDocumentTab(id, document.CollectionName, document.RequestName, exchangeCtrl.GetUI())
}

func (c *Controller) SaveTab(callback func()) error {
	id := c.view.getSelectedDocumentId()
	documentData := c.model.getDocumentData(id)
	exchangeState := c.exchangeCtrlMap[id].ToState()

	document := DocumentState{
		CollectionName: documentData.collectionName,
		RequestName:    documentData.requestName,
		ExchangeState:  exchangeState}
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
