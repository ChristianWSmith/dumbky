package controller

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/events"
	"dumbky/internal/features/exchange"
	"dumbky/internal/features/workspace/common"
	"dumbky/internal/features/workspace/model"
	"dumbky/internal/features/workspace/view"
	"dumbky/internal/log"
	"dumbky/internal/state"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/google/uuid"
)

type controllerImpl struct {
	model           model.Model
	view            view.View
	exchangeCtrlMap map[common.DocumentId]exchange.ExchangeController
}

func NewController() *controllerImpl {
	c := &controllerImpl{
		model:           model.NewModel(),
		view:            view.NewView(),
		exchangeCtrlMap: make(map[common.DocumentId]exchange.ExchangeController),
	}

	c.bindAll()

	c.view.SetDocumentSelectedHandler(func(tabItem *container.TabItem) {
		tabId := c.view.GetDocumentId(tabItem)
		documentData := c.model.GetDocumentData(tabId)
		c.model.SetRequestName(documentData.RequestName)
		c.model.SetCollectionName(documentData.CollectionName)
	})

	c.view.SetDocumentClosedHandler(func(tabItem *container.TabItem) {
		id := c.view.GetDocumentId(tabItem)
		c.view.DestroyDocumentTab(tabItem)
		c.model.DestroyDocumentData(id)
		delete(c.exchangeCtrlMap, id)
		if c.model.DocumentCount() == 0 {
			c.openTab(state.DocumentState{
				CollectionName: constants.DB_DEFAULT_COLLECTION_NAME,
				RequestName:    utils.SillyName()})
		}
	})

	c.view.SetImportHandler(func() {
		// TODO
	})

	c.view.SetExportHandler(func() {
		// TODO
	})

	c.model.SetRequestNameListener(func() {
		id := c.view.GetSelectedDocumentId()
		c.model.UpdateRequestName(id, c.model.GetRequestName())
		documentData := c.model.GetDocumentData(id)
		if documentData.RequestName == "" {
			c.view.SetSelectedDocumentText(constants.UI_PLACEHOLDER_UNTITLED)
		} else {
			c.view.SetSelectedDocumentText(documentData.RequestName)
		}
		c.view.RefreshTabs()

	})

	c.view.SetSaveHandler(func() {
		go c.saveTab()
	})

	events.Subscribe(func(requestSelectedEvent events.RequestSelected) {
		if !requestSelectedEvent.IsSelected {
			return
		}
		collectionSelectedEvent := events.Current[events.CollectionSelected]()
		if !collectionSelectedEvent.IsSelected {
			return
		}
		c.loadTab(collectionSelectedEvent.CollectionName, requestSelectedEvent.RequestName)
	})

	c.view.SetAddHandler(func() {
		collectionSelectedEvent := events.Current[events.CollectionSelected]()
		collectionName := collectionSelectedEvent.CollectionName
		if !collectionSelectedEvent.IsSelected {
			collectionName = constants.DB_DEFAULT_COLLECTION_NAME
		}
		c.openTab(state.DocumentState{
			CollectionName: collectionName,
			RequestName:    utils.SillyName()})
	})

	c.openTab(state.DocumentState{
		CollectionName: constants.DB_DEFAULT_COLLECTION_NAME,
		RequestName:    utils.SillyName()})

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) openTab(document state.DocumentState) {
	if c.view.SelectDocumentTabOnCondition(func(id common.DocumentId) bool {
		documentData := c.model.GetDocumentData(id)
		return documentData.RequestName == document.RequestName && documentData.CollectionName == document.CollectionName
	}) {
		return
	}

	id := common.DocumentId(uuid.New().String())

	exchangeCtrl := exchange.New()
	exchangeCtrl.LoadState(document.ExchangeState)
	c.exchangeCtrlMap[id] = exchangeCtrl

	c.model.SetDocumentData(id, common.DocumentData{
		CollectionName: document.CollectionName,
		RequestName:    document.RequestName,
	})

	c.view.AddDocumentTab(id, document.RequestName, exchangeCtrl.CanvasObject())
}

func (c *controllerImpl) saveTab() error {
	id := c.view.GetSelectedDocumentId()
	documentData := c.model.GetDocumentData(id)
	exchangeState := c.exchangeCtrlMap[id].ToState()

	document := state.DocumentState{
		CollectionName: documentData.CollectionName,
		RequestName:    documentData.RequestName,
		ExchangeState:  exchangeState}
	request, err := model.DocumentStateToRequest(document)
	if err != nil {
		log.Error(err)
		return err
	}
	go func() {
		err := db.SaveRequest(request)
		if err != nil {
			log.Error(err)
			return
		}
		events.Publish(events.RequestSaved{})
	}()
	return nil
}

func (c *controllerImpl) loadTab(collectionName, requestName string) {
	go func() {
		request, err := db.LoadRequest(collectionName, requestName)
		if err != nil {
			log.Error(err)
		}
		document, err := model.RequestToDocumentState(request)
		if err != nil {
			log.Error(err)
		}
		fyne.Do(func() {
			c.openTab(document)
		})
	}()
}

func (c *controllerImpl) bindAll() {
	bindings := c.model.GetBindings()
	bindables := c.view.GetBindables()
	bindables.RequestName.Bind(bindings.RequestName)
	bindables.CollectionName.Bind(bindings.CollectionName)
}
