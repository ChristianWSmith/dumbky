package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/request"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
)

type ExchangeView struct {
	UI           *fyne.Container
	HeaderView   ExchangeHeaderView
	RequestView  RequestView
	ResponseView ResponseView
}

type ExchangeState struct {
	Header  ExchangeHeaderState `json:"header"`
	Request RequestState        `json:"request"`
}

func (ev ExchangeView) ToState() (ExchangeState, error) {
	header, headerErr := ev.HeaderView.ToState()
	if headerErr != nil {
		log.Error(headerErr)
		return ExchangeState{}, headerErr
	}
	request, requestErr := ev.RequestView.ToState()
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
	headerErr := ev.HeaderView.LoadState(exchangeState.Header)
	if headerErr != nil {
		log.Error(headerErr)
		return headerErr
	}
	requestErr := ev.RequestView.LoadState(exchangeState.Request)
	if requestErr != nil {
		log.Error(requestErr)
		return requestErr
	}
	return nil
}

func (ev ExchangeView) ToRequestPayload() request.RequestPayload {
	url, urlGetErr := ev.HeaderView.URLBinding.Get()
	if urlGetErr != nil {
		log.Error(urlGetErr)
	}
	urlValidateErr := ev.HeaderView.ValidateURL()
	if urlValidateErr != nil {
		log.Warn(urlValidateErr)
	}

	method, methodGetErr := ev.HeaderView.MethodBinding.Get()
	if methodGetErr != nil {
		log.Error(methodGetErr)
	}

	useSSL, useSSLGetErr := ev.HeaderView.UseSSLBinding.Get()
	if useSSLGetErr != nil {
		log.Error(useSSLGetErr)
	}

	headers, headersGetErr := ev.RequestView.Headers.GetMap()
	if headersGetErr != nil {
		log.Error(headersGetErr)
	}
	headersValidatErr := ev.RequestView.Headers.Validate()
	if headersValidatErr != nil {
		log.Error(headersValidatErr)
	}

	params, paramsGetErr := ev.RequestView.Params.GetMap()
	if paramsGetErr != nil {
		log.Error(paramsGetErr)
	}
	paramsValidatErr := ev.RequestView.Params.Validate()
	if paramsValidatErr != nil {
		log.Error(paramsValidatErr)
	}

	bodyType, bodyTypeGetErr := ev.RequestView.Body.BodyTypeBinding.Get()
	if bodyTypeGetErr != nil {
		log.Error(bodyTypeGetErr)
	}

	bodyRaw, bodyRawGetErr := ev.RequestView.Body.BodyRawBinding.Get()
	if bodyRawGetErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Error(bodyRawGetErr)
	}
	bodyRawValidateErr := ev.RequestView.Body.ValidateBodyRaw()
	if bodyRawValidateErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Warn(bodyRawValidateErr)
	}

	bodyForm, bodyFormGetErr := ev.RequestView.Body.BodyKeyValueEditor.GetMap()
	if bodyFormGetErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn(bodyFormGetErr)
	}
	bodyFormValidateErr := ev.RequestView.Body.BodyKeyValueEditor.Validate()
	if bodyFormValidateErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn(bodyFormValidateErr)
	}

	return request.RequestPayload{
		URL:      url,
		Method:   method,
		UseSSL:   useSSL,
		Headers:  headers,
		Params:   params,
		BodyType: bodyType,
		BodyRaw:  bodyRaw,
		BodyForm: bodyForm,
	}
}

func (ev ExchangeView) sendRequestWorker(requestPayload request.RequestPayload) {
	defer fyne.Do(func() {
		ev.ResponseView.SetLoading(false)
		ev.HeaderView.SendButton.Enable()
	})

	fyne.Do(func() {
		bodyType, bodyTypeGetErr := ev.RequestView.Body.BodyTypeBinding.Get()
		if bodyTypeGetErr != nil {
			log.Error(bodyTypeGetErr)
		}
		if bodyType != constants.UI_BODY_TYPE_RAW && bodyTypeGetErr == nil {
			log.Debug("Skipping autoformat for non-raw bodyType")
			return
		}
		bodyRaw, bodyRawGetErr := ev.RequestView.Body.BodyRawBinding.Get()
		if bodyRawGetErr != nil {
			log.Error(bodyRawGetErr)
			return
		}
		bodyRawSetErr := ev.RequestView.Body.BodyRawBinding.Set(utils.SmartFormat(bodyRaw))
		if bodyRawSetErr != nil {
			log.Error(bodyRawSetErr)
			return
		}
	})

	responsePayload, err := request.SendRequest(requestPayload)
	if err != nil {
		log.Warn(err)
		dialog.ShowError(err, global.Window)
		return
	}

	fyne.Do(func() {
		statusErr := ev.ResponseView.StatusBinding.Set(responsePayload.Status)
		if statusErr != nil {
			log.Error(statusErr)
		}

		timeErr := ev.ResponseView.TimeBinding.Set(responsePayload.Time)
		if timeErr != nil {
			log.Error(timeErr)
		}

		bodyErr := ev.ResponseView.BodyBinding.Set(utils.SmartFormat(responsePayload.Body))
		if bodyErr != nil {
			log.Error(bodyErr)
		}
	})
}

func (ev ExchangeView) sendButtonHandler() {
	ev.HeaderView.SendButton.Disable()
	ev.ResponseView.SetLoading(true)

	statusErr := ev.ResponseView.StatusBinding.Set("")
	if statusErr != nil {
		log.Error(statusErr)
	}
	timeErr := ev.ResponseView.TimeBinding.Set("")
	if timeErr != nil {
		log.Error(timeErr)
	}
	bodyErr := ev.ResponseView.BodyBinding.Set("")
	if bodyErr != nil {
		log.Error(bodyErr)
	}

	requestPayload := ev.ToRequestPayload()

	go ev.sendRequestWorker(requestPayload)
}

func ComposeExchangeView() ExchangeView {
	headerView := ComposeExchangeHeaderView()
	requestView := ComposeRequestView()
	responseView := ComposeResponseView()

	headerView.MethodBinding.AddListener(binding.NewDataListener(func() {
		method, methodErr := headerView.MethodBinding.Get()
		if methodErr != nil {
			log.Error(methodErr)
			return
		}
		if method == constants.HTTP_METHOD_GET || method == constants.HTTP_METHOD_HEAD {
			bodyTypeErr := requestView.Body.BodyTypeBinding.Set(constants.UI_BODY_TYPE_NONE)
			if bodyTypeErr != nil {
				log.Error(bodyTypeErr)
			}
			requestView.Body.DisableBodyTypeSelect()
		} else {
			requestView.Body.EnableBodyTypeSelect()
		}
	}))

	requestResponseView := container.NewHSplit(requestView.UI, responseView.UI)
	ui := container.NewBorder(headerView.UI, nil, nil, nil, requestResponseView)

	ev := ExchangeView{
		ui,
		headerView,
		requestView,
		responseView,
	}

	headerView.SendButton.OnTapped = func() {
		ev.sendButtonHandler()
	}

	return ev
}
