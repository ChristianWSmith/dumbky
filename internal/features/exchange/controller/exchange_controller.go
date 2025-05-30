package controller

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/exchange/model"
	"dumbky/internal/features/exchange/view"
	"dumbky/internal/features/request"
	"dumbky/internal/features/response"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/restclient"
	"dumbky/internal/state"
	"dumbky/internal/utils"
	"dumbky/internal/validators"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type controllerImpl struct {
	model model.Model
	view  view.View

	requestCtrl  request.RequestController
	responseCtrl response.ResponseController
}

func NewController(
	requestCtrl request.RequestController,
	responseCtrl response.ResponseController) *controllerImpl {

	c := &controllerImpl{
		model:        model.NewModel(),
		view:         view.NewView(requestCtrl.CanvasObject(), responseCtrl.CanvasObject()),
		requestCtrl:  requestCtrl,
		responseCtrl: responseCtrl,
	}

	c.bindAll()

	c.view.SetUrlValidator(validators.ValidateURL)

	c.model.SetMethodListener(func() {
		method := c.model.GetMethod()
		if method == constants.HTTP_METHOD_GET ||
			method == constants.HTTP_METHOD_HEAD {
			requestCtrl.SetBodyTypeSelectEnabled(false)
		} else if method == constants.HTTP_METHOD_DELETE ||
			method == constants.HTTP_METHOD_OPTIONS ||
			method == constants.HTTP_METHOD_PATCH ||
			method == constants.HTTP_METHOD_POST ||
			method == constants.HTTP_METHOD_PUT {
			requestCtrl.SetBodyTypeSelectEnabled(true)
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
		Method:  c.model.GetMethod(),
		URL:     c.model.GetURL(),
		UseSSL:  c.model.GetUseSSL(),
		Request: c.requestCtrl.ToState(),
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
	c.requestCtrl.LoadState(exchangeState.Request)
}

func (c *controllerImpl) Validate() error {
	err := c.requestCtrl.Validate()
	if err != nil {
		return err
	}
	return c.view.ValidateURL()
}

func (c *controllerImpl) bindAll() {
	bindings := c.model.GetBindings()
	bindables := c.view.GetBindables()
	bindables.Method.Bind(bindings.Method)
	bindables.URL.Bind(bindings.URL)
	bindables.UseSSL.Bind(bindings.UseSSL)
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

	requestPayload, err := c.renderRequestConfig()
	if err != nil {
		// TODO: error feedback
		log.Error(err)
		c.setLoading(false)
		return
	}

	go c.sendRequestWorker(requestPayload)
}

func (c *controllerImpl) sendRequestWorker(requestConfig restclient.RequestConfig) {
	defer fyne.Do(func() {
		c.setLoading(false)
	})

	go func() {
		fyne.Do(func() {
			bodyType := c.requestCtrl.GetBodyType()
			if bodyType != constants.UI_BODY_TYPE_RAW {
				return
			}
			c.requestCtrl.FormatBodyRaw()
		})
	}()

	responsePayload, err := restclient.SendRequest(requestConfig)
	if err != nil {
		log.Warn(err)
		dialog.ShowError(err, global.Window)
		return
	}

	fyne.Do(func() {
		c.responseCtrl.Set(responsePayload)
	})
}

func (c *controllerImpl) renderRequestConfig() (restclient.RequestConfig, error) {
	err := c.Validate()
	if err != nil {
		log.Warn(err)
		return restclient.RequestConfig{}, err
	}

	url := c.model.GetURL()
	method := c.model.GetMethod()
	useSSL := c.model.GetUseSSL()

	headers := c.requestCtrl.GetHeadersMap()
	queryParams := c.requestCtrl.GetQueryParamsMap()
	pathParams := c.requestCtrl.GetPathParamsMap()
	bodyType := c.requestCtrl.GetBodyType()
	bodyRaw := c.requestCtrl.GetBodyRaw()
	bodyForm := c.requestCtrl.GetBodyFormMap()

	return restclient.RequestConfig{
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
