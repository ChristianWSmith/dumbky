package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/features/request"
	"dumbky/internal/features/response"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/requesthelper"
	"dumbky/internal/utils"
	"dumbky/internal/validators"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type controller struct {
	model *model
	view  *view

	requestCtrl  request.RequestController
	responseCtrl response.ResponseController
}

type ExchangeController interface {
	features.Controller
	ToState() ExchangeState
	LoadState(exchangeState ExchangeState)
	Validate() error
}

var _ ExchangeController = (*controller)(nil)

func New() ExchangeController {
	requestCtrl := request.New()
	responseCtrl := response.New()
	return newController(requestCtrl, responseCtrl)
}

func newController(
	requestCtrl request.RequestController,
	responseCtrl response.ResponseController) *controller {

	c := &controller{
		model:        newModel(),
		view:         newView(requestCtrl.CanvasObject(), responseCtrl.CanvasObject()),
		requestCtrl:  requestCtrl,
		responseCtrl: responseCtrl,
	}

	c.view.methodSelect.Bind(c.model.methodBinding)
	c.view.urlEntry.Bind(c.model.urlBinding)
	c.view.sslCheck.Bind(c.model.useSSLBinding)

	c.view.setUrlValidator(validators.ValidateURL)

	c.model.setMethodListener(func() {
		method := c.model.getMethod()
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

	c.view.setSendHandler(func() {
		c.sendButtonHandler()
	})

	return c
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) ToState() ExchangeState {
	return ExchangeState{
		Method:  c.model.getMethod(),
		URL:     c.model.getURL(),
		UseSSL:  c.model.getUseSSL(),
		Request: c.requestCtrl.ToState(),
	}
}

func (c *controller) LoadState(exchangeState ExchangeState) {
	method := exchangeState.Method
	if !utils.ElementInSlice(constants.HttpMethods(), method) {
		method = constants.HTTP_METHOD_DEFAULT
	}
	c.model.setMethod(method)
	c.model.setURL(exchangeState.URL)
	c.model.setUseSSL(exchangeState.UseSSL)
	c.requestCtrl.LoadState(exchangeState.Request)
}

func (c *controller) Validate() error {
	err := c.requestCtrl.Validate()
	if err != nil {
		return err
	}
	return c.view.validateURL()
}

func (c *controller) setSendEnabled(enabled bool) {
	if enabled {
		c.view.enableSend()
	} else {
		c.view.disableSend()
	}
}

func (c *controller) setLoading(loading bool) {
	c.setSendEnabled(!loading)
	c.responseCtrl.SetLoading(loading)
}

func (c *controller) sendButtonHandler() {
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

func (c *controller) sendRequestWorker(requestConfig requesthelper.RequestConfig) {
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

	responsePayload, err := requesthelper.SendRequest(requestConfig)
	if err != nil {
		log.Warn(err)
		dialog.ShowError(err, global.Window)
		return
	}

	fyne.Do(func() {
		c.responseCtrl.Set(responsePayload)
	})
}

func (c *controller) renderRequestConfig() (requesthelper.RequestConfig, error) {
	err := c.Validate()
	if err != nil {
		log.Warn(err)
		return requesthelper.RequestConfig{}, err
	}

	url := c.model.getURL()
	method := c.model.getMethod()
	useSSL := c.model.getUseSSL()

	headers := c.requestCtrl.GetHeadersMap()
	queryParams := c.requestCtrl.GetQueryParamsMap()
	pathParams := c.requestCtrl.GetPathParamsMap()
	bodyType := c.requestCtrl.GetBodyType()
	bodyRaw := c.requestCtrl.GetBodyRaw()
	bodyForm := c.requestCtrl.GetBodyFormMap()

	return requesthelper.RequestConfig{
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
