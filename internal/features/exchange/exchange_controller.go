package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/features/exchangeheader"
	"dumbky/internal/features/request"
	"dumbky/internal/features/response"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/requesthelper"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type controller struct {
	model *model
	view  *view

	exchangeHeaderCtrl exchangeheader.ExchangeHeaderController
	requestCtrl        request.RequestController
	responseCtrl       response.ResponseController
}

type ExchangeController interface {
	features.Controller
	ToState() ExchangeState
	LoadState(exchangeState ExchangeState)
}

var _ ExchangeController = (*controller)(nil)

func New() ExchangeController {
	exchangeHeaderCtrl := exchangeheader.NewController()
	requestCtrl := request.New()
	responseCtrl := response.New()
	return newController(exchangeHeaderCtrl, requestCtrl, responseCtrl)
}

func newController(
	exchangeHeaderCtrl exchangeheader.ExchangeHeaderController,
	requestCtrl request.RequestController,
	responseCtrl response.ResponseController) *controller {

	c := &controller{
		model:              newModel(),
		view:               newView(exchangeHeaderCtrl.CanvasObject(), requestCtrl.CanvasObject(), responseCtrl.CanvasObject()),
		exchangeHeaderCtrl: exchangeHeaderCtrl,
		requestCtrl:        requestCtrl,
		responseCtrl:       responseCtrl,
	}

	exchangeHeaderCtrl.SetMethodListener(func() {
		method := exchangeHeaderCtrl.GetMethod()
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
			log.Error(errors.New("invalid http method"))
		}
	})

	exchangeHeaderCtrl.SetSendHandler(func() {
		c.sendButtonHandler()
	})

	return c
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) ToState() ExchangeState {
	return ExchangeState{
		Header:  c.exchangeHeaderCtrl.ToState(),
		Request: c.requestCtrl.ToState(),
	}
}

func (c *controller) LoadState(exchangeState ExchangeState) {
	c.requestCtrl.LoadState(exchangeState.Request)
	c.exchangeHeaderCtrl.LoadState(exchangeState.Header)
}

func (c *controller) setLoading(loading bool) {
	c.exchangeHeaderCtrl.SetSendEnabled(!loading)
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
			bodyType := c.requestCtrl.GetRequestBodyType()
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
	err := c.exchangeHeaderCtrl.Validate()
	if err != nil {
		log.Warn(err)
		return requesthelper.RequestConfig{}, err
	}
	err = c.requestCtrl.Validate()
	if err != nil {
		log.Error(err)
		return requesthelper.RequestConfig{}, err
	}

	url := c.exchangeHeaderCtrl.GetURL()
	method := c.exchangeHeaderCtrl.GetMethod()
	useSSL := c.exchangeHeaderCtrl.GetUseSSL()

	headers := c.requestCtrl.GetHeadersMap()
	queryParams := c.requestCtrl.GetQueryParamsMap()
	pathParams := c.requestCtrl.GetPathParamsMap()
	bodyType := c.requestCtrl.GetRequestBodyType()
	bodyRaw := c.requestCtrl.GetRequestBodyRaw()
	bodyForm := c.requestCtrl.GetRequestBodyFormMap()

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
