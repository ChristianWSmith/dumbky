package exchangeview

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/exchangeheader"
	"dumbky/internal/features/request"
	"dumbky/internal/features/response"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/requesthelper"
	"dumbky/internal/utils"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

type ExchangeView struct {
	UI                 *fyne.Container
	exchangeHeaderCtrl *exchangeheader.Controller
	requestCtrl        *request.Controller
	responseCtrl       *response.Controller
}

type ExchangeState struct {
	Header  exchangeheader.ExchangeHeaderState `json:"header"`
	Request request.RequestState               `json:"request"`
}

func (ev ExchangeView) ToState() (ExchangeState, error) {
	header := ev.exchangeHeaderCtrl.ToState()
	request, requestErr := ev.requestCtrl.ToState()
	if requestErr != nil {
		log.Error(requestErr)
		return ExchangeState{}, requestErr
	}
	return ExchangeState{
		Header:  header,
		Request: request,
	}, nil
}

func (ev ExchangeView) LoadState(exchangeState ExchangeState) error {
	requestErr := ev.requestCtrl.LoadState(exchangeState.Request)
	if requestErr != nil {
		log.Error(requestErr)
		return requestErr
	}
	ev.exchangeHeaderCtrl.LoadState(exchangeState.Header)
	return nil
}

func (ev ExchangeView) ToRequestPayload() requesthelper.RequestPayload {
	url := ev.exchangeHeaderCtrl.GetURL()
	urlValidateErr := ev.exchangeHeaderCtrl.ValidateURL()
	if urlValidateErr != nil {
		log.Warn(urlValidateErr)
	}

	method := ev.exchangeHeaderCtrl.GetMethod()

	useSSL := ev.exchangeHeaderCtrl.GetUseSSL()

	headers := ev.requestCtrl.GetHeadersMap()
	headersValidatErr := ev.requestCtrl.ValidateHeaders()
	if headersValidatErr != nil {
		log.Error(headersValidatErr)
	}

	queryParams := ev.requestCtrl.GetQueryParamsMap()
	queryParamsValidatErr := ev.requestCtrl.ValidateQueryParams()
	if queryParamsValidatErr != nil {
		log.Error(queryParamsValidatErr)
	}

	pathParams := ev.requestCtrl.GetPathParamsMap()
	pathParamsValidatErr := ev.requestCtrl.ValidatePathParams()
	if pathParamsValidatErr != nil {
		log.Error(pathParamsValidatErr)
	}

	bodyType := ev.requestCtrl.GetRequestBodyType()
	bodyRaw := ev.requestCtrl.GetRequestBodyRaw()

	bodyRawValidateErr := ev.requestCtrl.ValidateRequestBodyRaw()
	if bodyRawValidateErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Warn(bodyRawValidateErr)
	}

	bodyForm := ev.requestCtrl.GetRequestBodyFormMap()
	bodyFormValidateErr := ev.requestCtrl.ValidateRequestBodyForm()
	if bodyFormValidateErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn(bodyFormValidateErr)
	}

	return requesthelper.RequestPayload{
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

func (ev ExchangeView) sendRequestWorker(requestPayload requesthelper.RequestPayload) {
	defer fyne.Do(func() {
		ev.responseCtrl.SetLoading(false)
		ev.exchangeHeaderCtrl.EnableSend()
	})

	fyne.Do(func() {
		bodyType := ev.requestCtrl.GetRequestBodyType()
		if bodyType != constants.UI_BODY_TYPE_RAW {
			return
		}
		bodyRaw := ev.requestCtrl.GetRequestBodyRaw()
		ev.requestCtrl.SetRequestBodyRaw(utils.SmartFormat(bodyRaw))
	})

	responsePayload, err := requesthelper.SendRequest(requestPayload)
	if err != nil {
		log.Warn(err)
		dialog.ShowError(err, global.Window)
		return
	}

	fyne.Do(func() {
		ev.responseCtrl.SetStatus(responsePayload.Status)
		ev.responseCtrl.SetTime(responsePayload.Time)
		ev.responseCtrl.SetBody(utils.SmartFormat(responsePayload.Body))
	})
}

func (ev ExchangeView) sendButtonHandler() {
	ev.exchangeHeaderCtrl.DisableSend()
	ev.responseCtrl.SetLoading(true)

	ev.responseCtrl.SetStatus(constants.UI_LOADING_RESPONSE_STATUS)
	ev.responseCtrl.SetTime(constants.UI_LOADING_RESPONSE_TIME)
	ev.responseCtrl.SetBody(constants.UI_LOADING_RESPONSE_BODY)

	requestPayload := ev.ToRequestPayload()

	go ev.sendRequestWorker(requestPayload)
}

func ComposeExchangeView() ExchangeView {
	exchangeHeaderCtrl := exchangeheader.NewController()
	requestCtrl := request.NewController()
	responseCtrl := response.NewController()

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

	requestResponseView := container.NewHSplit(requestCtrl.GetUI(), responseCtrl.GetUI())
	ui := container.NewBorder(exchangeHeaderCtrl.GetUI(), nil, nil, nil, requestResponseView)

	ev := ExchangeView{
		ui,
		exchangeHeaderCtrl,
		requestCtrl,
		responseCtrl,
	}

	exchangeHeaderCtrl.SetSendHandler(func() {
		ev.sendButtonHandler()
	})

	return ev
}
