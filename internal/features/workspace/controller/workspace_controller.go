package controller

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/events"
	"dumbky/internal/features/exchange"
	"dumbky/internal/features/workspace/common"
	"dumbky/internal/features/workspace/model"
	"dumbky/internal/features/workspace/view"
	"dumbky/internal/httputils"
	"dumbky/internal/log"
	"dumbky/internal/state"
	"dumbky/internal/utils"
	"strings"

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
			c.OpenTab(state.DocumentState{
				CollectionName: constants.DB_DEFAULT_COLLECTION_NAME,
				RequestName:    utils.SillyName()})
		}
	})

	c.view.SetImportHandler(func() {
		// TODO: split this out, do body
		curlCmd := `curl -X POST -H "Content-Type: application/json" -d '{"key":"value"}' https://example.com/api?a=b`

		request, _ := httputils.CurlToRequest(curlCmd)
		exchangeState := state.ExchangeState{}
		prefix := ""
		if strings.HasPrefix(request.URL.String(), "http://") {
			exchangeState.UseSSL = false
			prefix = "http://"
		} else if strings.HasPrefix(request.URL.String(), "https://") {
			exchangeState.UseSSL = true
			prefix = "https://"
		}
		exchangeState.Method = request.Method
		exchangeState.URL = prefix + request.URL.Host + request.URL.Path
		for key, value := range request.URL.Query() {
			exchangeState.Request.QueryParams.KeyValueStates =
				append(exchangeState.Request.Headers.KeyValueStates,
					state.KeyValueState{Enabled: true, Key: key, Value: strings.Join(value, ",")})
		}
		for key, value := range request.Header {
			exchangeState.Request.Headers.KeyValueStates =
				append(exchangeState.Request.Headers.KeyValueStates,
					state.KeyValueState{Enabled: true, Key: key, Value: strings.Join(value, ",")})
		}
		c.OpenTab(state.DocumentState{
			CollectionName: c.model.GetCollectionName(),
			RequestName:    "IMPORT " + utils.SillyName(),
			ExchangeState:  exchangeState,
		})
	})

	c.view.SetExportHandler(func() {
		// TODO: split this out
		id := c.view.GetSelectedDocumentId()
		exchangeCtrl := c.exchangeCtrlMap[id]
		requestConfig, _ := exchangeCtrl.RenderRequestConfig()
		request, _ := httputils.RenderRequest(requestConfig)
		curlCmd, _ := httputils.RequestToCurl(request)
		log.Debug(curlCmd) // TODO: display this with a copy to clipboard button
	})

	c.OpenTab(state.DocumentState{
		CollectionName: constants.DB_DEFAULT_COLLECTION_NAME,
		RequestName:    utils.SillyName()})

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

	c.SetSaveHandler(func() {
		go c.SaveTab()
	})

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) OpenTab(document state.DocumentState) {
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

func (c *controllerImpl) SaveTab() error {
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

func (c *controllerImpl) LoadTab(collectionName, requestName string) {
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
			c.OpenTab(document)
		})
	}()
}

func (c *controllerImpl) SetAddHandler(handler func()) {
	c.view.SetAddHandler(handler)
}

func (c *controllerImpl) SetSaveHandler(handler func()) {
	c.view.SetSaveHandler(handler)
}

func (c *controllerImpl) bindAll() {
	bindings := c.model.GetBindings()
	bindables := c.view.GetBindables()
	bindables.RequestName.Bind(bindings.RequestName)
	bindables.CollectionName.Bind(bindings.CollectionName)
}
