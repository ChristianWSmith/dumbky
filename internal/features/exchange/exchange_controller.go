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

func (c *Controller) setLoading(loading bool) {
	if loading {
		c.exchangeHeaderCtrl.DisableSend()
	} else {
		c.exchangeHeaderCtrl.EnableSend()
	}
	c.responseCtrl.SetLoading(loading)
}

func (c *Controller) sendButtonHandler() {
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

func (c *Controller) sendRequestWorker(requestConfig requesthelper.RequestConfig) {
	defer fyne.Do(func() {
		c.setLoading(false)
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
		c.responseCtrl.SetResponse(
			responsePayload.Status,
			responsePayload.Time,
			utils.SmartFormat(responsePayload.Body))
	})
}

func (c *Controller) renderRequestConfig() (requesthelper.RequestConfig, error) {
	url := c.exchangeHeaderCtrl.GetURL()
	err := c.exchangeHeaderCtrl.ValidateURL()
	if err != nil {
		log.Warn(err)
		return requesthelper.RequestConfig{}, err
	}

	method := c.exchangeHeaderCtrl.GetMethod()

	useSSL := c.exchangeHeaderCtrl.GetUseSSL()

	headers := c.requestCtrl.GetHeadersMap()

	err = c.requestCtrl.ValidateHeaders()
	if err != nil {
		log.Error(err)
		return requesthelper.RequestConfig{}, err
	}

	queryParams := c.requestCtrl.GetQueryParamsMap()
	err = c.requestCtrl.ValidateQueryParams()
	if err != nil {
		log.Error(err)
		return requesthelper.RequestConfig{}, err
	}

	pathParams := c.requestCtrl.GetPathParamsMap()
	err = c.requestCtrl.ValidatePathParams()
	if err != nil {
		log.Error(err)
		return requesthelper.RequestConfig{}, err
	}

	bodyType := c.requestCtrl.GetRequestBodyType()
	bodyRaw := c.requestCtrl.GetRequestBodyRaw()

	err = c.requestCtrl.ValidateRequestBodyRaw()
	if err != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Warn(err)
		return requesthelper.RequestConfig{}, err
	}

	bodyForm := c.requestCtrl.GetRequestBodyFormMap()
	err = c.requestCtrl.ValidateRequestBodyForm()
	if err != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn(err)
		return requesthelper.RequestConfig{}, err
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
	}, nil
}
