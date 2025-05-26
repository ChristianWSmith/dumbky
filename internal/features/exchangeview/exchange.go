package exchangeview

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/exchangeheaderview"
	"dumbky/internal/features/requestview"
	"dumbky/internal/features/response"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/requesthelper"
	"dumbky/internal/utils"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
)

type ExchangeView struct {
	UI                 *fyne.Container
	headerView         exchangeheaderview.ExchangeHeaderView
	requestView        requestview.RequestView
	responseController *response.Controller
}

type ExchangeState struct {
	Header  exchangeheaderview.ExchangeHeaderState `json:"header"`
	Request requestview.RequestState               `json:"request"`
}

func (ev ExchangeView) ToState() (ExchangeState, error) {
	header, headerErr := ev.headerView.ToState()
	if headerErr != nil {
		log.Error(headerErr)
		return ExchangeState{}, headerErr
	}
	request, requestErr := ev.requestView.ToState()
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
	requestErr := ev.requestView.LoadState(exchangeState.Request)
	if requestErr != nil {
		log.Error(requestErr)
		return requestErr
	}
	headerErr := ev.headerView.LoadState(exchangeState.Header)
	if headerErr != nil {
		log.Error(headerErr)
		return headerErr
	}
	return nil
}

func (ev ExchangeView) ToRequestPayload() requesthelper.RequestPayload {
	url, urlGetErr := ev.headerView.URLBinding.Get()
	if urlGetErr != nil {
		log.Error(urlGetErr)
	}
	urlValidateErr := ev.headerView.ValidateURL()
	if urlValidateErr != nil {
		log.Warn(urlValidateErr)
	}

	method, methodGetErr := ev.headerView.MethodBinding.Get()
	if methodGetErr != nil {
		log.Error(methodGetErr)
	}

	useSSL, useSSLGetErr := ev.headerView.UseSSLBinding.Get()
	if useSSLGetErr != nil {
		log.Error(useSSLGetErr)
	}

	headers := ev.requestView.Headers.GetMap()
	headersValidatErr := ev.requestView.Headers.Validate()
	if headersValidatErr != nil {
		log.Error(headersValidatErr)
	}

	queryParams := ev.requestView.QueryParams.GetMap()
	queryParamsValidatErr := ev.requestView.QueryParams.Validate()
	if queryParamsValidatErr != nil {
		log.Error(queryParamsValidatErr)
	}

	pathParams := ev.requestView.PathParams.GetMap()
	pathParamsValidatErr := ev.requestView.PathParams.Validate()
	if pathParamsValidatErr != nil {
		log.Error(pathParamsValidatErr)
	}

	bodyType, bodyTypeGetErr := ev.requestView.Body.BodyTypeBinding.Get()
	if bodyTypeGetErr != nil {
		log.Error(bodyTypeGetErr)
	}

	bodyRaw, bodyRawGetErr := ev.requestView.Body.BodyRawBinding.Get()
	if bodyRawGetErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Error(bodyRawGetErr)
	}
	bodyRawValidateErr := ev.requestView.Body.ValidateBodyRaw()
	if bodyRawValidateErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Warn(bodyRawValidateErr)
	}

	bodyForm := ev.requestView.Body.BodyKeyValueEditor.GetMap()
	bodyFormValidateErr := ev.requestView.Body.BodyKeyValueEditor.Validate()
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
		ev.responseController.SetLoading(false)
		ev.headerView.SendButton.Enable()
	})

	fyne.Do(func() {
		bodyType, bodyTypeGetErr := ev.requestView.Body.BodyTypeBinding.Get()
		if bodyTypeGetErr != nil {
			log.Error(bodyTypeGetErr)
		}
		if bodyType != constants.UI_BODY_TYPE_RAW && bodyTypeGetErr == nil {
			return
		}
		bodyRaw, bodyRawGetErr := ev.requestView.Body.BodyRawBinding.Get()
		if bodyRawGetErr != nil {
			log.Error(bodyRawGetErr)
			return
		}
		bodyRawSetErr := ev.requestView.Body.BodyRawBinding.Set(utils.SmartFormat(bodyRaw))
		if bodyRawSetErr != nil {
			log.Error(bodyRawSetErr)
			return
		}
	})

	responsePayload, err := requesthelper.SendRequest(requestPayload)
	if err != nil {
		log.Warn(err)
		dialog.ShowError(err, global.Window)
		return
	}

	fyne.Do(func() {
		ev.responseController.SetStatus(responsePayload.Status)
		ev.responseController.SetTime(responsePayload.Time)
		ev.responseController.SetBody(utils.SmartFormat(responsePayload.Body))
	})
}

func (ev ExchangeView) sendButtonHandler() {
	ev.headerView.SendButton.Disable()
	ev.responseController.SetLoading(true)

	ev.responseController.SetStatus(constants.UI_LOADING_RESPONSE_STATUS)
	ev.responseController.SetTime(constants.UI_LOADING_RESPONSE_TIME)
	ev.responseController.SetBody(constants.UI_LOADING_RESPONSE_BODY)

	requestPayload := ev.ToRequestPayload()

	go ev.sendRequestWorker(requestPayload)
}

func ComposeExchangeView() ExchangeView {
	headerView := exchangeheaderview.ComposeExchangeHeaderView()
	requestView := requestview.ComposeRequestView()
	responseController := response.NewController()

	headerView.MethodBinding.AddListener(binding.NewDataListener(func() {
		method, methodErr := headerView.MethodBinding.Get()
		if methodErr != nil {
			log.Error(methodErr)
			return
		}
		if method == constants.HTTP_METHOD_GET ||
			method == constants.HTTP_METHOD_HEAD {
			bodyTypeErr := requestView.Body.BodyTypeBinding.Set(constants.UI_BODY_TYPE_NONE)
			if bodyTypeErr != nil {
				log.Error(bodyTypeErr)
			}
			requestView.Body.DisableBodyTypeSelect()
		} else if method == constants.HTTP_METHOD_DELETE ||
			method == constants.HTTP_METHOD_OPTIONS ||
			method == constants.HTTP_METHOD_PATCH ||
			method == constants.HTTP_METHOD_POST ||
			method == constants.HTTP_METHOD_PUT {
			requestView.Body.EnableBodyTypeSelect()
		} else {
			log.Error(errors.New("invalid http method"))
		}
	}))

	requestResponseView := container.NewHSplit(requestView.UI, responseController.GetUI())
	ui := container.NewBorder(headerView.UI, nil, nil, nil, requestResponseView)

	ev := ExchangeView{
		ui,
		headerView,
		requestView,
		responseController,
	}

	headerView.SendButton.OnTapped = func() {
		ev.sendButtonHandler()
	}

	return ev
}
