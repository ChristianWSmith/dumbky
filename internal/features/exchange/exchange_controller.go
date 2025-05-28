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
	"dumbky/internal/utils"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type Controller struct {
	model *model
	view  *view

	exchangeHeaderCtrl *exchangeheader.Controller
	requestCtrl        *request.Controller
	responseCtrl       *response.Controller
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	exchangeHeaderCtrl := exchangeheader.NewController()
	requestCtrl := request.NewController()
	responseCtrl := response.NewController()

	c := &Controller{
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
			requestCtrl.SetRequestBodyType(constants.UI_BODY_TYPE_NONE)
			requestCtrl.DisableRequestBodyTypeSelect()
		} else if method == constants.HTTP_METHOD_DELETE ||
			method == constants.HTTP_METHOD_OPTIONS ||
			method == constants.HTTP_METHOD_PATCH ||
			method == constants.HTTP_METHOD_POST ||
			method == constants.HTTP_METHOD_PUT {
			requestCtrl.EnableRequestBodyTypeSelect()
		} else {
			log.Error(errors.New("invalid http method"))
		}
	})

	exchangeHeaderCtrl.SetSendHandler(func() {
		c.sendButtonHandler()
	})

	return c
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) ToState() ExchangeState {
	return ExchangeState{
		Header:  c.exchangeHeaderCtrl.ToState(),
		Request: c.requestCtrl.ToState(),
	}
}

func (c *Controller) LoadState(exchangeState ExchangeState) {
	c.requestCtrl.LoadState(exchangeState.Request)
	c.exchangeHeaderCtrl.LoadState(exchangeState.Header)
}

func (c *Controller) sendButtonHandler() {
	c.exchangeHeaderCtrl.DisableSend()
	c.responseCtrl.SetLoading(true)

	c.responseCtrl.SetStatus(constants.UI_LOADING_RESPONSE_STATUS)
	c.responseCtrl.SetTime(constants.UI_LOADING_RESPONSE_TIME)
	c.responseCtrl.SetBody(constants.UI_LOADING_RESPONSE_BODY)

	requestPayload := c.renderRequestConfig()

	go c.sendRequestWorker(requestPayload)
}

func (c *Controller) sendRequestWorker(requestConfig requesthelper.RequestConfig) {
	defer fyne.Do(func() {
		c.responseCtrl.SetLoading(false)
		c.exchangeHeaderCtrl.EnableSend()
	})

	fyne.Do(func() {
		bodyType := c.requestCtrl.GetRequestBodyType()
		if bodyType != constants.UI_BODY_TYPE_RAW {
			return
		}
		bodyRaw := c.requestCtrl.GetRequestBodyRaw()
		c.requestCtrl.SetRequestBodyRaw(utils.SmartFormat(bodyRaw))
	})

	responsePayload, err := requesthelper.SendRequest(requestConfig)
	if err != nil {
		log.Warn(err)
		dialog.ShowError(err, global.Window)
		return
	}

	fyne.Do(func() {
		c.responseCtrl.SetStatus(responsePayload.Status)
		c.responseCtrl.SetTime(responsePayload.Time)
		c.responseCtrl.SetBody(utils.SmartFormat(responsePayload.Body))
	})
}

func (c *Controller) renderRequestConfig() requesthelper.RequestConfig {
	url := c.exchangeHeaderCtrl.GetURL()
	err := c.exchangeHeaderCtrl.ValidateURL()
	if err != nil {
		log.Warn(err)
	}

	method := c.exchangeHeaderCtrl.GetMethod()

	useSSL := c.exchangeHeaderCtrl.GetUseSSL()

	headers := c.requestCtrl.GetHeadersMap()

	err = c.requestCtrl.ValidateHeaders()
	if err != nil {
		log.Error(err)
	}

	queryParams := c.requestCtrl.GetQueryParamsMap()
	err = c.requestCtrl.ValidateQueryParams()
	if err != nil {
		log.Error(err)
	}

	pathParams := c.requestCtrl.GetPathParamsMap()
	err = c.requestCtrl.ValidatePathParams()
	if err != nil {
		log.Error(err)
	}

	bodyType := c.requestCtrl.GetRequestBodyType()
	bodyRaw := c.requestCtrl.GetRequestBodyRaw()

	err = c.requestCtrl.ValidateRequestBodyRaw()
	if err != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Warn(err)
	}

	bodyForm := c.requestCtrl.GetRequestBodyFormMap()
	err = c.requestCtrl.ValidateRequestBodyForm()
	if err != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn(err)
	}

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
	}
}
