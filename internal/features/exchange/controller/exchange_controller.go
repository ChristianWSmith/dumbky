package controller

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/exchange/model"
	"dumbky/internal/features/exchange/view"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/features/response"
	"dumbky/internal/global"
	"dumbky/internal/httputils"
	"dumbky/internal/log"
	"dumbky/internal/state"
	"dumbky/internal/utils"
	"dumbky/internal/validators"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type controllerImpl struct {
	model model.Model
	view  view.View

	responseCtrl response.ResponseController

	queryParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController
	pathParamsKeyValueCtrl  keyvalueeditor.KeyValueEditorController
	headersKeyValueCtrl     keyvalueeditor.KeyValueEditorController
	bodyFormKeyValueCtrl    keyvalueeditor.KeyValueEditorController
}

func NewController(
	queryParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	pathParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	headersKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	bodyFormKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	responseCtrl response.ResponseController) *controllerImpl {

	c := &controllerImpl{
		model:                   model.NewModel(),
		view:                    view.NewView(queryParamsKeyValueCtrl.CanvasObject(), pathParamsKeyValueCtrl.CanvasObject(), headersKeyValueCtrl.CanvasObject(), bodyFormKeyValueCtrl.CanvasObject(), responseCtrl.CanvasObject()),
		responseCtrl:            responseCtrl,
		queryParamsKeyValueCtrl: queryParamsKeyValueCtrl,
		pathParamsKeyValueCtrl:  pathParamsKeyValueCtrl,
		headersKeyValueCtrl:     headersKeyValueCtrl,
		bodyFormKeyValueCtrl:    bodyFormKeyValueCtrl,
	}

	c.bindAll()
	c.model.SetBodyTypeListener(c.showBodyType)

	c.view.SetUrlValidator(validators.ValidateURL)

	c.model.SetMethodListener(func() {
		method := c.model.GetMethod()
		if method == constants.HTTP_METHOD_GET ||
			method == constants.HTTP_METHOD_HEAD {
			c.SetBodyTypeSelectEnabled(false)
		} else if method == constants.HTTP_METHOD_DELETE ||
			method == constants.HTTP_METHOD_OPTIONS ||
			method == constants.HTTP_METHOD_PATCH ||
			method == constants.HTTP_METHOD_POST ||
			method == constants.HTTP_METHOD_PUT {
			c.SetBodyTypeSelectEnabled(true)
		} else {
			log.Error(fmt.Errorf("invalid http method %s", method))
		}
	})

	c.view.SetSendHandler(func() {
		c.sendButtonHandler()
	})

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) ToState() state.ExchangeState {
	return state.ExchangeState{
		Method:      c.model.GetMethod(),
		URL:         c.model.GetURL(),
		UseSSL:      c.model.GetUseSSL(),
		QueryParams: c.queryParamsKeyValueCtrl.ToState(),
		PathParams:  c.pathParamsKeyValueCtrl.ToState(),
		Headers:     c.headersKeyValueCtrl.ToState(),
		BodyForm:    c.bodyFormKeyValueCtrl.ToState(),
		BodyType:    c.model.GetBodyType(),
		BodyRaw:     c.model.GetBodyRaw(),
	}
}

func (c *controllerImpl) LoadState(exchangeState state.ExchangeState) {
	method := exchangeState.Method
	if !utils.ElementInSlice(constants.HttpMethods(), method) {
		method = constants.HTTP_METHOD_DEFAULT
	}
	c.model.SetMethod(method)
	c.model.SetURL(exchangeState.URL)
	c.model.SetUseSSL(exchangeState.UseSSL)
	c.queryParamsKeyValueCtrl.LoadState(exchangeState.QueryParams)
	c.pathParamsKeyValueCtrl.LoadState(exchangeState.PathParams)
	c.headersKeyValueCtrl.LoadState(exchangeState.Headers)
	c.bodyFormKeyValueCtrl.LoadState(exchangeState.BodyForm)
	c.model.SetBodyType(exchangeState.BodyType)
	c.model.SetBodyRaw(exchangeState.BodyRaw)
}

func (c *controllerImpl) Validate() error {
	err := c.queryParamsKeyValueCtrl.Validate()
	if err != nil {
		return err
	}
	err = c.pathParamsKeyValueCtrl.Validate()
	if err != nil {
		return err
	}
	err = c.headersKeyValueCtrl.Validate()
	if err != nil {
		return err
	}
	err = c.bodyFormKeyValueCtrl.Validate()
	if err != nil {
		return err
	}
	return c.view.Validate()
}

func (c *controllerImpl) bindAll() {
	bindings := c.model.GetBindings()
	bindables := c.view.GetBindables()
	bindables.Method.Bind(bindings.Method)
	bindables.URL.Bind(bindings.URL)
	bindables.UseSSL.Bind(bindings.UseSSL)
	bindables.BodyRaw.Bind(bindings.BodyRaw)
	bindables.BodyType.Bind(bindings.BodyType)
}

func (c *controllerImpl) setSendEnabled(enabled bool) {
	if enabled {
		c.view.EnableSend()
	} else {
		c.view.DisableSend()
	}
}

func (c *controllerImpl) setLoading(loading bool) {
	c.setSendEnabled(!loading)
	c.responseCtrl.SetLoading(loading)
}

func (c *controllerImpl) sendButtonHandler() {
	c.setLoading(true)

	requestPayload, err := c.RenderRequestConfig()
	if err != nil {
		// TODO: error feedback
		log.Error(err)
		c.setLoading(false)
		return
	}

	go c.sendRequestWorker(requestPayload)
}

func (c *controllerImpl) sendRequestWorker(requestConfig httputils.RequestConfig) {
	defer fyne.Do(func() {
		c.setLoading(false)
	})

	go func() {
		fyne.Do(func() {
			bodyType := c.model.GetBodyType()
			if bodyType != constants.UI_BODY_TYPE_RAW {
				return
			}
			c.FormatBodyRaw()
		})
	}()

	responsePayload, err := httputils.SendRequest(requestConfig)
	if err != nil {
		log.Warn(err)
		dialog.ShowError(err, global.Window)
		return
	}

	fyne.Do(func() {
		c.responseCtrl.Set(responsePayload)
	})
}

func (c *controllerImpl) RenderRequestConfig() (httputils.RequestConfig, error) {
	err := c.Validate()
	if err != nil {
		log.Warn(err)
		return httputils.RequestConfig{}, err
	}

	url := c.model.GetURL()
	method := c.model.GetMethod()
	useSSL := c.model.GetUseSSL()

	headers := c.headersKeyValueCtrl.Get()
	queryParams := c.queryParamsKeyValueCtrl.Get()
	pathParams := c.pathParamsKeyValueCtrl.Get()
	bodyType := c.model.GetBodyType()
	bodyRaw := c.model.GetBodyRaw()
	bodyForm := c.bodyFormKeyValueCtrl.Get()

	return httputils.RequestConfig{
		URL:         url,
		Method:      method,
		UseSSL:      useSSL,
		Headers:     headers,
		QueryParams: queryParams,
		PathParams:  pathParams,
		BodyType:    bodyType,
		BodyRaw:     bodyRaw,
		BodyForm:    bodyForm,
	}, nil
}

func (c *controllerImpl) SetBodyTypeSelectEnabled(enabled bool) {
	if enabled {
		c.view.SetBodyTypeSelectEnabled(true)
	} else {
		c.view.SetBodyTypeSelectEnabled(false)
		c.model.SetBodyType(constants.UI_BODY_TYPE_NONE)
	}
}

func (c *controllerImpl) FormatBodyRaw() {
	c.model.SetBodyRaw(utils.SmartFormat(c.model.GetBodyRaw()))
}

func (c *controllerImpl) showBodyType() {
	bodyType := c.model.GetBodyType()
	if bodyType == constants.UI_BODY_TYPE_FORM {
		c.bodyFormKeyValueCtrl.SetVisible(true)
		c.view.SetBodyRawVisible(false)
	} else if bodyType == constants.UI_BODY_TYPE_RAW {
		c.bodyFormKeyValueCtrl.SetVisible(false)
		c.view.SetBodyRawVisible(true)
	} else if bodyType == constants.UI_BODY_TYPE_NONE {
		c.bodyFormKeyValueCtrl.SetVisible(false)
		c.view.SetBodyRawVisible(false)
	} else {
		log.Error(errors.New("invalid body type"))
	}
}
