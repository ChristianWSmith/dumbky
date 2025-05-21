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
		log.Error("", headerErr.Error())
		return ExchangeState{}, headerErr
	}
	request, requestErr := ev.RequestView.ToState()
	if requestErr != nil {
		log.Error("", requestErr.Error())
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
		log.Error("", headerErr.Error())
		return headerErr
	}
	requestErr := ev.RequestView.LoadState(exchangeState.Request)
	if requestErr != nil {
		log.Error("", requestErr.Error())
		return requestErr
	}
	return nil
}

func (ev ExchangeView) ToRequestPayload() request.RequestPayload {
	url, urlGetErr := ev.HeaderView.URLBinding.Get()
	if urlGetErr != nil {
		log.Error("Failed to Get URLBinding", urlGetErr.Error())
	}
	urlValidateErr := ev.HeaderView.ValidateURL()
	if urlValidateErr != nil {
		log.Warn("Failed to Validate URLEntry", urlValidateErr.Error())
	}

	method, methodGetErr := ev.HeaderView.MethodBinding.Get()
	if methodGetErr != nil {
		log.Error("Failed to Get MethodBinding", methodGetErr.Error())
	}

	useSSL, useSSLGetErr := ev.HeaderView.UseSSLBinding.Get()
	if useSSLGetErr != nil {
		log.Error("Failed to Get UseSSLBinding", useSSLGetErr.Error())
	}

	headers, headersGetErr := ev.RequestView.Headers.GetMap()
	if headersGetErr != nil {
		log.Error("Failed to Get Headers Map", headersGetErr.Error())
	}
	headersValidatErr := ev.RequestView.Headers.Validate()
	if headersValidatErr != nil {
		log.Error("Failed to Validate Headers Map", headersValidatErr.Error())
	}

	params, paramsGetErr := ev.RequestView.Params.GetMap()
	if paramsGetErr != nil {
		log.Error("Failed to Get Params Map", paramsGetErr.Error())
	}
	paramsValidatErr := ev.RequestView.Params.Validate()
	if paramsValidatErr != nil {
		log.Error("Failed to Validate Params Map", paramsValidatErr.Error())
	}

	bodyType, bodyTypeGetErr := ev.RequestView.Body.BodyTypeBinding.Get()
	if bodyTypeGetErr != nil {
		log.Error("Failed to Get BodyTypeBinding", bodyTypeGetErr.Error())
	}

	bodyRaw, bodyRawGetErr := ev.RequestView.Body.BodyRawBinding.Get()
	if bodyRawGetErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Error("Failed to Get BodyRawBinding", bodyRawGetErr.Error())
	}
	bodyRawValidateErr := ev.RequestView.Body.ValidateBodyRaw()
	if bodyRawValidateErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Warn("Failed to Validate BodyRawEntry", bodyRawValidateErr.Error())
	}

	bodyForm, bodyFormGetErr := ev.RequestView.Body.BodyKeyValueEditor.GetMap()
	if bodyFormGetErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn("Failed to Get Body Map", bodyFormGetErr.Error())
	}
	bodyFormValidateErr := ev.RequestView.Body.BodyKeyValueEditor.Validate()
	if bodyFormValidateErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn("Failed to Validate Body Map", bodyFormValidateErr.Error())
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
			log.Error("Failed to Get BodyTypeBinding during autoformat", bodyTypeGetErr.Error())
		}
		if bodyType != constants.UI_BODY_TYPE_RAW && bodyTypeGetErr == nil {
			log.Debug("Skipping autoformat for non-raw bodyType")
			return
		}
		bodyRaw, bodyRawGetErr := ev.RequestView.Body.BodyRawBinding.Get()
		if bodyRawGetErr != nil {
			log.Error("Failed to Get BodyRawBinding during autoformat", bodyRawGetErr.Error())
			return
		}
		bodyRawSetErr := ev.RequestView.Body.BodyRawBinding.Set(utils.SmartFormat(bodyRaw))
		if bodyRawSetErr != nil {
			log.Error("Failed to Set BodyRawBinding during autoformat", bodyRawSetErr.Error())
			return
		}
	})

	responsePayload, err := request.SendRequest(requestPayload)
	if err != nil {
		log.Error("Failed to SendRequest", err.Error())
		dialog.ShowError(err, global.Window)
		return
	}

	fyne.Do(func() {
		statusErr := ev.ResponseView.StatusBinding.Set(responsePayload.Status)
		if statusErr != nil {
			log.Error("Failed to set StatusBinding", statusErr.Error())
		}

		timeErr := ev.ResponseView.TimeBinding.Set(responsePayload.Time)
		if timeErr != nil {
			log.Error("Failed to set TimeBinding", timeErr.Error())
		}

		bodyErr := ev.ResponseView.BodyBinding.Set(utils.SmartFormat(responsePayload.Body))
		if bodyErr != nil {
			log.Error("Failed to set BodyBinding", bodyErr.Error())
		}
	})
}

func (ev ExchangeView) sendButtonHandler() {
	ev.HeaderView.SendButton.Disable()
	ev.ResponseView.SetLoading(true)

	statusErr := ev.ResponseView.StatusBinding.Set("")
	if statusErr != nil {
		log.Error("Failed to set StatusBinding", statusErr.Error())
	}
	timeErr := ev.ResponseView.TimeBinding.Set("")
	if timeErr != nil {
		log.Error("Failed to set TimeBinding", timeErr.Error())
	}
	bodyErr := ev.ResponseView.BodyBinding.Set("")
	if bodyErr != nil {
		log.Error("Failed to set BodyBinding", bodyErr.Error())
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
			log.Error("Failed to Get MethodBinding in DataListener", methodErr.Error())
			return
		}
		if method == constants.HTTP_METHOD_GET || method == constants.HTTP_METHOD_HEAD {
			bodyTypeErr := requestView.Body.BodyTypeBinding.Set(constants.UI_BODY_TYPE_NONE)
			if bodyTypeErr != nil {
				log.Error("Failed to Set BodyTypeBinding in DataListener", bodyTypeErr.Error())
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
